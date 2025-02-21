# Go HTTP Server

A simple, experimental HTTP 1.0 server written in Go.  
Feel free to explore and contribute.

## Requirements

This project relies solely on Go's standard library and does not require any external modules.

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/Salah2Eddin/go-http-server.git
   ```

2. Build the module:
   ```sh
   go build
   ```

3. Run the server:  
   - **Windows:**
     ```sh
     .\http.exe
     ```
   - **Linux/macOS:**
     ```sh
     ./http
     ```

## Usage

The `main.go` file provides an example usage of the server.  
To create a server, follow these steps:

1. Instantiate a `server.ServerAddress` object (leave it empty for automatic IP and port selection).
2. Create a server using `server.NewServer()`.
3. Use `Server.AddHandler(uri, method, handler)` to add route handlers.

### Example:

```go
package main

import (
    "ducky/http/pkg/request"
    "ducky/http/pkg/response"
    "ducky/http/pkg/server"
    statuscodes "ducky/http/pkg/response/statuscodes"
)

func index(request *request.Request) *response.Response {
    status := statuscodes.Status200()

    headers := response.NewResponseHeaders()
    headers.Set("content-type", "text/html")

    body := []byte("<h1>Hello, World!</h1>")

    return response.NewResponse(status, headers, &body)
}

func main() {
    server := server.NewServer(&server.ServerAddress{Ip: "127.0.0.1", Port: "8008"})
    server.AddHandler("/", "GET", index)
    server.Start()
}
```

## License

This project is open-source and available under the MIT License.

---

For issues and contributions, feel free to open a pull request or report an issue in the repository.