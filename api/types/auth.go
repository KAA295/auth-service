package types

type TokensResponse struct {
	AccessToken  string `json:"access_token" example:"access_token"`
	RefreshToken string `json:"refresh_token" example:"refresh_token"`
}

type GenerateTokensRequest struct {
	UserID string `json:"guid" example:"e8207e59-127e-4557-bd66-6c43c427c109"`
}

type RefreshTokensRequest struct {
	UserID       string `json:"guid" example:"e8207e59-127e-4557-bd66-6c43c427c109"`
	RefreshToken string `json:"refresh_token" example:"refresh_token"`
}
