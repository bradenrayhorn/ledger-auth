## Ledger - Auth

Auth service for the Ledger project. Handles user login and registration.

### Setup

Copy `config.example.yml` to `config.yml` and fill in the appropriate values.

### Creating Queries

The [sqlc](https://github.com/kyleconroy/sqlc) library is used to generate Go functions based
on provided queries.

Queries are placed in the `sql/query` folder and generate code is in the `internal/db` package.

To add a new query, add the query to the `sql/query` folder. Then run the `sqlc compile` command
to verify syntax. Finally, `sqlc generate` will generate the Go code.

### Testing

Tests are located in the tests package. These are integration tests that require a database connection.

There are two options to run tests:

Running in local environment:
```shell
go test ./...
```

Running with docker-compose:
```shell
make test
```

### Database Migrations

Install the [migrate CLI](https://github.com/golang-migrate/migrate).

#### Running Migrations

Make sure the connection is setup in Makefile.


Migrate up:

```shell
make migrate
```

Migrate down:

```shell
make migrate-down
```

#### Creating Migrations

```shell
migrate create -ext sql -dir sql/migrations -seq <migration name>
```

### Vault PKI

The following command is an example of how to create a Vault PKI role with proper restrictions:

```bash
write intermediateca/roles/ledger-auth allow_localhost=false allowed_domains=ledger-auth,ledger-auth.ledger-auth allow_bare_domains=true enforce_hostnames=false key_type=ec key_bits=256
```
