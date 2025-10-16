package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
    secretKey string
}

func NewJWTService(secretKey string) *jwtService {
    return &jwtService{secretKey: secretKey}
}

func (s *jwtService) GenerateToken(userID int64) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 1).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenString string) (int64, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.secretKey), nil
    })
    if err != nil || !token.Valid {
        return 0, err
    }

    claims := token.Claims.(jwt.MapClaims)
    userID := int64(claims["user_id"].(float64))
    return userID, nil
}
