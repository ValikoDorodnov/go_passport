package v1

import (
	"net/http"

	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/request"

	"github.com/ValikoDorodnov/go_passport/pkg/rest"
)

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestDto request.LoginByEmail

	err := rest.ParseRequestBody(r.Body, &requestDto)
	if err != nil {
		rest.ResponseError(w, err)
		return
	}

	if isOk, errs := h.validator.Validate(requestDto); !isOk {
		rest.ResponseErrors(w, errs)
		return
	}

	resp, err := h.auth.SignIn(ctx, &requestDto)
	if err != nil {
		rest.ResponseError(w, err)
		return
	}

	rest.ResponseOk(w, resp)
	return
}
