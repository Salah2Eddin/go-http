# Go HTTP Server

A lightweight, experimental HTTP 1.0 server implementation written in Go with a focus on simplicity and extensibility. This project demonstrates fundamental HTTP server concepts using only Go's standard library.

## Features

- **HTTP 1.0 Protocol Support**: Full implementation of HTTP 1.0 specification
- **Concurrent Request Handling**: Each connection is handled in a separate goroutine
- **Flexible Routing System**: Support for exact routes and wildcard patterns
- **URI Parameter Extraction**: Built-in support for extracting URI parameters
- **Custom Request/Response Handling**: Comprehensive request parsing and response building
- **Error Handling**: Structured error handling with appropriate HTTP status codes
- **Zero External Dependencies**: Built entirely with Go's standard library

## Requirements

This project relies solely on Go's standard library and does not require any external dependencies.

- Go 1.24.0 or later

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/Salah2Eddin/go-http.git
   cd go-http
   ```

2. Build the project:
   ```sh
   go build -o server example.go
   ```

3. Run the example server:
   - **Windows:**
     ```sh
     .\server.exe
     ```
   - **Linux/macOS:**
     ```sh
     ./server
     ```

## Usage

### Basic Server Setup

To create a server, follow these steps:

1. Import the necessary packages
2. Create handler functions that accept a `request.Request` and return a `response.Response`
3. Instantiate a `server.Address` object (optional - leave nil for automatic IP and port selection)
4. Create a server using `server.NewServer()`
5. Add route handlers using `Server.AddHandler(uri, method, handler)`
6. Start the server with `Server.Start()`

### Basic Example

```go
package main

import (
    "github.com/Salah2Eddin/go-http/pkg/request"
    "github.com/Salah2Eddin/go-http/pkg/response"
    "github.com/Salah2Eddin/go-http/pkg/response/statuscodes"
    "github.com/Salah2Eddin/go-http/pkg/server"
)

func helloHandler(req request.Request) response.Response {
    status := statuscodes.Status200()
    headers := response.NewResponseHeaders()
    headers.Add("content-type", "text/html")
    body := []byte("<h1>Hello, World!</h1>")
    return response.NewResponse(status, headers, &body)
}

func main() {
    // Create server with specific address
    app := server.NewServer(&server.Address{IP: "127.0.0.1", Port: "8008"})
    
    // Or use automatic address assignment
    // app := server.NewServer(nil)
    
    app.AddHandler("/", "GET", helloHandler)
    app.Start()
}
```

### Advanced Example with URI Parameters

```go
package main

import (
    "fmt"
    "github.com/Salah2Eddin/go-http/pkg/request"
    "github.com/Salah2Eddin/go-http/pkg/response"
    "github.com/Salah2Eddin/go-http/pkg/response/statuscodes"
    "github.com/Salah2Eddin/go-http/pkg/server"
)

func userHandler(req request.Request) response.Response {
    status := statuscodes.Status200()
    headers := response.NewResponseHeaders()
    headers.Add("content-type", "text/html")
    
    // Extract ID from URI segments (wildcard route)
    segments := req.Uri().GetSegments()
    id := segments[len(segments)-1] // Last segment is the ID
    
    var body []byte
    // Check for URI parameters (query string)
    if name, exists := req.GetUriParameter("name"); exists {
        body = []byte(fmt.Sprintf("<h1>Hello, %s!</h1><p>Your ID is: %s</p>", name, id))
    } else {
        body = []byte(fmt.Sprintf("<h1>User ID: %s</h1>", id))
    }
    
    return response.NewResponse(status, headers, &body)
}

func main() {
    app := server.NewServer(&server.Address{IP: "127.0.0.1", Port: "8008"})
    
    // Wildcard route - matches /user/123, /user/abc, etc.
    app.AddHandler("/user/*", "GET", userHandler)
    
    app.Start()
}
```

## Project Structure

```
go-http/
├── pkg/
│   ├── pkgerrors/      # Custom error types
│   ├── request/        # HTTP request parsing and handling
│   │   └── parsers/    # Request parsing utilities
│   ├── response/       # HTTP response building
│   │   └── statuscodes/ # HTTP status code definitions
│   ├── router/         # URL routing and handler management
│   ├── server/         # Core server implementation
│   ├── uri/            # URI parsing and validation
│   └── util/           # General utilities
├── example.go          # Example server implementation
├── go.mod             # Go module definition
└── README.md          # This file
```

## Routing

The server supports both exact and wildcard routing:

- **Exact routes**: `/api/users` matches only that exact path
- **Wildcard routes**: `/api/users/*` matches `/api/users/123`, `/api/users/john`, etc.

### Supported HTTP Methods

You can register handlers for any HTTP method (GET, POST, PUT, DELETE, etc.) by specifying the method string when calling `AddHandler()`.

### URI Parameters

Extract query parameters from the request using:
```go
if value, exists := request.GetUriParameter("param_name"); exists {
    // Use the parameter value
}
```

## Error Handling

The server automatically handles various error conditions:

- **400 Bad Request**: For malformed requests or invalid headers
- **404 Not Found**: When no route matches the requested URI
- **500 Internal Server Error**: For unexpected server errors

## Testing

You can test the server using curl or any HTTP client:

```sh
# Basic request
curl http://localhost:8008/

# Request with query parameters
curl "http://localhost:8008/user/123?name=John"
```

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

## License

This project is open-source and available under the MIT License.

---

For issues and contributions, feel free to open a pull request or report an issue in the repository.
