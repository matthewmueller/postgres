package postgres

import (
	"database/sql"
	"os"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

type film struct {
	name  string
	genre string
}

func TestMain(m *testing.M) {
	log.SetHandler(text.New(os.Stderr))
	setup()
	code := m.Run()
	if code == 0 {
		teardown()
	}
	os.Exit(code)
}

func setup() {
	conn := "postgres://localhost:5432/test-postgres?sslmode=disable"
	client, err := Postgres(conn)
	if err != nil {
		log.WithError(err).Fatalf("Couldnt connect to %s", conn)
	}
	defer client.Close()

	_, qerr := client.Raw(`
    drop table if exists films;
    create table films (
      name text,
      genre text,
      created timestamp with time zone default now(),
      actors text[]
    );
    insert into films ("name", "genre", "actors") values ('Vanilla Sky', 'Drama', '{"Tom Cruise","Penelope Cruz"}');
  `)

	if qerr != nil {
		log.WithError(qerr).Fatalf("Couldnt setup the table")
	}
}

func teardown() {
	conn := "postgres://localhost:5432/test-postgres?sslmode=disable"
	client, err := Postgres(conn)
	if err != nil {
		log.WithError(err).Fatalf("Couldnt connect to %s", conn)
	}
	defer client.Close()

	_, qerr := client.Raw(`drop table if exists films;`)

	if qerr != nil {
		log.WithError(qerr).Fatalf("Couldnt teardown the table")
	}
}

func TestConnect(t *testing.T) {
	client, err := Postgres("postgres://localhost:5432/test-postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	client.Close()
}

func TestAll(t *testing.T) {
	client, err := Postgres("postgres://localhost:5432/test-postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := client.All("select name, genre from films where name = :name", map[string]interface{}{"name": "Vanilla Sky"})
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()
	defer rows.Close()

	if err == sql.ErrNoRows {
		t.Fatalf("no rows found")
	}

	var films []film

	for rows.Next() {
		var f film
		if err := rows.Scan(&f.name, &f.genre); err != nil {
			t.Fatal(err)
		}
		films = append(films, f)
	}

	if films[0].name != "Vanilla Sky" {
		t.Fatalf("wrong name %s", films[0].name)
	}

	if films[0].genre != "Drama" {
		t.Fatalf("wrong drama %s", films[0].genre)
	}
}

func TestOne(t *testing.T) {
	client, err := Postgres("postgres://localhost:5432/test-postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	row := client.One("select name, genre from films where name = :name", map[string]interface{}{"name": "Vanilla Sky"})

	var f film
	if err := row.Scan(&f.name, &f.genre); err != nil {
		t.Fatal(err)
	}

	if f.name != "Vanilla Sky" {
		t.Fatalf("wrong name %s", f.name)
	}

	if f.genre != "Drama" {
		t.Fatalf("wrong drama %s", f.genre)
	}
}

func TestMultiple(t *testing.T) {
	client, err := Postgres("postgres://localhost:5432/test-postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	row := client.One("select name, genre from films where name = :name and genre = :genre", map[string]interface{}{"genre": "Drama", "name": "Vanilla Sky"})

	var f film
	if err := row.Scan(&f.name, &f.genre); err != nil {
		t.Fatal(err)
	}

	if f.name != "Vanilla Sky" {
		t.Fatalf("wrong name %s", f.name)
	}

	if f.genre != "Drama" {
		t.Fatalf("wrong drama %s", f.genre)
	}
}

func TestInaccessble(t *testing.T) {
	client, _ := Postgres("postgres://localhost:5432/test-postgres-noooo?sslmode=disable")
	err := client.Ping()
	if err == nil {
		t.Fatalf("Should have been unable to connect")
	}
}
