# Middleware & Authentication

## Overview

The `internal/middleware` package is responsible for cross-cutting HTTP concerns, specifically Authentication and Role-Based Access Control (RBAC). 

By isolating these concerns into middleware, the domain handlers (e.g., `product/handler`) remain completely decoupled from token parsing logic.

---

## Authentication (`AuthMiddleware`)

The `AuthMiddleware` intercepts incoming HTTP requests and ensures the client possesses a valid JWT Access Token.

### Flow

1. **Extraction**: Reads the `Authorization` header and extracts the Bearer token.
2. **Verification**: Uses the injected `token.Manager` to cryptographically verify the JWT signature and expiration.
3. **Context Injection**: If valid, the custom claims (containing the User ID and Role) are injected directly into the `http.Request` Context using a strongly-typed private context key.

```go
// Extracted claims are passed down the request chain
ctx := context.WithValue(r.Context(), claimsContextKey, customClaims)
next.ServeHTTP(w, r.WithContext(ctx))
```

---

## Role-Based Access Control (`RequireRole`)

Certain endpoints (like creating a product) must be strictly restricted to Administrators. The `RequireRole` middleware enforces this at the routing layer.

### Usage in Routing

```go
// Protect the POST /products route
r.With(
    authMiddleware.Auth, 
    authMiddleware.RequireRole(domain.RoleAdmin),
).Post("/", h.CreateProduct)
```

### Flow

1. **Context Extraction**: It extracts the `CustomClaims` previously injected by the `AuthMiddleware`.
2. **Authorization**: It compares the user's role against the required role.
3. **Rejection**: If the user is a `Customer` but the route requires `Admin`, the request is immediately terminated with a `403 Forbidden` response.

---

## Context Management (`context.go`)

To avoid "stringly-typed" keys and potential context collisions, the package defines a private struct type for context keys:

```go
type contextKey struct {
	name string
}

var claimsContextKey = &contextKey{"claims"}
```

It exposes a public helper `ClaimsFromContext(ctx)` which handlers and downstream services can use to safely retrieve the authenticated user's ID without needing to parse tokens themselves.
