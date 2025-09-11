package models

import "golang.org/x/crypto/bcrypt"

type UserRegistrationData struct {
	Username   string `json:"userName"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

func (u UserRegistrationData) GetHashPassword() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

type User struct {
	Id        int64  `json:"id" db:"id"`
	Username  string `json:"userName" db:"username"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"passwordhash"`
	CreatedAt string `json:"createdAt" db:"createdat"`
	UpdatedAt string `json:"updatedAt" db:"updatedat"`
}
