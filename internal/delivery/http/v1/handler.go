package v1

import (
	"net/http"

	"github.com/ValikoDorodnov/go_passport/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	userService *service.UserService
}

func NewHandler(userService *service.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login-by-email", h.LoginByEmail).Methods(http.MethodPost)

	return r
}
