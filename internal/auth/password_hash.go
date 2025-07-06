package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytePW := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePW, bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("Something went wrong")
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
