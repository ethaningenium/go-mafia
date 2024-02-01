package storage

import (
	"database/sql"
	"log"
	"mafia/services"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

var database *Database

func Start(){
	connStr := services.GetEnv("DSN")
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    log.Fatal(err)
  }
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	database = &Database{db}

	createTables(database)
}

func GetDatabase() *Database {
	return database
}

