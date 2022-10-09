package v1

import (
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/dto"
	"net/http"

	"github.com/ValikoDorodnov/go_passport/pkg/rest"
)

func (h *Handler) LoginByRefresh(w http.ResponseWriter, r *http.Request) {
	var requestDto dto.LoginByRefreshDto
	err := rest.ParseRequestBody(r.Body, &requestDto)
	if err != nil {
		rest.ResponseErrors(w, err)
		return
	}

	resp, err := h.userService.LoginByRefresh(&requestDto)
	if err != nil {
		rest.ResponseErrors(w, err)
		return
	}
	rest.ResponseOk(w, resp)
}
