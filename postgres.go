package postgres

import (
	"database/sql"
	"strconv"

	"github.com/matthewmueller/colon"

	// postgres side-effect
	_ "github.com/lib/pq"
)

// Client struct, extends DB struct
type Client struct {
	*sql.DB
}

// Connect to postgres and return a client
func Connect(conn string) (Client, error) {
	var client Client
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return client, err
	}

	err = db.Ping()
	if err != nil {
		return client, err
	}

	client = Client{db}
	return client, nil
}

// Prepare the query
func Prepare(query string, params map[string]interface{}) (string, []interface{}) {
	values := make([]interface{}, 0, len(params))
	render := colon.Compile(query)
	i := 0

	for key, value := range params {
		values = append(values, value)
		i = i + 1
		params[key] = "$" + strconv.Itoa(i)
	}

	return render(params), values
}

// All get all the results
func (c Client) All(query string, params map[string]interface{}) (*sql.Rows, error) {
	query, values := Prepare(query, params)
	return c.Query(query, values...)
}

// One get all the results
func (c Client) One(query string, params map[string]interface{}) *sql.Row {
	query, values := Prepare(query, params)
	return c.QueryRow(query, values...)
}

// Raw get all the results
func (c Client) Raw(query string) (*sql.Rows, error) {
	return c.Query(query)
}
