package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

type Character struct {
	ID     int32 `gorm:"column:ID"`
	Name   string
	RaceID int32 `gorm:"column:raceID"`
	Race   Race  `gorm:"foreignKey:raceID"`
}

func (Character) TableName() string {
	return "Character"
}

type Race struct {
	ID   int32 `gorm:"column:ID"`
	Name string
}

func (Race) TableName() string {
	return "Race"
}

func InitializeDB() (*gorm.DB, error) {
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

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s", host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	fmt.Println("Connected to postgres db")

	return db, nil
}
