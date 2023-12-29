package forms

type RegistrationForm struct {
	Name             string
	Email            string
	Password         string
	RepeatedPassword string
	ErrorLine        string
}

func NewRegistrationForm(name, email, password, repeatedPassword string) *RegistrationForm {
	return &RegistrationForm{
		Name:             name,
		Email:            email,
		Password:         password,
		RepeatedPassword: repeatedPassword,
	}
}
