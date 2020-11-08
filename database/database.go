package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	// required for database/sql to work with postgresql
	_ "github.com/lib/pq"
)

// Init starts up a postgresql database connection
func Init() (*sql.DB, error) {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		return nil, errors.New("error loading database .env file")
	}
	var (
		host     = os.Getenv("PGHOST")
		port, _  = strconv.Atoi(os.Getenv("PGPORT"))
		user     = os.Getenv("PGUSER")
		dbname   = os.Getenv("PGDATABASE")
		password = os.Getenv("PGPASSWORD")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to db")
	return db, nil
}
