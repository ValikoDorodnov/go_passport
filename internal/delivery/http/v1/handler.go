package v1

import (
	"net/http"

	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/middleware"
	"github.com/ValikoDorodnov/go_passport/pkg/validator"

	"github.com/ValikoDorodnov/go_passport/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	auth       *service.AuthService
	middleware *middleware.AuthMiddleware
	validator  *validator.Validator
}

func NewHandler(
	auth *service.AuthService,
	middleware *middleware.AuthMiddleware,
	validator *validator.Validator,
) *Handler {
	return &Handler{
		auth:       auth,
		middleware: middleware,
		validator:  validator,
	}
}

func (h *Handler) GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(h.middleware.CheckAuth)

	r.HandleFunc("/sign-in", h.SignIn).Methods(http.MethodPost)
	r.HandleFunc("/refresh-tokens", h.RefreshTokens).Methods(http.MethodPost)
	r.HandleFunc("/logout", h.Logout).Methods(http.MethodPost)

	return r
}
