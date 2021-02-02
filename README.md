## Ledger - Auth

Auth service for the Ledger project. Handles user login and registration.

### Database Migrations

Install the [migrate CLI](https://github.com/golang-migrate/migrate).

#### Running Migrations

```shell
make migrate
```

#### Creating Migrations

```shell
migrate create -ext sql -dir sql/migrations -seq <migration name>
```
