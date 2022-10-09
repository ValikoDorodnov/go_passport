package service

import (
	"errors"
	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtService struct {
	config config.JwtConfig
}

func NewJwtService(c config.JwtConfig) *JwtService {
	return &JwtService{config: c}
}

type MyClaims struct {
	jwt.RegisteredClaims
	CommonId  int    `json:"common_id"`
	Roles     string `json:"roles"`
	TokenType string `json:"token_type"`
}

func (s JwtService) IssueToken(user *entity.User, ttl int, tokenType string) string {
	token := newToken()

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["common_id"] = user.CommonId
	claims["roles"] = user.Roles
	claims["token_type"] = tokenType
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(ttl)).Unix()

	return signToken(token, []byte(s.config.SecretKey))
}

func (s JwtService) ParseToken(tokenStr string) (*entity.User, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.SecretKey), nil
	})

	claims, ok := token.Claims.(*MyClaims)

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	if ok && token.Valid && claims.VerifyExpiresAt(time.Now(), true) {
		return &entity.User{
			CommonId: claims.CommonId,
			Roles:    claims.Roles,
		}, nil
	}
	return nil, err
}

func newToken() *jwt.Token {
	return jwt.New(jwt.SigningMethodHS256)
}

func signToken(token *jwt.Token, secretKey []byte) string {
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
	}
	return tokenString
}
