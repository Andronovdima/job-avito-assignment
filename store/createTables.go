package store

import (
	"database/sql"
	"io/ioutil"
)

func CreateTables(db *sql.DB) error {

	file, err := ioutil.ReadFile("store/main.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(file))
	if err != nil {
		return err
	}

	return nil
}
