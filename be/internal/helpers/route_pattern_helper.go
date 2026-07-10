package helpers

import "strings"

// ValidateRoutePathPattern enforces that a Gateway Route's path_pattern uses "*" only as the
// trailing segment (e.g. "/files/*" is valid, "/files/*/meta" is not). proxy.matchSegments
// matches greedily and returns as soon as it hits a wildcard segment, silently ignoring any
// pattern segments after it — so an admin registering "/a/*/b" would get a route that actually
// behaves like "/a/*" with no error, which is surprising and worth rejecting up front.
func ValidateRoutePathPattern(pattern string) error {
	raw := strings.Split(strings.Trim(pattern, "/"), "/")
	segments := make([]string, 0, len(raw))
	for _, part := range raw {
		if part != "" {
			segments = append(segments, part)
		}
	}

	for i, part := range segments {
		if part == "*" && i != len(segments)-1 {
			return &FieldError{Field: "path_pattern", Message: "Wildcard (*) hanya boleh berada di segmen terakhir path_pattern"}
		}
	}

	return nil
}
