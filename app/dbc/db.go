package dbc

import (
	"database/sql"
	"log"
	"microStackServer/app/Config"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	conf := Config.Conf
	connectionString := conf.DB_USER + ":" + conf.DB_PASS + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME
	log.Println(conf)
	log.Println(connectionString)
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Println(err)
		log.Fatal("Could not connect to database")
	} else {
		log.Println("DB CONNECTION")
	}

	return db
}
