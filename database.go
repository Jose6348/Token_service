package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func SetupDatabase() *sql.DB {
	connStr := "user=youruser dbname=yourdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Conexão com o banco de dados estabelecida")
	return db
}
