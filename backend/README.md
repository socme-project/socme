## Documentation

- [ ] Add OAuth2 via github

## Rules

- The first user created is an admin
- When a user is created using OAuth, we gave him the "guest" role, it can't do anything.
- An admin can change the role of a user to "user" or "admin"

## Path

| **Category**  | **Method** | **Path**                          | **Description**                                      | **Permissions** |
|--------------|----------|--------------------------------|--------------------------------------------------|----------------|
| **Alerts**   | GET      | `/alerts`                     | Get all alerts (supports pagination & filters)  | User |
|              | GET      | `/alerts/:id`                 | Get one alert by ID                            | User |
|              | GET      | `/alerts/stats/:severity`     | Get graph info (last 24h, 12 bars)             | User |
| **Clients**  | GET      | `/clients`                    | Get all clients                                | User |
|              | GET      | `/clients/:id`                | Get one client by ID                           | User |
|              | POST     | `/clients`                    | Create a new client                           | Admin |
| **Users**    | GET      | `/users`                      | Get all users                                 | Admin |
|              | PATCH    | `/users/:id/role`             | Change user role                              | Admin |
|              | DELETE   | `/users/:id/session`          | Revoke a session token                        | Admin or THE user |
|              | DELETE   | `/users/:id`                  | Delete a user                                 | Admin or THE user |
| **Misc**     | GET      | `/certfr`                     | Get the CERT-FR alerts                        | User |
| **Auth**     | GET      | `/auth/refresh`               | Refresh user token & info                     | User |
| | GET      | `/auth/github`               | Auth with github | User |
| | GET      | `/auth/callback`               | Auth callback | User |

| **Category**  | **Method** | **Path**                          | **Description**                                      | **Permissions** |
|--------------|----------|--------------------------------|--------------------------------------------------|----------------|
| **Alerts**   | GET      | `/alerts`                     | Get all alerts (supports pagination & filters)  | User |
|              | GET      | `/alerts/:id`                 | Get one alert by ID                            | User |
|              | GET      | `/alerts/stats/:severity`     | Get graph info (last 24h, 12 bars)             | User |
| **Clients**  | GET      | `/clients`                    | Get all clients                                | User |
|              | GET      | `/clients/:id`                | Get one client by ID                           | User |
|              | POST     | `/clients`                    | Create a new client                           | Admin |
| **Users**    | GET      | `/users`                      | Get all users                                 | Admin |
|              | PATCH    | `/users/:id/role`             | Change user role                              | Admin |
|              | DELETE   | `/users/:id/session`          | Revoke a session token                        | Admin or THE user |
|              | DELETE   | `/users/:id`                  | Delete a user                                 | Admin or THE user |
