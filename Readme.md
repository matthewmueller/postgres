# Postgres

Little Postgres wrapper that makes querying just a little bit easier.

## Example

```go
client, err := Postgres("postgres://localhost:5432/db?sslmode=disable")
rows, qerr := client.All("select name, genre from films where name = :name", map[string]interface{}{"name": "Inception"})
```

## Install

```sh
go get github.com/matthewmueller/postgres
```

## License

MIT
