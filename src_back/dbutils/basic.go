package dbutils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB = nil

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
func CloseConnection() {
	db.Close()
}
