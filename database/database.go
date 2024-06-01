package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Connection struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func ConnectToDatabase(conn Connection) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conn.Host, conn.Port, conn.User, conn.Password, conn.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, nil
}
