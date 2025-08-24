package token

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func CreateToken(user_name string) (string, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var secretKey = []byte(os.Getenv("SECRETKEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user_name,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
