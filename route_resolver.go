package gbcobservehttp

import (
	"sort"
	"strings"
	"sync"

	goweb "goark.dev/goark/web"
)

type routeResolver struct {
	registry *goweb.Registry
	once     sync.Once
	routes   []routePattern
}

type routePattern struct {
	method, pattern string
	score, order    int
}

func (r *routeResolver) Resolve(method, path string) string {
	if r == nil || r.registry == nil {
		return ""
	}
	r.once.Do(r.compile)
	method = strings.ToUpper(method)
	for _, route := range r.routes {
		if route.method == method && matchRoute(route.pattern, path) {
			return route.pattern
		}
	}
	return ""
}
func (r *routeResolver) compile() {
	routes := r.registry.Routes()
	compiled := make([]routePattern, 0, len(routes))
	for index, route := range routes {
		compiled = append(compiled, routePattern{method: route.Method, pattern: route.Pattern, score: literalSegments(route.Pattern), order: index})
	}
	sort.SliceStable(compiled, func(i, j int) bool {
		if compiled[i].score != compiled[j].score {
			return compiled[i].score > compiled[j].score
		}
		return compiled[i].order < compiled[j].order
	})
	r.routes = compiled
}
func literalSegments(pattern string) int {
	score := 0
	for _, segment := range strings.Split(strings.Trim(pattern, "/"), "/") {
		if segment != "" && !strings.HasPrefix(segment, "{") {
			score++
		}
	}
	return score
}
func matchRoute(pattern, path string) bool {
	pattern = strings.Trim(pattern, "/")
	path = strings.Trim(path, "/")
	for {
		patternSegment, patternRest, patternMore := strings.Cut(pattern, "/")
		pathSegment, pathRest, pathMore := strings.Cut(path, "/")
		if patternSegment != pathSegment && !(strings.HasPrefix(patternSegment, "{") && strings.HasSuffix(patternSegment, "}") && pathSegment != "") {
			return false
		}
		if patternMore != pathMore {
			return false
		}
		if !patternMore {
			return true
		}
		pattern, path = patternRest, pathRest
	}
}
