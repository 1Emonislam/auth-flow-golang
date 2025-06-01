package database

import (
	"context"
	"fmt"
	"own-paynet/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedis initializes the Redis client
func InitRedis(cfg *config.Config) *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Ping Redis to check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	} else {
		fmt.Println("Connected to Redis")
	}

	return redisClient
}

// GetRedisClient returns the Redis client
func GetRedisClient() *redis.Client {
	return redisClient
}

// StoreToken stores a token in Redis with expiry
func StoreToken(ctx context.Context, userID uint, token string, expiry time.Duration) error {
	key := fmt.Sprintf("user:%d:token", userID)
	return redisClient.Set(ctx, key, token, expiry).Err()
}

// GetToken retrieves a token from Redis
func GetToken(ctx context.Context, userID uint) (string, error) {
	key := fmt.Sprintf("user:%d:token", userID)
	return redisClient.Get(ctx, key).Result()
}

// DeleteToken deletes a token from Redis
func DeleteToken(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("user:%d:token", userID)
	return redisClient.Del(ctx, key).Err()
}

// ValidateToken checks if a token exists and matches in Redis
func ValidateToken(ctx context.Context, userID uint, token string) (bool, error) {
	storedToken, err := GetToken(ctx, userID)
	if err != nil {
		return false, err
	}

	return storedToken == token, nil
}

// StorePasswordResetToken stores a password reset token in Redis with expiry
func StorePasswordResetToken(ctx context.Context, email string, token string, expiry time.Duration) error {
	key := fmt.Sprintf("password_reset:%s", email)
	return redisClient.Set(ctx, key, token, expiry).Err()
}

// GetPasswordResetToken retrieves a password reset token from Redis
func GetPasswordResetToken(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("password_reset:%s", email)
	return redisClient.Get(ctx, key).Result()
}

// DeletePasswordResetToken deletes a password reset token from Redis
func DeletePasswordResetToken(ctx context.Context, email string) error {
	key := fmt.Sprintf("password_reset:%s", email)
	return redisClient.Del(ctx, key).Err()
}

// ValidatePasswordResetToken checks if a token exists and matches in Redis
func ValidatePasswordResetToken(ctx context.Context, email string, token string) (bool, error) {
	storedToken, err := GetPasswordResetToken(ctx, email)
	if err != nil {
		if err == redis.Nil {
			return false, fmt.Errorf("password reset token not found for email: %s", email)
		}
		return false, fmt.Errorf("error retrieving password reset token: %v", err)
	}

	return storedToken == token, nil
}

// StoreEmailVerificationToken stores an email verification token in Redis with expiry
func StoreEmailVerificationToken(ctx context.Context, email string, token string, expiry time.Duration) error {
	key := fmt.Sprintf("email_verification:%s", email)
	return redisClient.Set(ctx, key, token, expiry).Err()
}

// GetEmailVerificationToken retrieves an email verification token from Redis
func GetEmailVerificationToken(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("email_verification:%s", email)
	return redisClient.Get(ctx, key).Result()
}

// DeleteEmailVerificationToken deletes an email verification token from Redis
func DeleteEmailVerificationToken(ctx context.Context, email string) error {
	key := fmt.Sprintf("email_verification:%s", email)
	return redisClient.Del(ctx, key).Err()
}

// ValidateEmailVerificationToken checks if a token exists and matches in Redis
func ValidateEmailVerificationToken(ctx context.Context, email string, token string) (bool, error) {
	storedToken, err := GetEmailVerificationToken(ctx, email)
	if err != nil {
		if err == redis.Nil {
			// Token does not exist
			return false, nil
		}
		return false, err
	}

	return storedToken == token, nil
}
