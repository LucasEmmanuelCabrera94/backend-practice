package db

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/core/port"
	"database/sql"
	"time"
)

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) port.SessionPort {
	return &sessionRepository{db: db}
}

func (s *sessionRepository) CreateSession(userId int64, token string) (entity.Session, error) {
	now := time.Now()
	expiration := now.Add(24 * time.Hour)

	query := `INSERT INTO sessions (user_id, token, created_at, expires_at)
			VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(query, userId, token, now, expiration)
	if err != nil {
		return entity.Session{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entity.Session{}, err
	}

	return entity.Session{
		ID:        id,
		UserID:    userId,
		Token:     token,
		CreatedAt: now,
		ExpiresAt: expiration,
	}, nil
}
