package jwt_generate

import (
	"context"
	"errors"

	"backend/pkg/cache"
	configs "backend/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtGenerate struct {
	secretKey             string
	issuer                string
	audience              string
	expiresAt             time.Duration
	refreshTokenSecretKey string
	refreshTokenExpiresAt time.Duration
	redisCache            cache.Cache
	ctx                   context.Context
	verifyEmailSecretKey  string
	verifyEmailExpiresAt  time.Duration
}
type TokenPayload struct {
	UserId uuid.UUID
	Email  string
}
type JwtGenerate interface {
	GenerateToken(user *TokenPayload) (string, error)
	GenerateVerifyEmailToken(user *TokenPayload) (string, error)
	GenerateRefreshToken(user *TokenPayload) (string, error)
	VerifyToken(refreshToken string, secretKey string) (*TokenPayload, error)
}

func NewJwtGenerate(ctx context.Context, config *configs.AppConfig, redisCache cache.Cache) JwtGenerate {
	return &jwtGenerate{
		secretKey:             config.Jwt.SecretKey,
		issuer:                config.Jwt.Issuer,
		audience:              config.Jwt.Audience,
		expiresAt:             time.Duration(config.Jwt.TokenExpire) * time.Minute,
		refreshTokenSecretKey: config.Jwt.RefreshSecretKey,
		refreshTokenExpiresAt: time.Duration(config.Jwt.RefreshTokenExpire) * 24 * time.Hour,
		redisCache:            redisCache,
		ctx:                   ctx,
		verifyEmailSecretKey:  config.Jwt.VerifyEmailSecretKey,
		verifyEmailExpiresAt:  time.Duration(config.Jwt.VerifyEmailTokenExpire) * time.Hour,
	}
}

func (j *jwtGenerate) generateTokenWithClaims(user *TokenPayload, secretKey string, expiresAt time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"email": user.Email,
		"id":    user.UserId,
		"exp":   time.Now().Add(expiresAt).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   j.issuer,
		"aud":   j.audience,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func (j *jwtGenerate) GenerateToken(user *TokenPayload) (string, error) {
	return j.generateTokenWithClaims(user, j.secretKey, j.expiresAt)
}

func (j *jwtGenerate) GenerateVerifyEmailToken(user *TokenPayload) (string, error) {
	return j.generateTokenWithClaims(user, j.verifyEmailSecretKey, j.verifyEmailExpiresAt)
}

func (j *jwtGenerate) GenerateRefreshToken(user *TokenPayload) (string, error) {
	refreshToken, err := j.generateTokenWithClaims(user, j.refreshTokenSecretKey, j.refreshTokenExpiresAt)
	if err != nil {
		return "", err
	}

	err = j.redisCache.Set(j.ctx, cache.RefreshTokenKey(user.UserId), refreshToken, int(j.refreshTokenExpiresAt))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
func (j *jwtGenerate) VerifyToken(refreshToken string, secretKey string) (*TokenPayload, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if iss, ok := claims["iss"].(string); !ok || iss != j.issuer {
		return nil, errors.New("invalid issuer")
	}
	if aud, ok := claims["aud"].(string); !ok || aud != j.audience {
		return nil, errors.New("invalid audience")
	}

	if exp, ok := claims["exp"].(float64); ok && exp < float64(time.Now().Unix()) {
		return nil, errors.New("token expired")
	}

	result := &TokenPayload{
		UserId: uuid.MustParse(claims["id"].(string)),
		Email:  claims["email"].(string),
	}
	return result, nil
}
