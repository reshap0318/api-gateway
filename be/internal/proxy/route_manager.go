package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/reshap0318/api-gateway/internal/database"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

// RateLimitConfig is a resolved (limit, window) pair ready to hand to
// helpers.RateLimiter.AllowWithLimit — no further chain logic needed at request time.
type RateLimitConfig struct {
	Limit      int
	WindowSecs int
}

// CachedRoute is the in-memory representation of a GatewayRoute ready for request matching.
type CachedRoute struct {
	ID          uint
	ServiceID   uint
	ServiceName string
	BaseURL     string
	// FullPath = BasePath+PathPattern; the actual path matched (see Refresh()).
	BasePath            string
	FullPath            string
	Protocol            string
	Method              string
	PathPattern         string
	segments            []pathSegment
	PermissionMatchMode string
	Permissions         []string

	// RateLimit is already fully resolved at Refresh()-time via the chain Route override →
	// Service default → global `.env` default (FSD §2.15) — the proxy handler just uses it
	// as-is, it never re-derives the chain per request.
	RateLimit RateLimitConfig

	// Proxy is built once per BaseURL at Refresh()-time and reused for every request matching
	// this route, instead of allocating a new httputil.ReverseProxy per request. Nil if
	// BaseURL failed to parse (proxyRequest reports 500 in that case).
	Proxy *httputil.ReverseProxy
}

// newReverseProxy builds a *httputil.ReverseProxy targeting baseURL, wired with the same
// upstream-unreachable JSON error response used across the Dynamic Proxy Engine.
func newReverseProxy(baseURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(baseURL)
	if err != nil || target.Scheme == "" || target.Host == "" {
		return nil, fmt.Errorf("invalid upstream base URL: %s", baseURL)
	}

	rp := httputil.NewSingleHostReverseProxy(target)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[gateway-proxy] upstream unreachable for %s %s -> %s: %v", r.Method, r.URL.Path, target.String(), err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"code":502,"message":"Upstream service unavailable"}`))
	}
	return rp, nil
}

// RouteManager holds the in-memory cache of active Service+Route config used by the
// Dynamic Proxy Engine. Refresh uses an atomic swap: a new snapshot is built in a local
// variable first, and only swapped into the active cache after the build succeeds — the
// cache that is currently serving traffic is never cleared/mutated in place, so it is
// never seen empty or partial (see docs/04_TDD.md §5 Konkurensi & Caching).
type RouteManager struct {
	mu     sync.RWMutex
	routes []*CachedRoute
	// byBasePath indexes routes by their owning Service's base_path, so Match() only scans
	// the routes of the matched Service instead of every route in the gateway.
	byBasePath map[string][]*CachedRoute

	repo   *repositories.GatewayServiceRepository
	logger *helpers.Logger

	redis   *database.RedisCache
	channel string

	globalRateLimit RateLimitConfig

	// instanceID identifies this process for the manual cache status endpoint (FSD §4.5) —
	// useful for confirming which instance answered when running multiple instances behind a
	// load balancer, since each instance keeps its own independent in-memory route cache.
	instanceID string
	// lastRefreshed is when Refresh() last completed successfully.
	lastRefreshed time.Time
	// lastVersion is the gateway_cache_meta version as of the last Refresh() — read before
	// the DB query for services/routes, so it's always <= the true data version (never
	// ahead of it), guaranteeing RefreshIfStale can't miss a mutation that raced with a Refresh.
	lastVersion uint64

	stopTicker chan struct{}
	stopPubSub chan struct{}
}

// NewRouteManager creates a new RouteManager.
func NewRouteManager(repo *repositories.GatewayServiceRepository, logger *helpers.Logger) *RouteManager {
	return &RouteManager{
		repo:       repo,
		logger:     logger,
		routes:     []*CachedRoute{},
		instanceID: newInstanceID(),
	}
}

// newInstanceID builds a per-process identifier from hostname + PID — stable for the life of
// the process, unique enough across instances/hosts for the cache status endpoint.
func newInstanceID() string {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	return fmt.Sprintf("%s:%d", host, os.Getpid())
}

// SetRedis wires the Redis client + Pub/Sub channel used for multi-instance cache sync
// (FSD §2.21). Must be called before StartRedisSubscriber/RefreshAndPublish; safe to skip
// entirely when Redis is disabled — RefreshAndPublish then behaves like a plain Refresh().
func (rm *RouteManager) SetRedis(redis *database.RedisCache, channel string) {
	rm.redis = redis
	rm.channel = channel
}

// SetGlobalRateLimit sets the `.env` fallback (RATE_LIMIT_REQUESTS/RATE_LIMIT_WINDOW) used
// to resolve CachedRoute.RateLimit whenever neither the Route nor its Service define an
// override. Must be called before the first Refresh() to take effect on that build.
func (rm *RouteManager) SetGlobalRateLimit(cfg RateLimitConfig) {
	rm.globalRateLimit = cfg
}

// Refresh rebuilds the route cache from the DB and atomically swaps it in.
// On query failure, the existing cache is kept untouched (fail-safe, not fail-empty).
func (rm *RouteManager) Refresh() error {
	// Read the version before querying services/routes: if a CUD races with this Refresh,
	// the version read here stays <= the true data version, so RefreshIfStale can never
	// wrongly think a racing mutation is already reflected in this build.
	version, verr := rm.repo.GetCacheVersion(nil)
	if verr != nil && rm.logger != nil {
		rm.logger.LogWarn("RouteManager.Refresh", "Failed to read cache version: %v", verr)
	}

	services, err := rm.repo.FindAllActiveWithRoutes(nil)
	if err != nil {
		if rm.logger != nil {
			rm.logger.LogWarn("RouteManager.Refresh", "Failed to load routes, keeping existing cache: %v", err)
		}
		return err
	}

	newRoutes := make([]*CachedRoute, 0)
	proxyCache := make(map[string]*httputil.ReverseProxy) // keyed by BaseURL, shared across routes of the same Service
	for _, svc := range services {
		proxy, ok := proxyCache[svc.BaseURL]
		if !ok {
			built, err := newReverseProxy(svc.BaseURL)
			if err != nil && rm.logger != nil {
				rm.logger.LogWarn("RouteManager.Refresh", "Skipping proxy for service %s: %v", svc.Name, err)
			}
			proxy = built
			proxyCache[svc.BaseURL] = proxy
		}

		for _, rt := range svc.Routes {
			// Resolution chain (FSD §2.15): Route override → Service default → global.
			// Route/Service values are "per minute" by definition, so window is fixed at 60s;
			// only the global fallback uses the configured RATE_LIMIT_WINDOW.
			rateLimit := rm.globalRateLimit
			if svc.RateLimitPerMinute != nil {
				rateLimit = RateLimitConfig{Limit: *svc.RateLimitPerMinute, WindowSecs: 60}
			}
			if rt.RateLimitPerMinute != nil {
				rateLimit = RateLimitConfig{Limit: *rt.RateLimitPerMinute, WindowSecs: 60}
			}

			// path_pattern is relative to the Service's base_path (FSD §2.13).
			fullPath := svc.BasePath + rt.PathPattern

			cr := &CachedRoute{
				ID:                  rt.ID,
				ServiceID:           svc.ID,
				ServiceName:         svc.Name,
				BaseURL:             svc.BaseURL,
				BasePath:            svc.BasePath,
				FullPath:            fullPath,
				Protocol:            svc.Protocol,
				Method:              rt.Method,
				PathPattern:         rt.PathPattern,
				segments:            parsePattern(fullPath),
				PermissionMatchMode: rt.PermissionMatchMode,
				RateLimit:           rateLimit,
				Proxy:               proxy,
			}
			for _, p := range rt.Permissions {
				cr.Permissions = append(cr.Permissions, p.Name)
			}
			newRoutes = append(newRoutes, cr)
		}
	}

	byBasePath := make(map[string][]*CachedRoute)
	for _, cr := range newRoutes {
		byBasePath[cr.BasePath] = append(byBasePath[cr.BasePath], cr)
	}

	// Atomic swap — replaced under lock only after the full build above has succeeded;
	// readers via Match() never observe a partial state.
	rm.mu.Lock()
	rm.routes = newRoutes
	rm.byBasePath = byBasePath
	rm.lastRefreshed = time.Now()
	if verr == nil {
		rm.lastVersion = version
	}
	rm.mu.Unlock()

	if rm.logger != nil {
		rm.logger.LogInfo("RouteManager.Refresh", "Cache refreshed: %d route(s) from %d service(s)", len(newRoutes), len(services))
	}

	return nil
}

// RefreshAndPublish refreshes the local cache, then broadcasts a refresh signal to all
// other Gateway instances via Redis Pub/Sub (FSD §2.18/§2.21 on-save trigger). This is the
// method the on-save trigger (services.RouteCacheRefresher) must call — NOT Refresh()
// directly — otherwise other instances never learn about the change. Conversely, the
// periodic ticker and the Pub/Sub subscriber itself must keep calling plain Refresh(),
// or every relayed message would re-publish and create an infinite broadcast loop.
func (rm *RouteManager) RefreshAndPublish() error {
	if err := rm.Refresh(); err != nil {
		return err
	}

	if rm.redis != nil && rm.redis.IsCacheAvailable() && rm.channel != "" {
		payload := fmt.Sprintf(`{"type":"route_refresh","triggered_at":%q}`, time.Now().Format(time.RFC3339))
		if err := rm.redis.Publish(rm.channel, payload); err != nil && rm.logger != nil {
			rm.logger.LogWarn("RouteManager.RefreshAndPublish", "Failed to publish refresh signal: %v", err)
		}
	}

	return nil
}

// StartRedisSubscriber subscribes to the configured Pub/Sub channel and refreshes the
// local cache (Refresh only, never RefreshAndPublish) whenever any instance publishes a
// change. No-op if SetRedis was never called or Redis is unavailable.
func (rm *RouteManager) StartRedisSubscriber() {
	if rm.redis == nil || !rm.redis.IsCacheAvailable() || rm.channel == "" {
		return
	}

	pubsub := rm.redis.Subscribe(rm.channel)
	rm.stopPubSub = make(chan struct{})

	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				if rm.logger != nil {
					rm.logger.LogInfo("RouteManager.StartRedisSubscriber", "Refresh signal received: %s", msg.Payload)
				}
				_ = rm.Refresh()
			case <-rm.stopPubSub:
				return
			}
		}
	}()

	if rm.logger != nil {
		rm.logger.LogInfo("RouteManager.StartRedisSubscriber", "Subscribed to channel: %s", rm.channel)
	}
}

// MustRefreshSync loads the cache synchronously and fails fast if the DB is unreachable.
// Intended to be called once at startup, before the server starts accepting traffic.
func (rm *RouteManager) MustRefreshSync() error {
	if err := rm.Refresh(); err != nil {
		return fmt.Errorf("initial route cache load failed: %w", err)
	}
	return nil
}

// Match finds the best (most specific) route for the given method + path. Only scans the
// routes of the Service whose base_path is the longest matching prefix of path, instead of
// every route in the gateway.
func (rm *RouteManager) Match(method, path string) (*CachedRoute, map[string]string, bool) {
	rm.mu.RLock()
	byBasePath := rm.byBasePath
	rm.mu.RUnlock()

	var routes []*CachedRoute
	bestBasePathLen := -1
	for basePath, group := range byBasePath {
		if path == basePath || strings.HasPrefix(path, basePath+"/") {
			if len(basePath) > bestBasePathLen {
				bestBasePathLen = len(basePath)
				routes = group
			}
		}
	}
	if routes == nil {
		return nil, nil, false
	}

	var best *CachedRoute
	var bestParams map[string]string
	bestScore := -1

	for _, r := range routes {
		if r.Method != "*" && r.Method != method {
			continue
		}

		params, score, ok := matchSegments(r.segments, path)
		if !ok {
			continue
		}

		// Exact method match outranks a wildcard ("*") method route.
		if r.Method == method {
			score += 1000
		}

		if score > bestScore {
			bestScore = score
			best = r
			bestParams = params
		}
	}

	if best == nil {
		return nil, nil, false
	}
	return best, bestParams, true
}

// RefreshIfStale checks gateway_cache_meta's version and only runs the full Refresh()
// rebuild if it moved since the last refresh. Used by the periodic ticker; on-save/Pub-Sub
// triggers already know a change happened and call Refresh() directly.
func (rm *RouteManager) RefreshIfStale() error {
	version, err := rm.repo.GetCacheVersion(nil)
	if err != nil {
		if rm.logger != nil {
			rm.logger.LogWarn("RouteManager.RefreshIfStale", "Version check failed, refreshing anyway: %v", err)
		}
		return rm.Refresh()
	}

	rm.mu.RLock()
	lastVersion := rm.lastVersion
	rm.mu.RUnlock()

	if version == lastVersion {
		return nil
	}

	if rm.logger != nil {
		rm.logger.LogInfo("RouteManager.RefreshIfStale", "Version changed (%d -> %d), refreshing", lastVersion, version)
	}
	return rm.Refresh()
}

// StartPeriodicRefresh starts a background ticker that refreshes the cache on an interval,
// as a fallback safety net independent of on-save/Pub-Sub triggers.
func (rm *RouteManager) StartPeriodicRefresh(interval time.Duration) {
	rm.stopTicker = make(chan struct{})
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = rm.RefreshIfStale()
			case <-rm.stopTicker:
				return
			}
		}
	}()
}

// Stop stops the periodic refresh ticker and Redis subscriber goroutines.
func (rm *RouteManager) Stop() {
	if rm.stopTicker != nil {
		close(rm.stopTicker)
	}
	if rm.stopPubSub != nil {
		close(rm.stopPubSub)
	}
}

// Stats returns the counters + metadata for the manual cache status endpoint (FSD §4.5:
// last_refreshed_at, total_services, total_routes, instance_id).
func (rm *RouteManager) Stats() (totalRoutes int, totalServices int, lastRefreshed time.Time, instanceID string) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	seen := make(map[uint]struct{})
	for _, r := range rm.routes {
		seen[r.ServiceID] = struct{}{}
	}
	return len(rm.routes), len(seen), rm.lastRefreshed, rm.instanceID
}
