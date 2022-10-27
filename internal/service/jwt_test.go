package service

import (
	"regexp"
	"testing"
	"time"

	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/stretchr/testify/assert"
)

const (
	accessTtl  = 10
	refreshTtl = 100
)

func TestJwtService_IssueAccess(t *testing.T) {
	jwt := initJwtService()
	user := initUser()

	token, err := jwt.IssueAccess(user)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	jwtRx := `(^[A-Za-z0-9-_]*\.[A-Za-z0-9-_]*\.[A-Za-z0-9-_]*$)`
	assert.Regexp(t, regexp.MustCompile(jwtRx), token.Value)

	exp := time.Now().Add(time.Second * time.Duration(accessTtl)).Unix()
	assert.Equal(t, exp, token.Exp)
}

func TestJwtService_IssueRefresh(t *testing.T) {
	jwt := initJwtService()
	token := jwt.IssueRefresh()

	jwtRx := `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
	assert.Regexp(t, regexp.MustCompile(jwtRx), token.Value)

	exp := time.Now().Add(time.Second * time.Duration(refreshTtl)).Unix()
	assert.Equal(t, exp, token.Exp)
}

func TestJwtService_ParseToken(t *testing.T) {
	jwt := initJwtService()
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4MjQ1NzA0NjMsInJvbGVzIjoidGVzdCxiZXN0Iiwic3ViIjoiNTU1In0.HUkzPzM_7aEpkDDGqsHGuL97gf5kcfZ3-oSr0og_iIQ"
	parsed, err := jwt.ParseToken(tokenStr)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, "555", parsed.Subject)
}

func initJwtService() *JwtService {
	conf := config.JwtConfig{
		AccessTtl:  accessTtl,
		RefreshTtl: refreshTtl,
		SecretKey:  "my_secret",
		Issuer:     "my_issuer",
	}
	return NewJwtService(conf)
}

func initUser() *entity.User {
	return &entity.User{
		CommonId: "test",
		Roles:    "best",
	}
}
