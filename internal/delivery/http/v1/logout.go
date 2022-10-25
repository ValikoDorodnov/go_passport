package v1

import (
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/request"
	"net/http"

	"github.com/ValikoDorodnov/go_passport/pkg/rest"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestDto request.Logout

	err := rest.ParseRequestBody(r.Body, &requestDto)
	if err != nil {
		rest.ResponseErrors(w, err)
		return
	}

	err = h.auth.Logout(ctx, &requestDto, h.middleware.ParsedToken)
	if err != nil {
		rest.ResponseErrors(w, err)
		return
	}
	rest.ResponseOk(w, "Success logout")
}
