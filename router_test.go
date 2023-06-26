package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	// Create a new router instance
	router := NewRouter()

	// Test handler for the "GET /hello" route
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Test handler for the "POST /users" route
	createUserHandler := func(w http.ResponseWriter, req *http.Request) {
		// Simulate creating a user
		// ...

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created"))
	}

	// Test handler for the not found route
	notFoundHandler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
	}

	// Add routes to the router
	router.AddRoute("GET", "/hello", helloHandler)
	router.AddRoute("POST", "/users", createUserHandler)
	router.SetNotFoundHandler(notFoundHandler)

	// Test cases

	t.Run("Valid GET request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/hello", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
		}

		// Check the response body
		expectedBody := "Hello, World!"
		if rr.Body.String() != expectedBody {
			t.Errorf("Expected response body %q, but got %q", expectedBody, rr.Body.String())
		}
	})

	t.Run("Valid POST request", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/users", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, but got %d", http.StatusCreated, rr.Code)
		}

		// Check the response body
		expectedBody := "User created"
		if rr.Body.String() != expectedBody {
			t.Errorf("Expected response body %q, but got %q", expectedBody, rr.Body.String())
		}
	})

	t.Run("Not found route", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/unknown", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, rr.Code)
		}

		// Check the response body
		expectedBody := "404 Not Found"
		if rr.Body.String() != expectedBody {
			t.Errorf("Expected response body %q, but got %q", expectedBody, rr.Body.String())
		}
	})

	t.Run("Middleware", func(t *testing.T) {
		// Test middleware
		middleware := func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("X-Middleware", "true")
				next(w, req)
			}
		}

		// Add middleware to the router
		router.Use(middleware)

		// Test handler with middleware
		handlerWithMiddleware := func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Handler with middleware"))
		}

		// Add route with middleware
		router.AddRoute("GET", "/middleware", handlerWithMiddleware)

		req, err := http.NewRequest("GET", "/middleware", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
		}

		// Check the response header set by the middleware
		middlewareValue := rr.Header().Get("X-Middleware")
		if middlewareValue != "true" {
			t.Errorf("Expected middleware value %q, but got %q", "true", middlewareValue)
		}

		// Check the response body
		expectedBody := "Handler with middleware"
		if rr.Body.String() != expectedBody {
			t.Errorf("Expected response body %q, but got %q", expectedBody, rr.Body.String())
		}
	})

	t.Run("Correlation ID", func(t *testing.T) {
		// Test handler with correlation ID
		handlerWithCorrelationID := func(w http.ResponseWriter, req *http.Request) {
			correlationID := router.GetCorrelationID(req)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(correlationID))
		}

		// Add route with correlation ID
		router.AddRoute("GET", "/correlation", handlerWithCorrelationID)

		req, err := http.NewRequest("GET", "/correlation", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
		}

		// Check the response body (correlation ID)
		correlationID := rr.Body.String()
		if correlationID == "" {
			t.Errorf("Expected non-empty correlation ID, but got an empty string")
		}
	})

	t.Run("Query Parameters", func(t *testing.T) {
		// Test handler with query parameters
		handlerWithQueryParams := func(w http.ResponseWriter, req *http.Request) {
			queryParams := router.GetQueryParams(req)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(queryParams.Get("name")))
		}

		// Add route with query parameters
		router.AddRoute("GET", "/query", handlerWithQueryParams)

		req, err := http.NewRequest("GET", "/query?name=John", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
		}

		// Check the response body (query parameter)
		expectedQueryParam := "John"
		if rr.Body.String() != expectedQueryParam {
			t.Errorf("Expected response body %q, but got %q", expectedQueryParam, rr.Body.String())
		}
	})

	t.Run("Form Parameters", func(t *testing.T) {
		// Test handler with form parameters
		handlerWithFormParams := func(w http.ResponseWriter, req *http.Request) {
			formParams, err := router.GetFormParams(req)
			if err != nil {
				t.Errorf("Failed to parse form parameters: %v", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(formParams.Get("username")))
		}

		// Add route with form parameters
		router.AddRoute("POST", "/form", handlerWithFormParams)

		// Create a test request with form data
		req, err := http.NewRequest("POST", "/form", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.PostForm = map[string][]string{
			"username": {"john_doe"},
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the response status code
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
		}

		// Check the response body (form parameter)
		expectedFormParam := "john_doe"
		if rr.Body.String() != expectedFormParam {
			t.Errorf("Expected response body %q, but got %q", expectedFormParam, rr.Body.String())
		}
	})
}
