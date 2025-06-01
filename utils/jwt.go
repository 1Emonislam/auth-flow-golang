package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"own-paynet/config"
	"own-paynet/database"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
)

// GenerateJWT creates a new JWT token for the given user ID
func GenerateJWT(userID uint) (string, error) {
	cfg := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(cfg.TokenExpiry)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	// Store token in Redis
	ctx := context.Background()
	err = database.StoreToken(ctx, userID, tokenString, time.Hour*time.Duration(cfg.TokenExpiry))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (uint, error) {
	cfg := config.LoadConfig()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))

		// Validate token against Redis
		ctx := context.Background()
		valid, err := database.ValidateToken(ctx, userID, tokenString)
		if err != nil {
			if err == redis.Nil {
				return 0, jwt.ErrSignatureInvalid // Token not found in Redis
			}
			return 0, err
		}

		if !valid {
			return 0, jwt.ErrSignatureInvalid // Token doesn't match stored token
		}

		return userID, nil
	}

	return 0, jwt.ErrSignatureInvalid
}

// InvalidateToken removes a token from Redis
func InvalidateToken(userID uint) error {
	ctx := context.Background()
	return database.DeleteToken(ctx, userID)
}

// GenerateRandomToken generates a random string of specified length
func GenerateRandomToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
