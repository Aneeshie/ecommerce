# Identity Module

## Responsibilities

- Register users
- Authenticate users
- Authorize users based on roles
- Manage user sessions
- Manage user profile
- Manage user addresses

## Business Rules

- Email must be unique.
- Passwords must be securely hashed.
- Users can have multiple addresses.
- A user can have only one default address.
- Users cannot delete their default address directly.
- Before deleting the default address, the user must choose another address as the default.
- Guests can browse the application without authentication.
- Only authenticated users can manage their profile.
- Only admins can perform admin-only actions.

## Owns

- Users
- Addresses
- Refresh Tokens

## Open Questions

- Should email verification be required before login?
- Should users log in using email only, or username too?
- How long should refresh tokens live?
- Should users be able to stay logged in on multiple devices?
- How are forgotten passwords handled?


## Domain Entities

### User
- id
- name
- email
- password_hash
- role
- email_verified
- created_at
- updated_at

### Address
- id
- user_id
- line1
- line2
- city
- state
- country
- postal_code
- is_default
- updated_at

### Refresh Tokens
- id
- user_id
- token_hash
- expires_at
- updated_at
- created_at


## Relationships
- One User can have many Addresses
- One User can have many Refresh Tokens.
- One Address belongs to exactly one User.
- One Refresh Token belongs to exactly one User.

## Constraints

### User
- Email must be unique
- Password is stored as hash
- role = 'Admin' | 'Customer'
### Address
- A user can have multiple Addresses
- Exactly one address is marked as default.

### Refresh Token
- Tokens must expire
- Tokens are stored as hashes

# Workflows

## User Registration

### Trigger
- Guest submits the registration form

### Flow
1. User enters name, email and password.
2. Validate the input.
3. Check if email already exists.
4. Hash the password
5. Create the user
6. Generate an email verification token.
7. Send the verification email.
8. Return success.


## Login

### Trigger
- User submits email and password.

## Flow
1. Validate the request.
2. Find the user by email.
3. Verify the password hash.
4. Ensure the account is active
5. Generate an access token.
6. Generate a refresh token.
7. Store the refresh token.
8. Return both tokens.

## Logout

### Trigger
- Authenticated user logs out.

## Flow
1. Receive the refresh tokens.
2. Validate the refresh tokens.
3. Remove teh refresh token.
4. return success.

## Add Address

### Trigger
- Authenticated user adds a new address.

## Flow
1. Validate the address.
2. if this is the user's first address, mark it as default.
3. otherwise, ask whether it should become the default.
4. Save the address.
5. return success.

## Delete Address

### Trigger
- Authenticated user deletes an address.

## Flow
1. Verify the address belongs to the user.
2. Check if it is the default address .
3. If it is the default, require the user to choose another default address.
4. Delete the address.
5. Return success.

# Database


## users

| Column         | Type      | Constraints      |
|----------------|-----------|------------------|
| id             | UUID      | PK               |
| name           | TEXT      | NOT NULL         |
| email          | TEXT      | UNIQUE, NOT NULL |
| password_hash  | TEXT      | NOT NULL         |
| role           | TEXT      | NOT NULL         |
| email_verified | BOOLEAN   | DEFAULT FALSE    |
| created_at     | TIMESTAMP | NOT NULL         |
| updated_at     | TIMESTAMP | NOT NULL         |

## addresses

| Column      | Type      | Constraints               |
|-------------|-----------|---------------------------|
| id          | UUID      | PK                        |
| user_id     | UUID      | FK -> users(id), NOT NULL |
| line_1      | TEXT      | NOT NULL                  |
| line_2      | TEXT      | NULL                      |
| label       | TEXT      | NOT NULL                  |
| city        | TEXT      | NOT NULL                  |
| state       | TEXT      | NOT NULL                  |
| country     | TEXT      | NOT NULL                  |
| postal_code | TEXT      | NOT NULL                  |
| is_default  | BOOLEAN   | NOT NULL DEFAULT FALSE    |
| created_at  | TIMESTAMP | NOT NULL                  |
| updated_at  | TIMESTAMP | NOT NULL                  |


## refresh_tokens
| Column     | Type      | Constraints               |
|------------|-----------|---------------------------|
| id         | UUID      | PK                        |
| user_id    | UUID      | FK -> users(id), NOT NULL |
| token_hash | TEXT      | NOT NULL                  |
| expires_at | TIMESTAMP | NOT NULL                  |
| updated_at | TIMESTAMP | NOT NULL                  |
| created_at | TIMESTAMP | NOT NULL                  |
| revoked_at | TIMESTAMP | NULL                      |


# API ENDPOINTS

POST /auth/register

Request
```json
    {
      "name": "",
      "email": "",
      "password": ""
    }
```


Response
```json
    {
      "message": "Registration successful. Please verify your email."
    }
```


GET /auth/verify-email?token=...

POST /auth/login

Request
```json
    {
      "email": "",
      "password": ""
    }
```

Response
```json
    {
      "access_token": "",
      "refresh_token": "",
      "expires_in": 900,
      "user": {
        ....
      }
    }
```

POST /auth/forgot-password

POST /auth/reset-password

POST /auth/logout

Request
```json
    {
      "refresh_token": ""
    }
```

POST /auth/refresh

Request
```json
    {
      "refresh_token":  ""
    }
```

Response
```json
    {
      "access_token": "...",
      "expires_in": 900
    }
```

GET /me

Response

```json
    {
      "id": "",
      "name": "",
      "email": "",
      "role": ""
    }
```

PATCH /me

GET /me/addresses

POST  /me/addresses

DELETE /me/addresses/{addressID}

PATCH /me/addresses/{addressID}


# Authentication Middleware Design

## Goal

Protect routes by verifying the user's JWT before the request reaches the handler.

---

## Responsibilities

- Read the `Authorization` header.
- Validate the `Bearer` token format.
- Verify the JWT signature.
- Verify the token has not expired.
- Extract the JWT claims.
- Store the claims in the request context.
- Pass the request to the next middleware/handler.

---

## Non-Responsibilities

The middleware should **not**:

- Query the database.
- Check whether the user is an admin.
- Perform business logic.
- Generate tokens.

These belong elsewhere.

---

## Request Flow

Incoming Request

↓

Read `Authorization` header

↓

Missing?

→ Return `401 Unauthorized`

↓

Validate `Bearer <token>` format

↓

Extract JWT

↓

Verify signature

↓

Verify expiration

↓

Extract claims

↓

Store claims in request context

↓

Call `next.ServeHTTP()`

---

## Claims

The middleware should extract:

- User ID (`sub`)
- Role
- Issued At (`iat`)
- Expiration (`exp`)

---

## Error Cases

| Condition | Response |
|-----------|----------|
| Missing Authorization header | 401 |
| Invalid Bearer format | 401 |
| Invalid signature | 401 |
| Expired token | 401 |
| Malformed JWT | 401 |

---

## Context

The middleware should attach the authenticated user's claims to the request context.

Handlers should retrieve authentication information from the context instead of parsing the JWT again.

---

## Public Routes

Authentication middleware should **not** protect:

- POST `/auth/register`
- POST `/auth/login`
- POST `/auth/refresh`

---

## Protected Routes

Examples:

- GET `/auth/me`
- POST `/products`
- PUT `/products/{id}`
- DELETE `/products/{id}`
- POST `/orders`
- GET `/cart`

---

