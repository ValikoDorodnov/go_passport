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

func (s *JwtService) IssueAccess(user *entity.User) (*entity.Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	now := time.Now()
	exp := now.Add(time.Second * time.Duration(s.config.AccessTtl))

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.CommonId
	claims["roles"] = user.Roles
	claims["exp"] = exp.Unix()

	signedToken, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return nil, err
	}

	return &entity.Token{
		Value: signedToken,
		Exp:   exp.Sub(now),
	}, nil
}

func (s *JwtService) IssueRefresh() *entity.Token {
	now := time.Now()
	exp := now.Add(time.Second * time.Duration(s.config.RefreshTtl))
	return &entity.Token{
		Value: uuid.New().String(),
		Exp:   exp.Sub(now),
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
