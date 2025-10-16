package main

import (
	"log"

	"backend-practice/internal/core/usecase"
	"backend-practice/internal/infra/config"
	"backend-practice/internal/infra/db"
	infraDB "backend-practice/internal/infra/db"
	infraHTTP "backend-practice/internal/infra/http"
	"backend-practice/internal/infra/jwt"
	"backend-practice/internal/infra/transport"
	"backend-practice/internal/infra/transport/handler"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	sqlDB, err := infraDB.ConnectMySQL(cfg.MySQLDSN, 10)
	if err != nil {
		log.Fatalf("could not connect to MySQL: %v", err)
	}
	defer sqlDB.Close()

	userRepo := db.NewUserRepository(sqlDB)
	jwtService := jwt.NewJWTService(cfg.JWTSecret)

	createUserUC := *usecase.NewCreateUserUseCase(userRepo)
	loginUseCase := *usecase.NewLoginUseCase(userRepo, jwtService)

	healthHandler := handler.NewHealthHandler()
	userHandler := handler.NewUserHandler(createUserUC)
	loginHandler := handler.NewLoginHandler(loginUseCase)

	router := transport.NewRouter(healthHandler, userHandler, loginHandler)

	infraHTTP.Run(":"+cfg.AppPort, router)
}
