package middleware

import (
	"net/http"

	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/ValikoDorodnov/go_passport/internal/repository"
	"github.com/ValikoDorodnov/go_passport/internal/service"
)

type AuthMiddleware struct {
	ParsedToken *entity.ParsedToken
	jwt         *service.JwtService
	sessionRepo *repository.SessionRepository
}

func NewAuthMiddleware(jwt *service.JwtService, sessionRepo *repository.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwt:         jwt,
		sessionRepo: sessionRepo,
	}
}

func (r *AuthMiddleware) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		r.ParsedToken = nil
		token := request.Header.Get("X-Access-Token")
		if token != "" {
			ctx := request.Context()

			revoked := r.sessionRepo.CheckTokenIsInBlackList(ctx, token)
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
