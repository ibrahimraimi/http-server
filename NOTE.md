# Code Explanation

This Go program creates a simple web server that responds to HTTP requests with JSON data. Let's break down the key parts:

## Import Statements

```go
import (
  "encoding/json"     // For working with JSON data
  "fmt"                // For printing output to the console
  "net"                // For network operations (getting IP addresses)
  "net/http"           // For building HTTP servers
  "os"                 // For interacting with the operating system (signals)
  "os/signal"          // For handling operating system signals
  "syscall"            // For system calls related to signals
)
```

These lines bring in the necessary packages from the Go standard library to handle JSON, network requests, system signals, and more.

## `main` Function

```go
func main() {
  // ...
}
```

The `main` function is the entry point of the program:

I. **HTTP Handler Setup:**

- `http.HandleFunc("/", handler)` tells the server to use the `handler` function whenever a request comes in for the root path `(/)`.

II. **Server Start:**

- `server := &http.Server{Addr: ":8000"}` creates an HTTP server that will listen on port `8000`.

- `go func() { ... }()` starts the server in a separate goroutine (like a lightweight thread) so it can run concurrently.

- `if err := server.ListenAndServe(); err != nil { ... }` starts listening for incoming requests and prints an error if it fails.

III. **Graceful Shutdown:**

- `stop := make(chan os.Signal, 1)` creates a channel to receive operating system signals.

- `signal.Notify(stop, os.Interrupt, syscall.SIGTERM)` tells the program to send interrupt (Ctrl+C) and termination signals to the stop channel.

- `<-stop` waits for a signal to be received.

- `fmt.Println("Shutting down server...")` prints a message.

- `server.Shutdown(nil)` initiates the shutdown process, allowing existing requests to complete.

## `handler` Function

```go
func handler(w http.ResponseWriter, r *http.Request) {
  // ...
}
```
I. **IP Extraction:**

- `ip, _,_ := net.SplitHostPort(r.RemoteAddr)` gets the client's IP address from the request.

- `_,_, _` is used to discard the port number and any errors.

II. **Response Construction:**

- `response := struct { ... } { ... }` creates an anonymous structure to hold the response data (IP, location, greeting).

III. **JSON Encoding:**

- `json.NewEncoder(w).Encode(response)` converts the response struct into JSON format and sends it as the response to the client.
