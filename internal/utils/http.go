package utils

import (
	"context"
	"net/http"
	"strings"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const pathParamsKey contextKey = "pathParams"

// Route holds information about a registered route.
type Route struct {
	Method  string
	Path    string // Can contain placeholders like /path/{id}
	Handler http.HandlerFunc
}

// Router is a simple HTTP multiplexer.
type Router struct {
	routes          []Route
	NotFoundHandler http.HandlerFunc
}

// NewRouter creates a new Router.
func NewRouter() *Router {
	return &Router{}
}

// HandleFunc registers a new route with a handler function for a specific method.
func (rt *Router) HandleFunc(path string, handler http.HandlerFunc, method string) {
	rt.routes = append(rt.routes, Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

// PathPrefix registers a handler for a path prefix (used for static files).
func (rt *Router) PathPrefix(prefix string, handler http.Handler) {
	rt.routes = append(rt.routes, Route{
		Method: "", // Match any method
		Path:   strings.TrimSuffix(prefix, "/") + "/*",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
		},
	})
}

// PathVariable extracts a path variable from the request context.
func PathVariable(r *http.Request, key string) string {
	if vars, ok := r.Context().Value(pathParamsKey).(map[string]string); ok {
		return vars[key]
	}
	return ""
}

// ServeHTTP dispatches the request to the handler whose path pattern matches.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rt.routes {
		// Check method first
		if route.Method != "" && route.Method != r.Method {
			continue
		}

		// Check if it's a prefix route (for static assets)
		if strings.HasSuffix(route.Path, "/*") {
			prefix := strings.TrimSuffix(route.Path, "/*")
			if strings.HasPrefix(r.URL.Path, prefix) {
				route.Handler(w, r)
				return
			}
		} else {
			// Check for exact path match or path with parameters
			if pathParams, matched := matchPath(route.Path, r.URL.Path); matched {
				ctx := context.WithValue(r.Context(), pathParamsKey, pathParams)
				route.Handler(w, r.WithContext(ctx))
				return
			}
		}
	}

	// No route matched, use NotFoundHandler or default
	if rt.NotFoundHandler != nil {
		rt.NotFoundHandler(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// matchPath checks if a route path matches a request path and extracts parameters.
func matchPath(routePath, requestPath string) (map[string]string, bool) {
	// Handle root path
	if routePath == "/" && requestPath == "/" {
		return make(map[string]string), true
	}

	routeParts := strings.Split(strings.Trim(routePath, "/"), "/")
	requestParts := strings.Split(strings.Trim(requestPath, "/"), "/")

	if len(routeParts) != len(requestParts) {
		return nil, false
	}

	params := make(map[string]string)
	for i, routePart := range routeParts {
		if strings.HasPrefix(routePart, "{") && strings.HasSuffix(routePart, "}") {
			// Extract parameter name and store value
			paramName := strings.Trim(routePart, "{}")
			params[paramName] = requestParts[i]
		} else if routePart != requestParts[i] {
			// Exact match required for non-parameter segments
			return nil, false
		}
	}
	return params, true
}
