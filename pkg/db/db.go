package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func init() {
	db := RetrieveDbConn()
	defer db.Close()

	_, err := db.Exec(initStatement)
	if err != nil {
		panic(err)
	}

}

func RetrieveDbConn() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/sqlite.db")
	if err != nil {
		log.Error().Msgf("Could not connect to db: %s", err.Error())
		panic("Fatal")
	}
	return db
}
