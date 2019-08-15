# Inventario

Inventory Management API

## ENV Variable

| ENV var name    | Default       | Note                                  |
|-----------------|---------------|---------------------------------------|
| APP_ADDRESS     | :8080         |                                       |
| APP_DEBUG       | 1             | available value: `0` or `1`           |
| APP_NAME        | Inventario    |                                       |
| APP_LOG_ENGINE  | logrus        | available value: `logrus` or `stdlib` |
| DATABASE_ENGINE | sqlite3       | available value: `sqlite3`            |
| DATABASE_URL    | inventario.db |                                       |

## Requirements

You need this installed on your system to run locally
- [Go](https://golang.org/)
- [Dep](https://golang.github.io/dep/)

## Database Migration

Run migration to update your database schema

```bash
$ go build -tags migrate -o ./bin/migrate
$ ./bin/migrate up
```

or you can run `migrate.go` directly

```bash
$ go run migrate.go up
```

## How To Run

**NOTE:** Before you run this app, make sure your database has latest schema.
You can run [database migration](#database-migration) to update its schema

```bash
$ go build -tags main -o ./bin/inventario
$ ./bin/inventario
```

or you can run `main.go` directly

```bash
$ go run main.go
```