package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/beevik/guid"

	"auth_service/api/types"
)

type GenerateTokensTestCase struct {
	Request        types.GenerateTokensRequest
	ExpectedStatus int
}

func TestGenerateTokens(t *testing.T) {
	GUID := guid.NewString()
	tests := []GenerateTokensTestCase{
		{
			Request: types.GenerateTokensRequest{
				UserID: GUID,
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Request: types.GenerateTokensRequest{
				UserID: "not guid",
			},
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	for i, test := range tests {
		body, err := json.Marshal(test.Request)
		if err != nil {
			t.Fatalf("could not marshal to JSON: %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, "http://service:8000/generate_tokens", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("error while formatting request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != test.ExpectedStatus {
			t.Errorf("response number %d was incorrect, got: %d, want: %d.", i, resp.StatusCode, test.ExpectedStatus)
		}

	}
}

type RefreshTokensTestCase struct {
	Request        types.RefreshTokensRequest
	AuthHeader     string
	ExpectedStatus int
}

func TestRefreshTokens(t *testing.T) {
	GUID := guid.NewString()

	secondGUID := guid.NewString()

	responseStruct, err := getTokens(GUID)
	if err != nil {
		t.Fatal(err)
	}
	secondResponseStruct, err := getTokens(GUID)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []RefreshTokensTestCase{
		{
			Request: types.RefreshTokensRequest{
				UserID: GUID,
			},
			AuthHeader:     responseStruct.AccessToken,
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Request: types.RefreshTokensRequest{
				UserID:       secondGUID,
				RefreshToken: responseStruct.RefreshToken,
			},
			AuthHeader:     responseStruct.AccessToken,
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Request: types.RefreshTokensRequest{
				UserID:       GUID,
				RefreshToken: responseStruct.RefreshToken,
			},
			AuthHeader:     secondResponseStruct.AccessToken,
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Request: types.RefreshTokensRequest{
				UserID:       GUID,
				RefreshToken: secondResponseStruct.RefreshToken,
			},
			AuthHeader:     responseStruct.AccessToken,
			ExpectedStatus: http.StatusUnauthorized,
		},
		{
			Request: types.RefreshTokensRequest{
				UserID:       GUID,
				RefreshToken: responseStruct.RefreshToken,
			},
			AuthHeader:     responseStruct.AccessToken,
			ExpectedStatus: http.StatusOK,
		},
	}

	for i, test := range testCases {
		body, err := json.Marshal(test.Request)
		if err != nil {
			t.Fatalf("could not marshal to JSON: %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, "http://service:8000/refresh_tokens", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("error while formatting request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.AuthHeader)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != test.ExpectedStatus {
			t.Errorf("response number %d was incorrect, got: %d, want: %d.", i, resp.StatusCode, test.ExpectedStatus)
		}

	}
}

func getTokens(userID string) (types.TokensResponse, error) {
	responseStruct := new(types.TokensResponse)
	request := types.GenerateTokensRequest{
		UserID: userID,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return *responseStruct, fmt.Errorf("could not marshal to JSON: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://service:8000/generate_tokens", bytes.NewBuffer(body))
	if err != nil {
		return *responseStruct, fmt.Errorf("error while formatting request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return *responseStruct, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return *responseStruct, fmt.Errorf("error reading body: %w", err)
	}

	err = json.Unmarshal(respBody, responseStruct)
	if err != nil {
		return *responseStruct, fmt.Errorf("could not unmarshal JSON: %w", err)
	}
	if responseStruct.AccessToken == "" || responseStruct.RefreshToken == "" {
		return *responseStruct, fmt.Errorf("wrong response")
	}
	return *responseStruct, nil
}
