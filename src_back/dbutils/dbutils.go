package dbutils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, nil
}

func GetAuthTokenByID(db *sql.DB, id int) (string, error) {
	var registrationDate string
	var email string
	err := db.QueryRow("SELECT registration_date, email FROM users WHERE id = $1", id).Scan(&registrationDate, &email)
	if err != nil {
		return "", err
	}
	result := registrationDate + " | " + email
	return result, nil
}
