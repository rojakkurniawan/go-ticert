package jwt

import (
	"errors"
	"strconv"
	"ticert/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func getAccessTokenSecret() []byte {
	cfg := config.GetConfig()
	return []byte(cfg.JWTAccessSecret)
}

func getRefreshTokenSecret() []byte {
	cfg := config.GetConfig()
	return []byte(cfg.JWTRefreshSecret)
}

func GetAccessTokenExpiry() (time.Duration, error) {
	cfg := config.GetConfig()
	hours, err := strconv.Atoi(cfg.JWTAccessExpiry)
	if err != nil {
		return 0, errors.New("invalid JWT_ACCESS_EXPIRY value: must be a valid integer")
	}
	return time.Hour * time.Duration(hours), nil
}

func GetRefreshTokenExpiry() (time.Duration, error) {
	cfg := config.GetConfig()
	hours, err := strconv.Atoi(cfg.JWTRefreshExpiry)
	if err != nil {
		return 0, errors.New("invalid JWT_REFRESH_EXPIRY value: must be a valid integer")
	}
	return time.Hour * time.Duration(hours), nil
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uuid.UUID, email, role string) (string, error) {
	expiry, err := GetAccessTokenExpiry()
	if err != nil {
		return "", err
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getAccessTokenSecret())
}

func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	expiry, err := GetRefreshTokenExpiry()
	if err != nil {
		return "", err
	}

	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getRefreshTokenSecret())
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getAccessTokenSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getRefreshTokenSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}
