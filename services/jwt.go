package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte(GetEnv("SECRET_KEY"))

type CustomClaims struct {
	Email string    `json:"email"`
	jwt.StandardClaims
}

func GenerateJWTToken(email string, expirationTime time.Time) (string, error) {

	
	// Создаем структуру CustomClaims с пользовательскими полями
	claims := &CustomClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Создаем новый токен с настройками подписи
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с использованием секретного ключа
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseJWTToken проверяет и парсит JWT токен
func ParseJWTToken(tokenString string) (*CustomClaims, error) {
	// Парсим токен с использованием секретного ключа
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Проверяем, что токен действителен и преобразуем его в пользовательские данные
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("неверный токен")
}