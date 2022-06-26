package mysql

import (
	"database/sql"
	"math/rand"

	"github.com/google/uuid"
)

type RouteModel struct {
	DB *sql.DB
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (m *RouteModel) Insert(endpoint_id int, name string) (int, error) {

	UUID := uuid.NewString()
	token := RandomString(8)

	stmt := `INSERT INTO routes (uuid, server_id, domain_id, endpoint_id, endpoint_type, name, spam_mode, created_at, updated_at, token, mode)
	VALUES(?, 1, 1, ?, "AddressEndpoint", ?, "Mark", UTC_TIMESTAMP(), UTC_TIMESTAMP(), ?, "Endpoint")`

	result, err := m.DB.Exec(stmt, UUID, endpoint_id, name, token)
	if err != nil {
		return 0, err
	}
	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}
