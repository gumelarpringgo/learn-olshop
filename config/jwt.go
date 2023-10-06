package config

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte(os.Getenv("KEY_JWT"))

func CreateToken(userId int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"role":    role,
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Parse(tokenString string) (any, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("signing method invalid")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("unauthorized token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("unauthorized validation")
	}

	return claims, nil
}
