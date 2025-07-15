package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func CheckPassword(input, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
	return err == nil
}
