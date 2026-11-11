// Package kwai provides an unofficial SDK for the overseas Kwai Marketing API.
package kwai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://developers.kwai.com"
const DefaultTokenURL = "https://business.kwai.com/oauth/token"
const maxResponseBytes = 32 << 20

// Params preserves API field names, optional zero values, and nested models.
type Params map[string]any

type Config struct {
	AccessToken string
	BaseURL     string
	TokenURL    string
	HTTPClient  *http.Client
}

// Client is safe for concurrent calls. WithAccessToken returns a new client.
type Client struct {
	Business    *BusinessService
	baseURL     string
	tokenURL    string
	accessToken string
	http        *http.Client
	Reports     *ReportsService
	Campaigns   *CampaignsService
	Units       *UnitsService
	Creatives   *CreativesService
	Materials   *MaterialsService
	Organic     *OrganicService
	Finance     *FinanceService
}

func New(config Config) (*Client, error) {
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	if config.TokenURL == "" {
		config.TokenURL = DefaultTokenURL
	}
	if err := validateURL(config.BaseURL); err != nil {
		return nil, fmt.Errorf("kwai: invalid base URL: %w", err)
	}
	if err := validateURL(config.TokenURL); err != nil {
		return nil, fmt.Errorf("kwai: invalid token URL: %w", err)
	}
	c := &Client{baseURL: strings.TrimRight(config.BaseURL, "/"), tokenURL: config.TokenURL, accessToken: config.AccessToken, http: newHTTPClient(config.HTTPClient)}
	c.initServices()
	return c, nil
}

func validateURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" || (u.Scheme != "https" && u.Scheme != "http") || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return errors.New("expected an http(s) URL without credentials, query, or fragment")
	}
	return nil
}

func newHTTPClient(original *http.Client) *http.Client {
	c := http.Client{Timeout: 30 * time.Second}
	if original != nil {
		c = *original
	}
	// Never forward API credentials to a redirect target.
	c.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &c
}

func (c *Client) WithAccessToken(token string) *Client {
	next := *c
	next.accessToken = token
	next.initServices()
	return &next
}

// Response retains the full envelope and unmodified JSON data.
type Response struct {
	Status    int             `json:"status"`
	Message   string          `json:"message"`
	Data      json.RawMessage `json:"data"`
	RequestID string          `json:"requestId"`
	Raw       json.RawMessage `json:"-"`
}

// DecodeData uses json.Number when decoding into interface/map values so large IDs
// and monetary values are not silently rounded through float64.
func (r *Response) DecodeData(target any) error { return decodeJSON(r.Data, target) }

type APIError struct {
	HTTPStatus int
	Status     int
	Message    string
	RequestID  string
	Body       json.RawMessage
}

func (e *APIError) Error() string {
	return fmt.Sprintf("kwai: HTTP %d, status %d, requestId %s: %s", e.HTTPStatus, e.Status, e.RequestID, e.Message)
}

func decodeJSON(body []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return errors.New("kwai: trailing JSON content")
	}
	return nil
}

// Call performs a single POST. It never automatically retries a mutation.
func (c *Client) Call(ctx context.Context, path string, params any) (*Response, error) {
	if !strings.HasPrefix(path, "/rest/n/mapi/") || strings.ContainsAny(path, "?#\\") || strings.Contains(path, "..") {
		return nil, errors.New("kwai: expected an absolute /rest/n/mapi/ path")
	}
	if c.accessToken == "" {
		return nil, errors.New("kwai: access token is required")
	}
	if params == nil {
		params = Params{}
	}
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("kwai: encode request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return nil, errors.New("kwai: invalid request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", c.accessToken)
	body, status, err := execute(c.http, req)
	if err != nil {
		return nil, err
	}
	var response Response
	decodeErr := decodeJSON(body, &response)
	if status < 200 || status >= 300 {
		return nil, &APIError{status, response.Status, response.Message, response.RequestID, body}
	}
	if decodeErr != nil {
		return nil, fmt.Errorf("kwai: invalid response JSON: %w", decodeErr)
	}
	if response.Status != 200 {
		return nil, &APIError{status, response.Status, response.Message, response.RequestID, body}
	}
	response.Raw = body
	return &response, nil
}

func execute(client *http.Client, req *http.Request) ([]byte, int, error) {
	res, err := client.Do(req)
	if err != nil {
		if req.Context().Err() != nil {
			return nil, 0, req.Context().Err()
		}
		// Transport errors may include URLs containing OAuth/upload credentials.
		return nil, 0, errors.New("kwai: HTTP transport failed")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(io.LimitReader(res.Body, maxResponseBytes+1))
	if err != nil {
		return nil, res.StatusCode, errors.New("kwai: read response failed")
	}
	if len(body) > maxResponseBytes {
		return nil, res.StatusCode, errors.New("kwai: response exceeds 32 MiB")
	}
	return body, res.StatusCode, nil
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

// RefreshAccessToken explicitly refreshes credentials. Persist the returned refresh
// token, and use WithAccessToken to create a client with the new access token.
func (c *Client) RefreshAccessToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*Token, error) {
	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return nil, errors.New("kwai: client ID, secret and refresh token are required")
	}
	u, _ := url.Parse(c.tokenURL)
	q := u.Query()
	q.Set("grant_type", "refresh_token")
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("refresh_token", refreshToken)
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.New("kwai: invalid OAuth request")
	}
	body, status, err := execute(c.http, req)
	if err != nil {
		return nil, err
	}
	var token Token
	if status != 200 || decodeJSON(body, &token) != nil || token.AccessToken == "" || token.ExpiresIn <= 0 {
		return nil, fmt.Errorf("kwai: OAuth refresh failed (HTTP %d)", status)
	}
	return &token, nil
}
