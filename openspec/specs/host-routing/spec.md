## ADDED Requirements

### Requirement: Route requests by Host header
The proxy SHALL inspect the `Host` header of incoming HTTP requests and route to the configured upstream if the host matches `<host>.localhost`. The port portion of the `Host` header (if present) SHALL be stripped before matching.

#### Scenario: Matching host
- **WHEN** a request arrives with `Host: myapp.localhost` or `Host: myapp.localhost:7890`
- **THEN** the proxy SHALL forward the request to `http://127.0.0.1:<app-port>`

#### Scenario: Non-matching host
- **WHEN** a request arrives with a `Host` header that does not match `<host>.localhost`
- **THEN** the proxy SHALL respond with HTTP 404 and a plain-text body indicating no route was found

### Requirement: Forward requests using reverse proxy
The proxy SHALL use `net/http/httputil.ReverseProxy` to forward matching requests to the upstream. The original request path, query string, and body SHALL be preserved.

#### Scenario: Request path preserved
- **WHEN** a request to `myapp.localhost:7890/api/users?page=2` matches the route
- **THEN** the proxy SHALL forward to `http://127.0.0.1:<app-port>/api/users?page=2`

#### Scenario: Request method preserved
- **WHEN** a POST request with a body to `myapp.localhost:7890/submit` matches the route
- **THEN** the proxy SHALL forward as a POST with the same body to the upstream

### Requirement: Set X-Forwarded headers
The proxy SHALL set `X-Forwarded-For`, `X-Forwarded-Host`, and `X-Forwarded-Proto` headers on proxied requests.

#### Scenario: Forwarded headers added
- **WHEN** a request from `192.168.1.10` to `myapp.localhost:7890` is proxied
- **THEN** the upstream request SHALL include `X-Forwarded-For: 192.168.1.10`, `X-Forwarded-Host: myapp.localhost`, and `X-Forwarded-Proto: http`

### Requirement: Return 502 when upstream is unreachable
When the upstream app is not listening on `<app-port>`, the proxy SHALL return HTTP 502 Bad Gateway.

#### Scenario: Upstream not listening
- **WHEN** a request matches the route but the upstream at `127.0.0.1:<app-port>` refuses the connection
- **THEN** the proxy SHALL respond with HTTP 502 and a plain-text body indicating the upstream is unreachable
