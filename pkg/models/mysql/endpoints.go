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

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
