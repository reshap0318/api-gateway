package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/helpers"
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
		if !authenticate(c, svcs) {
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

		proxyRequest(c, route, params)
	}
}

// enforceRateLimit checks CachedRoute.RateLimit — already fully resolved at Refresh()-time
// via the chain Route override → Service default → global `.env` default (FSD §2.15) — via
// a per-(route,client) token bucket, so overriding the limit for one Route never affects
// another Route's or client's budget.
func enforceRateLimit(c *gin.Context, route *CachedRoute, limiter *helpers.RateLimiter) bool {
	ip := helpers.GetClientIP(c.Request, c.ClientIP())
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

// authenticate validates the Bearer JWT and, on success, attaches the caller ID to the
// request context (mirrors middleware.JWTAuth — NoRoute bypasses group-scoped middleware,
// so this must be done inline here).
func authenticate(c *gin.Context, svcs *services.Services) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		helpers.Unauthorized(c, "Authorization header required")
		return false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		helpers.Unauthorized(c, "Invalid authorization header format")
		return false
	}

	claims, err := svcs.AuthValidateToken(parts[1])
	if err != nil {
		helpers.Unauthorized(c, "Invalid or expired token")
		return false
	}

	ctx := context.WithValue(c.Request.Context(), helpers.KeyUserID, claims.UserID)
	c.Request = c.Request.WithContext(ctx)
	return true
}

// proxyRequest forwards the request as-is to {service.base_url}{original path} and streams
// the upstream response back to the client. params (extracted :param values) are resolved
// but not rewritten into the forwarded path — the full original path is proxied verbatim.
func proxyRequest(c *gin.Context, route *CachedRoute, params map[string]string) {
	_ = params

	target, err := url.Parse(route.BaseURL)
	if err != nil || target.Scheme == "" || target.Host == "" {
		helpers.InternalServerError(c, "Invalid upstream base URL")
		return
	}

	rp := httputil.NewSingleHostReverseProxy(target)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[gateway-proxy] upstream unreachable for %s %s -> %s: %v", r.Method, r.URL.Path, target.String(), err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"code":502,"message":"Upstream service unavailable"}`))
	}

	rp.ServeHTTP(c.Writer, c.Request)
}
