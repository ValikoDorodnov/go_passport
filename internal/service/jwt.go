package service

import (
	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type JwtService struct {
	config config.JwtConfig
}

func NewJwtService(c config.JwtConfig) *JwtService {
	return &JwtService{config: c}
}

func (s JwtService) IssueAccess(user *entity.User) *entity.Token {
	token := jwt.New(jwt.SigningMethodHS256)
	exp := time.Now().Add(time.Minute * time.Duration(s.config.AccessTtl)).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["subject"] = user.CommonId
	claims["roles"] = user.Roles
	claims["exp"] = exp

	return &entity.Token{
		Value: signToken(token, []byte(s.config.SecretKey)),
		Exp:   exp,
	}
}

func (s JwtService) IssueRefresh() *entity.Token {
	exp := time.Now().Add(time.Minute * time.Duration(s.config.RefreshTtl)).Unix()
	token := uuid.New()

	return &entity.Token{
		Value: token.String(),
		Exp:   exp,
	}
}

func signToken(token *jwt.Token, secretKey []byte) string {
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
	}
	return tokenString
}
