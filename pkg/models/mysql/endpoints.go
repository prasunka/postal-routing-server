package mysql

import (
	"database/sql"

	"github.com/google/uuid"
)

type EndpointModel struct {
	DB *sql.DB
}

func (m *EndpointModel) Insert(address string) (int, error) {

	UUID := uuid.NewString()

	stmt := `INSERT INTO address_endpoints (server_id, uuid, address, last_used_at, created_at, updated_at)
	VALUES(1, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP(),UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, UUID, address)
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
