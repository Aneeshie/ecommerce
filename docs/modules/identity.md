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