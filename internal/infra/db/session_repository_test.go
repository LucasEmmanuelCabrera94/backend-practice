package db_test

import (
	"backend-practice/internal/infra/db"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func setupSessionTestDB(t *testing.T) *sql.DB {
	database, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	_, err = database.Exec(`
		CREATE TABLE sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			expires_at DATETIME NOT NULL
		)
	`)
	assert.NoError(t, err)

	return database
}

func TestCreateSession(t *testing.T) {
	dbConn := setupSessionTestDB(t)
	repo := db.NewSessionRepository(dbConn)

	userID := int64(1)
	token := "test-token"

	session, err := repo.CreateSession(userID, token)

	assert.NoError(t, err)
	assert.NotZero(t, session.ID)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, token, session.Token)
	assert.WithinDuration(t, time.Now(), session.CreatedAt, time.Second*2)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), session.ExpiresAt, time.Second*2)
}

func TestCreateMultipleSessions(t *testing.T) {
	dbConn := setupSessionTestDB(t)
	repo := db.NewSessionRepository(dbConn)

	userID := int64(1)

	s1, err := repo.CreateSession(userID, "token1")
	assert.NoError(t, err)

	s2, err := repo.CreateSession(userID, "token2")
	assert.NoError(t, err)

	assert.NotEqual(t, s1.ID, s2.ID)
	assert.NotEqual(t, s1.Token, s2.Token)
}
