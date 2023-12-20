package src

import "time"

type User struct {
	Id               int64
	Name             string
	Email            string
	Password         string
	RegistrationDate time.Time
}

func NewUser(id int64, name string) User {
	return User{
		Id:   id,
		Name: name,
	}
}
func NewUserFullInfo(id int64, name, email, password string, registrationDate time.Time) User {
	return User{
		Id:               id,
		Name:             name,
		Email:            email,
		Password:         password,
		RegistrationDate: registrationDate,
	}
}
