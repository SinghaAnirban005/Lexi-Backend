package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	isPasswordCorrect := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return isPasswordCorrect == nil
}
