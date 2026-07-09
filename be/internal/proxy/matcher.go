package proxy

import "strings"

type segmentKind int

const (
	segmentLiteral segmentKind = iota
	segmentParam
	segmentWildcard
)

type pathSegment struct {
	kind  segmentKind
	value string
}

// parsePattern splits a route path_pattern (e.g. "/user/:id", "/user/*") into segments.
func parsePattern(pattern string) []pathSegment {
	raw := strings.Split(strings.Trim(pattern, "/"), "/")
	segments := make([]pathSegment, 0, len(raw))
	for _, part := range raw {
		if part == "" {
			continue
		}
		switch {
		case part == "*":
			segments = append(segments, pathSegment{kind: segmentWildcard})
		case strings.HasPrefix(part, ":"):
			segments = append(segments, pathSegment{kind: segmentParam, value: part[1:]})
		default:
			segments = append(segments, pathSegment{kind: segmentLiteral, value: part})
		}
	}
	return segments
}

// matchSegments matches a request path against parsed pattern segments.
// Returns extracted params, a specificity score (higher = more specific), and whether it matched.
// Specificity: literal segment > param segment > wildcard segment.
func matchSegments(pattern []pathSegment, requestPath string) (map[string]string, int, bool) {
	rawPath := strings.Split(strings.Trim(requestPath, "/"), "/")
	pathParts := make([]string, 0, len(rawPath))
	for _, p := range rawPath {
		if p != "" {
			pathParts = append(pathParts, p)
		}
	}

	params := make(map[string]string)
	score := 0

	for i, seg := range pattern {
		if seg.kind == segmentWildcard {
			// Wildcard must be the trailing segment; matches everything remaining (including nothing).
			score++
			return params, score, true
		}

		if i >= len(pathParts) {
			return nil, 0, false
		}

		switch seg.kind {
		case segmentLiteral:
			if seg.value != pathParts[i] {
				return nil, 0, false
			}
			score += 3
		case segmentParam:
			params[seg.value] = pathParts[i]
			score += 2
		}
	}

	// No wildcard consumed — path must be fully consumed too (exact segment count match).
	if len(pathParts) != len(pattern) {
		return nil, 0, false
	}

	return params, score, true
}
