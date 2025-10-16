package port

import (
	"backend-practice/internal/core/entity"
)

type SessionPort interface {
	CreateSession(userId int64, token string) (entity.Session, error)
}
