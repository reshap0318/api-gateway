package services

// RouteCacheRefresher is implemented by proxy.RouteManager. Defined here (instead of
// importing the proxy package directly) to avoid an import cycle: proxy.Handler needs
// *services.Services for JWT validation, so services cannot import proxy back.
type RouteCacheRefresher interface {
	Refresh() error
}

// RefreshRouteCache triggers an on-save local cache refresh (FSD §2.18). Errors are logged
// only — a stale cache is safer than failing the CUD operation that already committed to DB.
func (s *Services) RefreshRouteCache(caller string) {
	if s.RouteCache == nil {
		return
	}
	if err := s.RouteCache.Refresh(); err != nil {
		s.Logger.LogWarn(caller, "On-save route cache refresh failed: %v", err)
	}
}
