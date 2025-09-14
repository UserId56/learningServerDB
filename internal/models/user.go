package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserRegistrationData struct {
	Username   string `json:"user_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	BIO        string `json:"BIO"`
	AvatarURL  string `json:"avatar_url"`
}

func (u UserRegistrationData) GetHashPassword() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (u UserRegistrationData) Validate() error {
	switch {
	case u.Username == "":
		return fmt.Errorf("user_name обязательное поле")
	case u.Email == "":
		return fmt.Errorf("email обязательное поле")
	case u.Password == "":
		return fmt.Errorf("password обязательное поле")
	case len(u.Password) < 6:
		return fmt.Errorf("password должен быть не менее 6 символов")
	}
	return nil
}

type User struct {
	Id        int64  `json:"id" db:"id"`
	Username  string `json:"user_name" db:"username"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"passwordhash"`
	CreatedAt string `json:"created_at" db:"createdat"`
	UpdatedAt string `json:"updated_at" db:"updatedat"`
}
