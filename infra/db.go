package infra

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

const (
	Driver = "postgres"
	DSN    = "postgres://$DB_USER_NAME:$DB_PASSWORD@$DB_HOST/$DB_DATABASE?sslmode=disable"
)

func NewDB() (*sql.DB, func(), error) {
	db, err := sql.Open(Driver, os.ExpandEnv(DSN))
	if err != nil {
		return nil, nil, fmt.Errorf("sql.Open: %w", err)
	}

	closeDB := func() {
		errC := db.Close()
		if errC != nil {
			log.Println(errC)
		}
	}

	err = db.Ping()
	if err != nil {
		closeDB()
		return nil, nil, fmt.Errorf("db.Ping: %w", err)
	}

	return db, closeDB, nil
}
