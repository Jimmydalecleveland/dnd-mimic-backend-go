package pgconnect

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	// required for database/sql to work with postgresql
	_ "github.com/lib/pq"
)

// InitializeDB starts a database connection with an initial ping and returns that instance.
func InitializeDB() *sql.DB {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatalf("Error loading database .env file")
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
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to db")
	return db
}
