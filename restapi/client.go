package restapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client is a Jexia REST API client
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// ClientResp is the Jexia REST API client response
type ClientResp struct {
	// Success is a true|false status that indicates the response token validation result.
	Success bool
	// Error is the user response error codes list.
	Error []string `json:"error"`
	// Jexia user token
	Token string `json:"token"`
	// RefreshTokenn is long lived (30 days) and can be used to obtain a new access token.
	RefreshToken string `json:"refresh_token"`
}

// New instantiate new client properties
func New(baseURL *url.URL) *Client {
	if baseURL == nil {
		baseURL = defaultBaseURL
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

var (
	defaultBaseURL     = &url.URL{Host: "projectid.app.jexia.com", Scheme: "https", Path: "/"}
	authRequest        = "auth"
	authRefreshRequest = "auth/refresh"
	ds                 = "ds/"
)

func (c *Client) post(projectID, request string, values []byte, refreshToken string) (cr *ClientResp, err error) {
	requestURL := strings.Replace(c.baseURL.String()+request, "projectid", projectID, -1)

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(values))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", refreshToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cr = &ClientResp{}
	if len(refreshToken) == 0 {
		_ = json.Unmarshal(body, cr)
		if err != nil {
			return nil, err
		}
	}
	cr.Success = true

	return cr, nil
}
