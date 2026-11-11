package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-acme/lego/v5/internal/errutils"
	"github.com/go-acme/lego/v5/internal/useragent"
)

// DefaultBaseURL the default API endpoint.
const DefaultBaseURL = "https://www.gyrohosting.com/api/v1/client/"

// Client the Gyro DNS API client.
type Client struct {
	token string

	baseURL    *url.URL
	HTTPClient *http.Client
}

// NewClient creates a new Client.
func NewClient(baseURL, token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("credentials missing")
	}

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}

	return &Client{
		token:      token,
		baseURL:    base,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// AddChallenge adds a value to the TXT record set of an _acme-challenge name.
// Other values of the same name are kept (wildcard and apex challenges can coexist).
// https://www.gyrohosting.com/developers
func (c *Client) AddChallenge(ctx context.Context, name, value string) (*Challenge, error) {
	endpoint := c.baseURL.JoinPath("cloud", "dns", "acme-challenge")

	req, err := newJSONRequest(ctx, http.MethodPost, endpoint, ChallengeRequest{Name: name, Value: value})
	if err != nil {
		return nil, err
	}

	var result ChallengeResponse

	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// RemoveChallenge removes a value from the TXT record set of an _acme-challenge name.
// https://www.gyrohosting.com/developers
func (c *Client) RemoveChallenge(ctx context.Context, name, value string) (*Challenge, error) {
	endpoint := c.baseURL.JoinPath("cloud", "dns", "acme-challenge")

	query := endpoint.Query()
	query.Set("name", name)
	query.Set("value", value)
	endpoint.RawQuery = query.Encode()

	req, err := newJSONRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var result ChallengeResponse

	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) do(req *http.Request, result any) error {
	useragent.SetHeader(req.Header)

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return errutils.NewHTTPDoError(req, err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return parseError(req, resp)
	}

	if result == nil {
		return nil
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return errutils.NewReadResponseError(req, resp.StatusCode, err)
	}

	err = json.Unmarshal(raw, result)
	if err != nil {
		return errutils.NewUnmarshalError(req, resp.StatusCode, raw, err)
	}

	return nil
}

func newJSONRequest(ctx context.Context, method string, endpoint *url.URL, payload any) (*http.Request, error) {
	buf := new(bytes.Buffer)

	if payload != nil {
		err := json.NewEncoder(buf).Encode(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to create request JSON body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func parseError(req *http.Request, resp *http.Response) error {
	raw, _ := io.ReadAll(resp.Body)

	errAPI := new(APIError)

	err := json.Unmarshal(raw, errAPI)
	if err != nil || errAPI.ErrorInfo.Code == "" {
		return errutils.NewUnexpectedStatusCodeError(req, resp.StatusCode, raw)
	}

	return errAPI
}
