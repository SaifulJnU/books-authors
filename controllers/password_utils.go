// controllers/password_utils.go

package controllers

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password and returns the hashed password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword compares a plain text password with a hashed password and returns true if they match.
func VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
