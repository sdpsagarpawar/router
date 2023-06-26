# Router

This package provides a basic implementation of an HTTP router in Go. The router allows you to define routes and associate them with handler functions to handle incoming HTTP requests.

## Usage

To use the `router` package, follow these steps:

1. Import the package in your Go file:

   ```go
   import "github.com/sdpsagarpawar/router"
   ```

1. Create a new instance of the router:
```
router := router.NewRouter()
```
2. Add routes to the router using the AddRoute method:
```
router.AddRoute("GET", "/hello", helloHandler)
```
3. Set up your handler functions to handle the requests:
```
helloHandler := func(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello, World!"))
}
```
4. Run the router and send HTTP requests to test the routes:
```
req, err := http.NewRequest("GET", "/hello", nil)
if err != nil {
    log.Fatal(err)
}

rr := httptest.NewRecorder()
router.ServeHTTP(rr, req)
```
5. Assert and validate the responses:
```
if rr.Code != http.StatusOK {
    log.Fatalf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
}

expectedBody := "Hello, World!"
if rr.Body.String() != expectedBody {
    log.Fatalf("Expected response body %q, but got %q", expectedBody, rr.Body.String())
}
```

## Features
Support for handling different HTTP methods (GET, POST, PUT, DELETE, etc.)
Routing based on URL paths and query parameters
Customizable not found (404) handler
Middleware support for intercepting and modifying requests
Parsing and retrieval of query parameters and form data

## Testing
The router package provides a set of test cases to demonstrate the usage and functionality of the router. To run the tests, execute the following command:
```
go test -v github.com/sdpsagarpawar/router
```

## License
This package is distributed under the MIT License. See the LICENSE file for more information.

Feel free to modify and extend the package according to your needs. Contributions are always welcome!