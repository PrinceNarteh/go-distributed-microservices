package db

import (
	"authentication/internals/config"
	"database/sql"
	"log"
	"time"
)

func ConnectDB() *sql.DB {
	dsn := config.Env.DB.DNS
	for range 10 {
		db, err := sql.Open("pgx", dsn)
		if err == nil && db.Ping() == nil {
			log.Println("connected to database successfully.")
			return db
		}
	}

	log.Println("waiting for Postgres")
	time.Sleep(2 * time.Second)
	return nil
}
