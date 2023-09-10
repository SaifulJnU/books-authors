// auth/auth.go

package auth

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var vjwtSecret = []byte("mysecretkey")

func VerifyToken(tokenString string) (jwt.MapClaims, error) {

	//Parse the tokenString and Check the signing method
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return vjwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	//verify claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
