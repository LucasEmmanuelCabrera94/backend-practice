package db

import (
	"database/sql"
	"log"
	"time"
)

func ConnectMySQL(dsn string, maxAttempts int) (*sql.DB, error) {
	var sqlDB *sql.DB
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		dbConn, err := sql.Open("mysql", dsn)
		if err != nil {
			lastErr = err
		} else {
			lastErr = dbConn.Ping()
			if lastErr == nil {
				sqlDB = dbConn
				break
			}
		}
		log.Printf("mysql connect attempt %d/%d failed: %v", attempt, maxAttempts, lastErr)
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return sqlDB, nil
}
