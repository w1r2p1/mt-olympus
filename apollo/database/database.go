package database

import (
	"apollo/env"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const dbType = "postgres"

var db *sql.DB

func Init(sslMode string) (*sql.DB, error) {
	var err error

	db = nil

	connStr := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v sslmode=%v",
		env.PostgresHost,
		env.PostgresUser,
		env.PostgresPass,
		env.PostgresDB,
		sslMode,
	)

	db, err = sql.Open(dbType, connStr)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(0)

	for i := uint(0); i < env.RetryTimes; i++ {
		log.Println("database retry #", i+1)

		if err := db.Ping(); err != nil {
			time.Sleep(time.Duration(env.RetrySeconds) * time.Second)
			continue
		}

		return db, nil
	}

	return db, nil
}

func GetDB() *sql.DB {
	return db
}
