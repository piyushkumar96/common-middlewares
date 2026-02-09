# common-middlewares

Reusable Gin HTTP middlewares: **trace**, **cors**, **authentication**, **context**, and **openapi** (request validation). Uses [piyushkumar96/app-error](https://github.com/piyushkumar96/app-error), [piyushkumar96/app-monitoring](https://github.com/piyushkumar96/app-monitoring) (optional), and [piyushkumar96/generic-logger](https://github.com/piyushkumar96/generic-logger).

## Layout

```
common-middlewares/
├── authentication/   # Static-token auth (Authorization header)
│   └── examples/
├── context/          # Request context, request ID, response helpers
│   └── examples/
├── cors/             # CORS with configurable headers and origin regex
│   └── examples/
├── openapi/          # OpenAPI/Swagger request validation (kin-openapi)
│   └── examples/
├── trace/            # Request context + trace/response meta (for monitoring)
│   └── examples/
├── go.mod
└── README.md
```

## Installation

```bash
go get github.com/piyushkumar96/common-middlewares/...
```

## Packages

| Package | Import | Description |
|--------|--------|-------------|
| **trace** | `github.com/piyushkumar96/common-middlewares/trace` | Initializes request context (context meta, response meta, trace meta); use with app-monitoring for metrics. |
| **cors** | `github.com/piyushkumar96/common-middlewares/cors` | CORS middleware with configurable headers and origin regex. |
| **authentication** | `github.com/piyushkumar96/common-middlewares/authentication` | Auth via static token in `Authorization` header. |
| **context** | `github.com/piyushkumar96/common-middlewares/context` | Request ID, `InitRequestContext`, `GetRequestContext`, `RespondJSON`, `MessageFailure`, context meta. |
| **openapi** | `github.com/piyushkumar96/common-middlewares/openapi` | OpenAPI request and optional response validation; `OpenAPIValidatorRequest` (request only), `OpenAPIValidatorRequestAndResponse` (request + response; response failures logged). |

## Examples

Each package has an `examples/` folder with a runnable `examples.go` (main package).

| Example | Command | Port |
|--------|--------|------|
| **trace** | `go run ./trace/examples` | 8080 |
| **cors** | `go run ./cors/examples` | 8081 |
| **authentication** | `go run ./authentication/examples` | 8082 |
| **context** | `go run ./context/examples` | 8083 |
| **openapi** | `go run ./openapi/examples` (optional: add `openapi.yaml` in that dir) | 8084 |

From repo root:

```bash
# Trace: request ID and context meta
go run ./trace/examples
# GET http://localhost:8080/ping

# CORS: configurable headers and origin regex
go run ./cors/examples
# GET http://localhost:8081/ping

# Auth: static token
go run ./authentication/examples
# curl -H "Authorization: my-secret-token" http://localhost:8082/ping

# Context: request ID and response helpers
go run ./context/examples
# GET http://localhost:8083/ping or /fail

# OpenAPI: validator (needs openapi.yaml / openapi.json in openapi/examples/ to enable)
go run ./openapi/examples
# GET http://localhost:8084/ping
```

## Quick usage

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/piyushkumar96/common-middlewares/authentication"
    "github.com/piyushkumar96/common-middlewares/context"
    "github.com/piyushkumar96/common-middlewares/cors"
    "github.com/piyushkumar96/common-middlewares/trace"
)

r := gin.New()
r.Use(trace.Trace(nil))  // or pass app-monitoring AppMetricsInterface
r.Use(cors.CORS(&cors.CORSHeaders{AccessControlAllowOrigin: ".*", ...}))
r.Use(authentication.Auth(&authentication.AuthConfig{Token: "your-token"}))

r.GET("/ping", func(c *gin.Context) {
    ctx := context.GetRequestContext(c)
    // ...
})
```

## License

See [LICENSE](./LICENSE).
