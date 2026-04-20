package token

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey []byte = nil

func setRandomSecretKey() {
	if secretKey == nil {
		secretKey = make([]byte, 64)
		_, err := rand.Read(secretKey)
		if err != nil {
			fmt.Printf("Failed to generate secret key: %v\n", err)
			return
		}

		log.Println(string(secretKey))
	}
}

func CreateToken(uuid string) (string, error) {
	setRandomSecretKey()
	// Create claims with multiple fields
	claims := jwt.MapClaims{
		"uuid": uuid,
		"exp":  time.Now().Add(time.Hour * 8).Unix(), // Expires in 8 hours
		"iat":  time.Now().Unix(),                    // Issued at
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
