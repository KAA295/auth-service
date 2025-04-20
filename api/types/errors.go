package types

import (
	"auth_service/domain"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Message string `json:"message" example:"internal error"`
	Err     error  `json:"-"`
}

func BadRequest(w http.ResponseWriter, r *http.Request, errResp ErrorResponse) {
	render.Status(r, http.StatusBadRequest)

	render.JSON(w, r, errResp)
}

func NotFound(w http.ResponseWriter, r *http.Request, errResp ErrorResponse) {
	render.Status(r, http.StatusNotFound)

	render.JSON(w, r, errResp)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, errResp ErrorResponse) {
	render.Status(r, http.StatusUnauthorized)

	render.JSON(w, r, errResp)
}

func Internal(w http.ResponseWriter, r *http.Request, errResp ErrorResponse) {
	render.Status(r, http.StatusInternalServerError)

	render.JSON(w, r, errResp)
}

func ProcessError(w http.ResponseWriter, r *http.Request, errResp ErrorResponse) {
	if errors.Is(errResp.Err, domain.ErrNotFound) {
		NotFound(w, r, errResp)
		return
	}
	if errors.Is(errResp.Err, domain.ErrBadRequest) {
		BadRequest(w, r, errResp)
		return
	}
	if errors.Is(errResp.Err, domain.ErrUnauthorized) {
		Unauthorized(w, r, errResp)
		return
	}
	fmt.Println(errResp)
	Internal(w, r, ErrorResponse{Message: "internal error"})
}
