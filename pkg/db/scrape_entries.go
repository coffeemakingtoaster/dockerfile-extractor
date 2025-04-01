package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func AddToDB(db *sql.DB, hash, url string) error {
	stmt, err := db.Prepare("INSERT INTO retrieved_repositories(hash, url) VALUES(?, ?)")
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

func GetPresentHashCount() (int, error) {
	db := RetrieveDbConn()

	defer db.Close()

	stmt, err := db.Prepare("SELECT count(*) FROM retrieved_repositories")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var presentItems int
	err = stmt.QueryRow().Scan(&presentItems)
	if err != nil {
		return 0, err
	}

	return presentItems, nil

}
