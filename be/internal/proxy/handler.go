package proxy

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/middleware"
	"github.com/reshap0318/api-gateway/internal/services"
)

// Handler returns a gin.HandlerFunc meant to be registered via gin.Engine.NoRoute — any
// request that doesn't match a Gateway Management API route (§4.1–§4.5) falls through here
// and is resolved dynamically against the RouteManager cache, then reverse-proxied to the
// matched upstream Service. REST and WebSocket share this exact same code path:
// httputil.ReverseProxy natively hijacks the connection on `Connection: Upgrade` (Go 1.12+),
// so no separate WebSocket handling is required.
func Handler(rm *RouteManager, svcs *services.Services, acc *helpers.Access, limiter *helpers.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := authenticate(c, svcs)
		if !ok {
			return
		}

		if hasDotDotSegment(c.Request.URL.Path) {
			helpers.BadRequest(c, "Invalid path")
			return
		}

		route, params, found := rm.Match(c.Request.Method, c.Request.URL.Path)
		if !found {
			helpers.NotFound(c, "Route not found")
			return
		}

		if len(route.Permissions) > 0 {
			ctx := c.Request.Context()
			var allowed bool
			if route.PermissionMatchMode == "all" {
				allowed = acc.HasAllPermissions(ctx, route.Permissions...)
			} else {
				allowed = acc.HasPermission(ctx, route.Permissions...)
			}
			if !allowed {
				helpers.Forbidden(c, "Permission denied")
				return
			}
		}

		if !enforceRateLimit(c, route, limiter) {
			return
		}

		injectCallerHeaders(c, claims, acc)
		proxyRequest(c, route, params)
	}
}

// enforceRateLimit checks CachedRoute.RateLimit — already fully resolved at Refresh()-time
// via the chain Route override → Service default → global `.env` default (FSD §2.15) — via
// a per-(route,client) token bucket, so overriding the limit for one Route never affects
// another Route's or client's budget.
func enforceRateLimit(c *gin.Context, route *CachedRoute, limiter *helpers.RateLimiter) bool {
	// c.ClientIP() only trusts X-Forwarded-For/X-Real-IP when the immediate peer is a
	// configured trusted proxy (main.go), preventing clients from spoofing a fresh IP per
	// request to get an unlimited supply of new token buckets.
	ip := c.ClientIP()
	key := fmt.Sprintf("gwrl:route:%d:%s", route.ID, ip)

	info := limiter.AllowWithLimit(key, route.RateLimit.Limit, route.RateLimit.WindowSecs)

	c.Header("X-RateLimit-Limit", strconv.Itoa(info.Limit))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(info.Remaining))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(info.Reset, 10))

	if !info.Allowed {
		retryAfter := int(info.Reset - time.Now().Unix())
		if retryAfter < 0 {
			retryAfter = 0
		}
		c.Header("Retry-After", strconv.Itoa(retryAfter))
		c.JSON(http.StatusTooManyRequests, gin.H{"code": http.StatusTooManyRequests, "message": "Rate limit exceeded"})
		return false
	}

	return true
}

// authenticate validates the Bearer JWT (via the same middleware.ExtractBearerClaims used by
// middleware.JWTAuth for the Management API — NoRoute bypasses group-scoped middleware, so this
// must be invoked inline here) and, on success, attaches the caller ID to the request context
// and returns the decoded claims for the caller to use (e.g. injectCallerHeaders).
func authenticate(c *gin.Context, svcs *services.Services) (*helpers.JWTClaims, bool) {
	claims, err := middleware.ExtractBearerClaims(c, svcs)
	if err != nil {
		helpers.Unauthorized(c, err.Error())
		return nil, false
	}

	ctx := context.WithValue(c.Request.Context(), helpers.KeyUserID, claims.UserID)
	c.Request = c.Request.WithContext(ctx)

	c.Set("user_id", claims.UserID)
	c.Set("user_email", claims.Email)

	return claims, true
}

// injectCallerHeaders sets X-User-Id/X-User-Email/X-User-Name/X-User-Roles/X-User-Permissions/
// X-Trace-Id on the request being forwarded upstream, so the upstream service can identify the
// caller without having to decode/verify the JWT itself, and can tag its own logs with the same
// trace ID middleware.TraceID assigned at request entry — so one request's path through gateway
// + service is correlatable end-to-end. Roles/Permissions come from acc.GetUserAccess (live,
// cached DB data) rather than the JWT claims, so a Role/Permission change takes effect on the
// very next proxied request instead of waiting for the caller to refresh their token; if the
// live lookup fails (e.g. user deleted after login), the JWT claims are used as a fallback so
// the request still proxies. Roles/Permissions are comma-joined (role/permission names never
// contain commas, e.g. "toko.publish"), matching the standard HTTP list-header convention (cf.
// Accept-Encoding). Header.Set (not Add) is used deliberately: it replaces any of these the
// client may have sent themselves, so a caller can't spoof another user's identity/permissions
// to the upstream.
func injectCallerHeaders(c *gin.Context, claims *helpers.JWTClaims, acc *helpers.Access) {
	roles, permissions, ok := acc.GetUserAccess(c.Request.Context())
	if !ok {
		roles, permissions = claims.Roles, claims.Permissions
	}

	c.Request.Header.Set("X-User-Id", strconv.FormatUint(uint64(claims.UserID), 10))
	c.Request.Header.Set("X-User-Email", claims.Email)
	c.Request.Header.Set("X-User-Name", claims.Name)
	c.Request.Header.Set("X-Trace-Id", helpers.GetTraceID(c.Request.Context()))
	c.Request.Header.Set("X-User-Roles", strings.Join(roles, ","))
	c.Request.Header.Set("X-User-Permissions", strings.Join(permissions, ","))
}

// proxyRequest forwards the request to {service.base_url}{path_pattern-relative path} and
// streams the upstream response back to the client. The Service's base_path is a
// Gateway-only routing prefix (FSD §2.13) used to pick the matching Service — it is stripped
// before forwarding, so the upstream only ever sees its own path_pattern (e.g. a request to
// "/master/category" reaches the upstream as "/category"). params (extracted :param values)
// are resolved but not rewritten — the remainder of the path is proxied verbatim.
// route.Proxy is built once per BaseURL at RouteManager.Refresh()-time and shared across
// requests/routes, rather than allocating a new httputil.ReverseProxy per request.
func proxyRequest(c *gin.Context, route *CachedRoute, params map[string]string) {
	_ = params

	if route.Proxy == nil {
		helpers.InternalServerError(c, "Invalid upstream base URL")
		return
	}

	trimmed := strings.TrimPrefix(c.Request.URL.Path, route.BasePath)
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	c.Request.URL.Path = trimmed
	c.Request.URL.RawPath = ""

	route.Proxy.ServeHTTP(c.Writer, c.Request)
}
