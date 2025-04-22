package server

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	answer, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(answer), err
}
