package repository

import (
	"context"
	"fmt"
	"ticert/config"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthRepository interface {
	StoreAccessToken(userID uuid.UUID, token string, expiry time.Duration) error
	StoreRefreshToken(userID uuid.UUID, token string, expiry time.Duration) error
	ValidateAccessToken(userID uuid.UUID, token string) (bool, error)
	ValidateRefreshToken(userID uuid.UUID, token string) (bool, error)
	RevokeAccessToken(userID uuid.UUID) error
	RevokeRefreshToken(userID uuid.UUID) error
	RevokeAllUserTokens(userID uuid.UUID) error
}

type authRepository struct {
	redisClient *redis.Client
}

func NewAuthRepository() AuthRepository {
	return &authRepository{
		redisClient: config.GetRedisClient(),
	}
}

func (r *authRepository) StoreAccessToken(userID uuid.UUID, token string, expiry time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("access_token:%s", userID.String())
	return r.redisClient.Set(ctx, key, token, expiry).Err()
}

func (r *authRepository) StoreRefreshToken(userID uuid.UUID, token string, expiry time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", userID.String())
	return r.redisClient.Set(ctx, key, token, expiry).Err()
}

func (r *authRepository) ValidateAccessToken(userID uuid.UUID, token string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("access_token:%s", userID.String())

	storedToken, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return storedToken == token, nil
}

func (r *authRepository) ValidateRefreshToken(userID uuid.UUID, token string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", userID.String())

	storedToken, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return storedToken == token, nil
}

func (r *authRepository) RevokeAccessToken(userID uuid.UUID) error {
	ctx := context.Background()
	key := fmt.Sprintf("access_token:%s", userID.String())
	return r.redisClient.Del(ctx, key).Err()
}

func (r *authRepository) RevokeRefreshToken(userID uuid.UUID) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", userID.String())
	return r.redisClient.Del(ctx, key).Err()
}

func (r *authRepository) RevokeAllUserTokens(userID uuid.UUID) error {
	ctx := context.Background()
	accessKey := fmt.Sprintf("access_token:%s", userID.String())
	refreshKey := fmt.Sprintf("refresh_token:%s", userID.String())

	pipe := r.redisClient.Pipeline()
	pipe.Del(ctx, accessKey)
	pipe.Del(ctx, refreshKey)
	_, err := pipe.Exec(ctx)
	return err
}
