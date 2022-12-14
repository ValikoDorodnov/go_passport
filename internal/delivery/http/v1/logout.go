package v1

import (
	"net/http"

	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/request"

	"github.com/ValikoDorodnov/go_passport/pkg/rest"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestDto request.Logout

	err := rest.ParseRequestBody(r.Body, &requestDto)
	if err != nil {
		rest.ResponseError(w, err)
		return
	}

	err = h.auth.Logout(ctx, &requestDto, h.middleware.ParsedToken)
	if err != nil {
		rest.ResponseError(w, err)
		return
	}

	rest.ResponseOk(w, "Success logout")
	return
}
