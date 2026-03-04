## 1. Project Setup

- [x] 1.1 Initialize Go module with `go mod init` and create `go.mod`
- [x] 1.2 Create `main.go` with package declaration and empty `main()` function

## 2. CLI Flag Parsing

- [x] 2.1 Add `-p` and `-port` flags using the `flag` package, both pointing to the same variable with default value 7890
- [x] 2.2 Parse flags in `main()` and construct the listen address as `0.0.0.0:<port>`

## 3. HTTP Server

- [x] 3.1 Create an HTTP handler that responds with status 200, `Content-Type: text/plain`, and body `hello world`
- [x] 3.2 Configure `http.Server` with the listen address and handler
- [x] 3.3 Start the server in a goroutine and log the listen address to stdout

## 4. Graceful Shutdown

- [x] 4.1 Set up `os/signal` to capture SIGINT and SIGTERM
- [x] 4.2 On signal, call `server.Shutdown()` with a 5-second timeout context
- [x] 4.3 Log shutdown and exit cleanly with code 0

## 5. Verification

- [x] 5.1 `go build ./...` compiles without errors
- [x] 5.2 `go vet ./...` passes
- [x] 5.3 Server starts on default port 7890 and responds to `curl http://localhost:7890` with `hello world`
- [x] 5.4 `-p` and `-port` flags override the default port
