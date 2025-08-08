# SOCME Backend

## Rules

- Le premier utilisateur créé est un admin
- Lorsqu'un utilisateur est créé avec OAuth, nous lui avons donné le rôle "guest", il ne peut rien faire, il sera redirigé vers une page d'invité
- Un administrateur peut changer le rôle d'un utilisateur en "user" ou "admin".

- Les utilisateurs et les clients sont stockés dans des fichiers JSON.
- Les alertes sont stockées dans une base de données (SQLite pour dev, PostgreSQL pour prod)

## Auth

- Les utilisateurs peuvent s'authentifier en utilisant OAuth (GitHub)

## Paths

| **Category**  | **Method** | **Path**                          | **Description**                                      | **Permissions** |
|--------------|----------|--------------------------------|--------------------------------------------------|----------------|
| **Alerts**   | GET      | `/alerts`                     | Récupérer toutes les alertes (avec pagination et filtres) | User |
|              | GET      | `/alerts/:id`                 | Récupérer une alerte par ID                         | User |
|              | GET      | `/alerts/stats/:severity`     | Obtenir des infos sur le graphique (dernières 24 heures, 12 barres)   | User |
| **Clients**  | GET      | `/clients`                    | Récupérer tous les clients                                | User |
|              | GET      | `/clients/:id`                | Récupérer un client par ID                           | User |
|              | POST     | `/clients`                    | Créer un nouveau client                           | Admin |
| **Users**    | GET      | `/users`                      | Récupérer tous les users                                 | Admin |
|              | PATCH    | `/users/:id`             | Changer l'user (sauf son rôle)                              | Admin ou l'user |
|              | PATCH    | `/users/:id/role`             | Changer le rôle de l'user                              | Admin |
|              | DELETE   | `/users/:id/session`          | Revoquer un token de session                    | Admin ou l'user |
|              | DELETE   | `/users/:id`                  | Supprimer un user                                 | Admin ou l'user |
| **Misc**     | GET      | `/certfr`                     | Récupérer les alertes du CERT-FR                        | User |
| **Auth**     | GET      | `/auth/refresh`               | Rafraîchir le jeton et les infos de l'user                     | User |
| | GET      | `/auth/github`               | Auth avec Github | User |
| | GET      | `/auth/callback`               | Auth callback | User |

## Utilisation

Démarrez le backend :

```bash
go run ./cmd/main.go
```
