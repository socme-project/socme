# SOCME Backend

## Rules

- The first user created is an admin
- When a user is created with OAuth, we give them the "guest" role, they can't do anything, they will be redirected to a guest page
- An administrator can change a user's role to "user" or "admin".

- Users and clients are stored in JSON files.
- Alerts are stored in a database (SQLite for dev, PostgreSQL for prod)

## Auth

- Users can authenticate using OAuth (GitHub)

## Paths

| **Category**  | **Method** | **Path**                          | **Description**                                      | **Permissions** |
|--------------|----------|--------------------------------|--------------------------------------------------|----------------|
| **Alerts**   | GET      | `/alerts`                     | Retrieve all alerts (with pagination and filters) | User |
|              | GET      | `/alerts/:id`                 | Retrieve an alert by ID                         | User |
|              | GET      | `/alerts/stats/:severity`     | Get info on the graph (last 24 hours, 12 bars)   | User |
| **Clients**  | GET      | `/clients`                    | Retrieve all clients                                | User |
|              | GET      | `/clients/:id`                | Retrieve a client by ID                           | User |
|              | POST     | `/clients`                    | Create a new client                           | Admin |
| **Users**    | GET      | `/users`                      | Retrieve all users                                 | Admin |
|              | PATCH    | `/users/:id`             | Change the user (except their role)                              | Admin ou l'user |
|              | PATCH    | `/users/:id/role`             | Change the user's role                              | Admin |
|              | DELETE   | `/users/:id/session`          | Revoke a session token                    | Admin ou l'user |
|              | DELETE   | `/users/:id`                  | Delete a user                                 | Admin ou l'user |
| **Misc**     | GET      | `/certfr`                     | Retrieve alerts from CERT-FR                        | User |
| **Auth**     | GET      | `/auth/refresh`               | Refresh the token and user info                     | User |
| | GET      | `/auth/github`               | Auth with Github | User |
| | GET      | `/auth/callback`               | Auth callback | User |
