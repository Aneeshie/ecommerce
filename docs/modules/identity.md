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