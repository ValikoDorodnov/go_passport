package service

import (
	"errors"
	"time"

	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type MyClaims struct {
	jwt.RegisteredClaims
}

type JwtService struct {
	config config.JwtConfig
}

func NewJwtService(c config.JwtConfig) *JwtService {
	return &JwtService{config: c}
}

func (s *JwtService) IssueAccess(user *entity.User) *entity.Token {
	token := jwt.New(jwt.SigningMethodHS256)
	exp := time.Now().Add(time.Second * time.Duration(s.config.AccessTtl)).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.CommonId
	claims["roles"] = user.Roles
	claims["exp"] = exp

	return &entity.Token{
		Value: signToken(token, []byte(s.config.SecretKey)),
		Exp:   exp,
	}
}

func (s *JwtService) IssueRefresh() *entity.Token {
	return &entity.Token{
		Value: uuid.New().String(),
		Exp:   time.Now().Add(time.Second * time.Duration(s.config.RefreshTtl)).Unix(),
	}
}

func (s *JwtService) ParseToken(tokenStr string) (*entity.ParsedToken, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, errors.New("invalid signing method")
		}
		return []byte(s.config.SecretKey), nil
	})

	claims, ok := token.Claims.(*MyClaims)
	now := time.Now()

	exp := claims.RegisteredClaims.ExpiresAt.Time
	diff := exp.Sub(now)

	if ok && token.Valid && claims.VerifyExpiresAt(now, true) {
		return &entity.ParsedToken{
			Subject: claims.RegisteredClaims.Subject,
			ExpTtl:  diff,
			Jwt:     tokenStr,
		}, nil
	}
	return nil, err
}

func signToken(token *jwt.Token, secretKey []byte) string {
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
	}
	return tokenString
}
