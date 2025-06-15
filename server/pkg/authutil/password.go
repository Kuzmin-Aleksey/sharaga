package authutil

import (
	"crypto/sha256"
	"fmt"
)

var passwordSalt = []byte{55, 83, 78, 126, 35, 65}

func HashPassword(password string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(passwordSalt)), nil
}

func EqualPassword(password, hashed string) (error, bool) {
	hash := sha256.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return err, false
	}
	passwordHash := fmt.Sprintf("%x", hash.Sum(passwordSalt))
	return nil, passwordHash == hashed
}
