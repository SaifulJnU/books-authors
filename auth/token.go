// auth/token.go

package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("mysecretkey")

// Generate a token
func GenerateToken(username string) (string, error) {
	//Define signed method
	token := jwt.New(jwt.SigningMethodHS256)

	//Define claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	//Sign the JWT token with the secret key to generate the final JWT string
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
