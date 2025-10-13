package db_test

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/infra/db"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	database, err := sql.Open("sqlite", ":memory:")

	assert.NoError(t, err)

	_, err = database.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT
	)`)
	assert.NoError(t, err)

	return database
}

func TestAddUser(t *testing.T) {
	dbConn := setupTestDB(t)
	repo := db.NewUserRepository(dbConn)

	user := entity.User{Name: "Lucas", Email: "lucas@example.com"}
	createdUser, err := repo.AddUser(user)

	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)
	assert.Equal(t, "Lucas", createdUser.Name)
}
