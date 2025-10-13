package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	userrepo "backend-practice/internal/core/repository/user"
	"backend-practice/internal/core/usecase"
	"backend-practice/internal/infra/db"
	infrahttp "backend-practice/internal/infra/transport"
	"backend-practice/internal/infra/transport/handler"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load(".env")

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		user := getEnv("MYSQL_USER", "backend")
		pass := getEnv("MYSQL_PASSWORD", "backendpw")
		dbName := getEnv("MYSQL_DATABASE", "backend")
		host := getEnv("MYSQL_HOST", "")
		port := getEnv("MYSQL_PORT", "")

		if host == "" {
			if os.Getenv("DOCKER") == "1" {
				host = "mysql"
			} else {
				host = "127.0.0.1"
			}
		}

		if port == "" {
			if host == "127.0.0.1" {
				port = "3307"
			} else {
				port = "3306"
			}
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbName)
	}

	log.Printf("constructed MYSQL_DSN: %s", dsn)

	var sqlDB *sql.DB
	var repo userrepo.Repository
	var lastErr error
	maxAttempts := 10

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		dbConn, err := sql.Open("mysql", dsn)
		if err != nil {
			lastErr = err
		} else {
			lastErr = dbConn.Ping()
			if lastErr == nil {
				sqlDB = dbConn
				repo = db.NewUserRepository(sqlDB)
				break
			}
		}
		log.Printf("mysql connect attempt %d/%d failed: %v", attempt, maxAttempts, lastErr)
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	if lastErr != nil {
		log.Fatalf("could not connect to mysql after %d attempts: %v", maxAttempts, lastErr)
	}

	uc := usecase.NewCreateUserUseCase(repo)
	healthH := handler.NewHealthHandler()
	userH := handler.NewUserHandler(uc)
	router := infrahttp.NewRouter(healthH, userH)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", getEnv("APP_PORT", "8080")),
		Handler: router,
	}

	go func() {
		fmt.Printf("Server running on http://localhost:%s\n", getEnv("APP_PORT", "8080"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("error closing db: %v", err)
		}
	}
	log.Println("Server exiting")
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
