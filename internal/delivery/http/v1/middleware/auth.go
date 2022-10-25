package middleware

import (
	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/ValikoDorodnov/go_passport/internal/repository"
	"github.com/ValikoDorodnov/go_passport/internal/service"
	"net/http"
)

type AuthMiddleware struct {
	jwt           *service.JwtService
	ParsedToken   *entity.ParsedToken
	accessSession *repository.AccessSessionRepository
}

func NewAuthMiddleware(jwt *service.JwtService, accessSession *repository.AccessSessionRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwt:           jwt,
		accessSession: accessSession,
	}
}

func (r *AuthMiddleware) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		r.ParsedToken = nil
		token := request.Header.Get("X-Access-Token")
		if token != "" {
			ctx := request.Context()

			revoked := r.accessSession.CheckTokenIsInBlackList(ctx, token)
			if !revoked {
				parsedToken, err := r.jwt.ParseToken(token)
				if err == nil {
					r.ParsedToken = parsedToken
				}
			}
		}
		next.ServeHTTP(w, request)
	})
}
