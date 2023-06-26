package router

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/sdpsagarpawar/logger"
)

type Router struct {
	routes          map[string]map[string]*Route
	notFoundHandler http.HandlerFunc
	middleware      []func(http.HandlerFunc) http.HandlerFunc
	logger          *logger.Logger
}

type Route struct {
	HandlerFunc http.HandlerFunc
	Response    http.HandlerFunc
}

// NewRouter creates a new instance of Router.
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]*Route),
		logger: logger.NewLogger(), // Create a new logger instance
	}
}

// AddRoute adds a new route to the router with the specified HTTP method.
func (r *Router) AddRoute(method string, path string, handler http.HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]*Route)
	}
	r.routes[method][path] = &Route{
		HandlerFunc: handler,
	}
}

// SetResponse sets the response for a specific route.
func (r *Router) SetResponse(method string, path string, response http.HandlerFunc) {
	if r.routes[method] != nil && r.routes[method][path] != nil {
		r.routes[method][path].Response = response
	}
}

// SetNotFoundHandler sets the handler for the not found route.
func (r *Router) SetNotFoundHandler(handler http.HandlerFunc) {
	r.notFoundHandler = handler
}

// Use adds middleware to the router.
func (r *Router) Use(middleware ...func(http.HandlerFunc) http.HandlerFunc) {
	r.middleware = append(r.middleware, middleware...)
}

// ServeHTTP handles the incoming HTTP requests.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var route *Route

	// Determine the appropriate route based on the requested method and path
	if routes, ok := r.routes[req.Method]; ok {
		if r, ok := routes[req.URL.Path]; ok {
			route = r
		}
	}

	// If no route found, use the not found handler or default to http.NotFound
	if route == nil {
		if r.notFoundHandler != nil {
			route = &Route{
				HandlerFunc: r.notFoundHandler,
			}
		} else {
			route = &Route{
				HandlerFunc: http.NotFound,
			}
		}
	}

	// Apply middleware in reverse order
	handler := route.HandlerFunc
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}

	// Generate correlation ID using UUID
	correlationID := uuid.New().String()

	// Set correlation ID in request context
	ctx := req.Context()
	ctx = context.WithValue(ctx, "correlationID", correlationID)
	req = req.WithContext(ctx)

	// Parse query parameters
	queryParams, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		r.logger.Errorf("Failed to parse query parameters: %v", err)
	}

	// Add query parameters to the request context
	ctx = req.Context()
	ctx = context.WithValue(ctx, "queryParams", queryParams)
	req = req.WithContext(ctx)

	// Call the handler with the modified request
	handler(w, req)

	// Set the response for the route
	if route.Response != nil {
		route.Response(w, req)
	}
}

// GetQueryParams retrieves the query parameters from the request.
func (r *Router) GetQueryParams(req *http.Request) url.Values {
	queryParams, ok := req.Context().Value("queryParams").(url.Values)
	if !ok {
		return nil
	}
	return queryParams
}

// GetFormParams retrieves the form parameters from the request.
func (r *Router) GetFormParams(req *http.Request) (url.Values, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, err
	}
	return req.Form, nil
}

// GetCorrelationID retrieves the correlation ID from the request.
func (r *Router) GetCorrelationID(req *http.Request) string {
	correlationID, ok := req.Context().Value("correlationID").(string)
	if !ok {
		return ""
	}
	return correlationID
}
