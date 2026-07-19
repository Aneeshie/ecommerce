# HTTP Extensions (httpx)

## Overview

The `internal/httpx` module is a small, specialized utility package that provides consistent JSON serialization and robust error mapping across the entire API.

By centralizing how HTTP responses are formed, we guarantee that API consumers receive a predictable, unified response structure regardless of which module (Identity, Product, Order) they are interacting with.

---

## JSON Responses (`httpx.json.go`)

This file exposes generic, type-safe helpers for writing JSON payloads and decoding JSON requests.

### Standard Error Payload

Whenever an error occurs, the API never sends raw strings. Instead, it serializes a standard JSON error payload:

```json
{
  "error": "The requested product could not be found."
}
```

This is enforced by `httpx.WriteError(w http.ResponseWriter, status int, message string)`.

---

## Error Translation (`httpx.errors.go`)

One of the strict rules in Domain-Driven Design (DDD) is that the **Domain** and **Service** layers must never know anything about HTTP status codes. They return purely logical domain errors (e.g., `ErrNoRows`, `ErrInsufficientInventory`).

The `httpx.HandleError` function acts as an **Anti-Corruption Layer** between the domain and the HTTP transport.

### Dynamic Mapping

When a Handler receives an error from a Service, it simply calls:
```go
httpx.HandleError(w, err)
```

`httpx` intelligently maps the domain error to the correct HTTP status code:
- `ErrNoRows` (from the repository) → `404 Not Found`
- `ErrInvalidCredentials` (from Identity) → `401 Unauthorized`
- `ErrInsufficientInventory` (from Inventory) → `400 Bad Request`
- Unknown panics or database connection failures → `500 Internal Server Error`

### Benefits
- **Clean Handlers**: Handlers are kept extremely thin. They don't have large `switch` statements checking error types.
- **Pure Domain**: The domain models remain completely agnostic of the web transport layer.
- **Consistent Security**: Raw database errors or stack traces are automatically sanitized into generic `500` errors, preventing accidental data leakage to the client.
