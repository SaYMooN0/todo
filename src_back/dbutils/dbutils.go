package dbutils

import (
	"database/sql"
	"errors"
	"fmt"
	src "my-todo-app/src_back"
	"time"

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

func GetAuthTokenByID(id int) (string, error) {
	var registrationDate time.Time
	var email string
	err := db.QueryRow("SELECT registration_date, email FROM users WHERE id = $1", id).Scan(&registrationDate, &email)
	if err != nil {
		return "", err
	}

	// yyyymmddhhmmss format
	formattedDate := registrationDate.Format("20060102150405")
	result := formattedDate + " | " + email
	return result, nil
}
func GetIDByEmail(email string) (int, error) {
	var id int
	query := `SELECT id FROM users WHERE email = $1;`
	err := db.QueryRow(query, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("no user with the given email found")
		}
		return 0, err
	}
	return id, nil
}
func AddNewUserToDb(user src.User) (int64, error) {
	query := `INSERT INTO users (name, email, password, registration_date) VALUES ($1, $2, $3, $4) RETURNING id;`
	var userId int64
	err := db.QueryRow(query, user.Name, user.Email, user.Password, user.RegistrationDate).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
