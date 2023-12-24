package dbutils

import (
	"database/sql"
	"errors"
	"fmt"
	structs "my-todo-app/src_back/structs"
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
	result := formattedDate + "|" + email
	return result, nil
}
func GetIDByEmail(email string) (int64, error) {
	var id int64
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

func AuthenticateUser(email, password string) (structs.User, error) {
	var user structs.User
	userId, err := GetIDByEmail(email)
	if err != nil {
		return user, err
	}
	query := `SELECT id, name, email, password, registration_date FROM users WHERE id = $1;`
	err = db.QueryRow(query, userId).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.RegistrationDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	if password != user.Password {
		return user, errors.New("invalid password")
	}
	return user, nil
}

func AddNewUserToDb(user structs.User) (int64, error) {
	query := `INSERT INTO users (name, email, password, registration_date) VALUES ($1, $2, $3, $4) RETURNING id;`
	var userId int64
	err := db.QueryRow(query, user.Name, user.Email, user.Password, user.RegistrationDate).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
func UpsertEmailConfirmation(code string, user structs.User) (int64, error) {
	var existingID int64
	queryCheck := `SELECT id FROM email_confirmation WHERE email = $1;`
	err := db.QueryRow(queryCheck, user.Email).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if existingID > 0 {
		queryUpdate := `UPDATE email_confirmation SET name = $1, password = $2, confirmation_code = $3 WHERE id = $4 RETURNING id;`
		err := db.QueryRow(queryUpdate, user.Name, user.Password, code, existingID).Scan(&existingID)
		if err != nil {
			return 0, err
		}
		return existingID, nil
	} else {
		queryInsert := `INSERT INTO email_confirmation (name, email, password, confirmation_code) VALUES ($1, $2, $3, $4) RETURNING id;`
		var newID int64
		err := db.QueryRow(queryInsert, user.Name, user.Email, user.Password, code).Scan(&newID)
		if err != nil {
			return 0, err
		}
		return newID, nil
	}
}
