package models

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
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

func (urd UserRegistrationData) GetHashPassword() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(urd.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (urd UserRegistrationData) Validate() error {
	switch {
	case urd.Username == "":
		return fmt.Errorf("user_name обязательное поле")
	case urd.Email == "":
		return fmt.Errorf("email обязательное поле")
	case urd.Password == "":
		return fmt.Errorf("password обязательное поле")
	case len(urd.Password) < 6:
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

func (u User) JWTGeneration(secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   u.Id,
		"username": u.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}
