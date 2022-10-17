package v1

import (
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/request"
	"net/http"

	"github.com/ValikoDorodnov/go_passport/pkg/rest"
)

func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestDto request.Refresh

	err := rest.ParseRequestBody(r.Body, &requestDto)
	if err != nil {
		rest.ResponseErrors(w, err)
		return
	}

	resp, err := h.auth.RefreshTokens(ctx, &requestDto)
	if err != nil {
		rest.ResponseErrors(w, err)
		return
	}
	rest.ResponseOk(w, resp)
}
