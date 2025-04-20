package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/beevik/guid"
	"github.com/go-chi/render"

	"auth_service/api/types"
	"auth_service/usecases"
)

type authHandler struct {
	service usecases.AuthService
}

func NewAuthHandler(authService usecases.AuthService) *authHandler {
	return &authHandler{service: authService}
}

// @Summary Generate Tokens
// @Description Generate access_token and refresh_token
// @Param json-body body types.GenerateTokensRequest true "guid"
// @Success 200 {object} types.TokensResponse "Token"
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /generate_tokens [post]
func (h *authHandler) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	var req types.GenerateTokensRequest
	ip := r.RemoteAddr
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		types.BadRequest(w, r, types.ErrorResponse{Message: "invalid request"})
		return
	}
	if req.UserID == "" {
		types.BadRequest(w, r, types.ErrorResponse{Message: "missing guid"})
		return
	}
	if !guid.IsGuid(req.UserID) {
		types.BadRequest(w, r, types.ErrorResponse{Message: "guid is not valid"})
		return
	}

	tokens, err := h.service.GenerateTokens(req.UserID, ip)
	if err != nil {
		types.ProcessError(w, r, types.ErrorResponse{Message: err.Error(), Err: err})
		return
	}

	resp := types.TokensResponse{
		AccessToken:  string(tokens.AccessToken),
		RefreshToken: string(tokens.RefreshToken),
	}

	render.Status(r, http.StatusOK)

	render.JSON(w, r, resp)
}

// @Summary Refresh Tokens
// @Description Refresh access_token and refresh_token
// @Param json-body body types.RefreshTokensRequest true "guid and refresh_token"
// @Param Authorization header string true "access_token with Bearer"
// @Success 200 {object} types.TokensResponse "Token"
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /refresh_tokens [post]
func (h *authHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var req types.RefreshTokensRequest
	ip := r.RemoteAddr
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		types.BadRequest(w, r, types.ErrorResponse{Message: "invalid request"})
		return
	}
	if req.UserID == "" {
		types.BadRequest(w, r, types.ErrorResponse{Message: "missing guid"})
		return
	}
	if !guid.IsGuid(req.UserID) {
		types.BadRequest(w, r, types.ErrorResponse{Message: "guid is not valid"})
		return
	}
	if req.RefreshToken == "" {
		types.BadRequest(w, r, types.ErrorResponse{Message: "missing refresh_token"})
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		types.BadRequest(w, r, types.ErrorResponse{Message: "missing authorization header"})
		return
	}
	t := strings.Split(authHeader, " ")
	if len(t) != 2 {
		types.BadRequest(w, r, types.ErrorResponse{Message: "authorization header has a wrong format"})
		return
	}

	accessToken := t[1]

	tokens, err := h.service.RefreshTokens(req.UserID, ip, accessToken, req.RefreshToken)
	if err != nil {
		types.ProcessError(w, r, types.ErrorResponse{Message: err.Error(), Err: err})
		return
	}

	resp := types.TokensResponse{
		AccessToken:  string(tokens.AccessToken),
		RefreshToken: string(tokens.RefreshToken),
	}

	render.Status(r, http.StatusOK)

	render.JSON(w, r, resp)
}
