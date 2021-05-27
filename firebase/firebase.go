package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/xerrors"
)

const contentTypeJSON = "application/json"

type Client struct {
	apiKey string
	Client *http.Client
}

type ErrorResponse struct {
	ErrorPayload struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("firebase error. code = %d, message = %s", e.ErrorPayload.Code, e.ErrorPayload.Message)
}

func NewClient(apiKey string, opts ...func(client *Client)) *Client {
	c := &Client{apiKey: apiKey, Client: http.DefaultClient}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type VerifyCustomTokenRequest struct {
	Token             string `json:"token"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type VerifyCustomTokenResponse struct {
	Kind         string `json:"kind"`
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	IsNewUser    bool   `json:"isNewUser"`
}

func (c *Client) VerifyCustomToken(ctx context.Context, req VerifyCustomTokenRequest) (*VerifyCustomTokenResponse, error) {
	u := fmt.Sprintf("https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s", c.apiKey)

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&req); err != nil {
		return nil, xerrors.Errorf("failed to encode: %w", err)
	}

	resp, err := c.postJSON(ctx, u, &b)
	if err != nil {
		return nil, xerrors.Errorf("failed to postJSON: %w", err)
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var tokenResponse VerifyCustomTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, xerrors.Errorf("failed to decode response: %w", err)
	}

	return &tokenResponse, nil
}

func handleErrorResponse(resp *http.Response) error {
	var errorResponse ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
		return xerrors.Errorf("failed to decode error response: %w", err)
	}
	return errorResponse
}

func (c *Client) postJSON(ctx context.Context, url string, r io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentTypeJSON)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
