package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func AddToDB(db *sql.DB, hash, url string) error {
	stmt, err := db.Prepare("INSERT INTO retrieved_repositories(hash, repo) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(hash, url)
	if err != nil {
		return err
	}
	return nil
}

func CheckIfAlreadyPresent(hash string) (bool, error) {
	db := RetrieveDbConn()

	defer db.Close()

	stmt, err := db.Prepare("SELECT hash FROM retrieved_repositories WHERE hash = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var retrievedHash string
	err = stmt.QueryRow(hash).Scan(&retrievedHash)
	if err != nil {
		return false, err
	}

	return hash == retrievedHash, nil
}

func GetPresentHashCount() (int, int, error) {
	db := RetrieveDbConn()

	defer db.Close()

	stmt, err := db.Prepare("SELECT count(*), sum(scraped) FROM retrieved_repositories")
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	var presentItemsTotal int
	var presentItemsDone int
	err = stmt.QueryRow().Scan(&presentItemsTotal, &presentItemsDone)
	if err != nil {
		return 0, 0, err
	}

	return presentItemsTotal, presentItemsDone, nil
}

func SetScrapedByHash(hash string) error {
	db := RetrieveDbConn()
	defer db.Close()
	stmt, err := db.Prepare("UPDATE retrieved_repositories SET scraped = 1 WHERE hash = ?")

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hash)
	return err
}

func GetRepoIterator() *sql.Rows {
	db := RetrieveDbConn()
	rows, err := db.Query("select hash, repo from retrieved_repositories WHERE scraped=0")
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return rows
}
