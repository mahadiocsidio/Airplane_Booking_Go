package utils

import(
	"time"
	"os"

	"github.com/golang-jwt/jwt/v5"
)
var secret_key = os.Getenv("secretkey")
var jwtSecret = []byte(secret_key)

func GenerateToken(userId string, role string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // expire 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}