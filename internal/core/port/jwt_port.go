package port

type JWTService interface {
    GenerateToken(userID int64) (string, error)
    ValidateToken(token string) (int64, error)
}
