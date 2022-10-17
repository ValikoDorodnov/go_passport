package v1

import (
	"net/http"

	"github.com/ValikoDorodnov/go_passport/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	auth *service.AuthService
}

func NewHandler(auth *service.AuthService) *Handler {
	return &Handler{
		auth: auth,
	}
}

func (h *Handler) GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/sign-in", h.SignIn).Methods(http.MethodPost)
	r.HandleFunc("/refresh-tokens", h.RefreshTokens).Methods(http.MethodPost)

	return r
}
