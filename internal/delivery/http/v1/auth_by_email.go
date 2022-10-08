package v1

import (
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/dto"
	"net/http"

	"github.com/ValikoDorodnov/go_passport/pkg/rest"
)

func (h *Handler) AuthByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestDto dto.RequestDto
	err := rest.ParseRequestBody(r.Body, &requestDto)
	if err != nil {
		rest.ResponseErrors(w, err)
	}

	resp, err := h.userService.AuthByEmailProcess(ctx, &requestDto)
	if err != nil {
		rest.ResponseErrors(w, err)
	}
	rest.ResponseOk(w, resp)
}
