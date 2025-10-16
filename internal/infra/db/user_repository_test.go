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
		email TEXT,
		passwordhash TEXT
	)`)
	assert.NoError(t, err)

	return database
}

func TestAddUser(t *testing.T) {
	dbConn := setupTestDB(t)
	repo := db.NewUserRepository(dbConn)

	user := entity.User{Name: "Lucas", Email: "lucas@example.com", PasswordHash: "hashedpassword"}
	createdUser, err := repo.CreateUser(user)

	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)
	assert.Equal(t, "Lucas", createdUser.Name)
}

func TestGetUserByEmail(t *testing.T) {
	dbConn := setupTestDB(t)
	repo := db.NewUserRepository(dbConn)

	user := entity.User{Name: "Lucas", Email: "lucas@example.com", PasswordHash: "hashedpassword"}
	createdUser, err := repo.CreateUser(user)
	assert.NoError(t, err)

	fetchedUser, err := repo.GetUserByEmail("lucas@example.com")
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, fetchedUser.ID)
	assert.Equal(t, createdUser.Name, fetchedUser.Name)

	_, err = repo.GetUserByEmail("noone@example.com")
	assert.Error(t, err)
}
