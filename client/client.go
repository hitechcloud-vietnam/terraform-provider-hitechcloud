// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

// Package client implements a small, focused HTTP client for the HiTechCloud
// User API (the API described by the Postman collection shipped with this
// repository).
//
// Design notes, derived from the API specification:
//
//   - Authentication uses `Authorization: Bearer <token>` after obtaining a
//     token from `POST /api/login`.
//   - Every request parameter is sent as a query parameter; the API has no
//     JSON request bodies.
//   - The specification ships no response examples, so decoding is lenient:
//     responses may be wrapped in a common envelope (e.g. {"data": ...},
//     {"result": ...}) or be bare objects/arrays. The helpers in decode.go
//     normalise these shapes.
//   - The API has no idempotency-key support. To avoid duplicate resources,
//     unsafe methods (POST) are never retried automatically; idempotent
//     methods (GET/PUT/DELETE) are retried on transient failures (429, 5xx,
//     network errors) with exponential backoff.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultEndpoint is the production API base URL taken from the Postman
	// collection variable `baseUrl`.
	DefaultEndpoint = "https://api.hitechcloud.vn"

	// DefaultTimeout is the default per-request timeout.
	DefaultTimeout = 60 * time.Second

	// DefaultMaxRetries is the default number of retries for idempotent calls.
	DefaultMaxRetries = 3

	// DefaultRetryWaitMin/Max bound the exponential backoff.
	DefaultRetryWaitMin = 500 * time.Millisecond
	DefaultRetryWaitMax = 10 * time.Second

	defaultUserAgent = "terraform-provider-hitechcloud"
)

// Config configures a Client.
type Config struct {
	// Endpoint is the API base URL, e.g. "https://api.hitechcloud.vn".
	Endpoint string

	// Token is the bearer token obtained from POST /api/login. Required for
	// all authenticated endpoints. Never logged.
	Token string

	// RefreshToken, when set, is used to obtain a fresh access token via
	// POST /api/token whenever the current one expires (HTTP 401/403).
	RefreshToken string

	// UserAgent overrides the default User-Agent header.
	UserAgent string

	// Timeout is the per-request timeout (including body read).
	Timeout time.Duration

	// MaxRetries is the number of additional attempts for idempotent
	// requests that fail with a retryable error.
	MaxRetries int

	// RetryWaitMin / RetryWaitMax bound the exponential backoff duration.
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration

	// HTTPClient optionally overrides the underlying *http.Client. Useful
	// for tests and for advanced TLS configuration.
	HTTPClient *http.Client
}

// Client is a HiTechCloud User API client.
type Client struct {
	baseURL      *url.URL
	tokenMu      sync.RWMutex
	token        string
	refreshToken string
	userAgent    string
	http         *http.Client
	maxRetries   int
	retryWaitMin time.Duration
	retryWaitMax time.Duration

	// Logf, when non-nil, receives debug lines (never secrets). Intended to
	// be wired to a logger (e.g. terraform-plugin-log) in production.
	Logf func(format string, args ...any)
}

// New validates cfg and returns a ready-to-use Client.
func New(cfg Config) (*Client, error) {
	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint %q: %w", endpoint, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid endpoint %q: scheme must be http or https", endpoint)
	}
	// Normalise: strip trailing slash so path joins are predictable.
	u.Path = strings.TrimSuffix(u.Path, "/")

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: timeout}
	}

	maxRetries := cfg.MaxRetries
	if maxRetries < 0 {
		maxRetries = 0
	} else if cfg.MaxRetries == 0 {
		maxRetries = DefaultMaxRetries
	}

	waitMin := cfg.RetryWaitMin
	if waitMin <= 0 {
		waitMin = DefaultRetryWaitMin
	}
	waitMax := cfg.RetryWaitMax
	if waitMax <= 0 {
		waitMax = DefaultRetryWaitMax
	}

	userAgent := cfg.UserAgent
	if userAgent == "" {
		userAgent = defaultUserAgent
	}

	return &Client{
		baseURL:      u,
		token:        cfg.Token,
		refreshToken: cfg.RefreshToken,
		userAgent:    userAgent,
		http:         httpClient,
		maxRetries:   maxRetries,
		retryWaitMin: waitMin,
		retryWaitMax: waitMax,
	}, nil
}

// BaseURL returns the configured API base URL (without trailing slash).
func (c *Client) BaseURL() string {
	return c.baseURL.String()
}

// SetToken replaces the bearer token. Used by Login.
func (c *Client) SetToken(token string) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.token = token
}

// TokenString returns the current bearer token.
func (c *Client) TokenString() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token
}

// SetRefreshToken stores the refresh token used for automatic token renewal.
func (c *Client) SetRefreshToken(token string) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.refreshToken = token
}

// RefreshTokenString returns the current refresh token.
func (c *Client) RefreshTokenString() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.refreshToken
}

// Verb helpers ---------------------------------------------------------------

// Get performs a GET request and decodes the (unwrapped) JSON response into
// out when out is non-nil.
func (c *Client) Get(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodGet, path, query, out)
}

// Post performs a POST request. Per the API specification all inputs travel
// as query parameters.
func (c *Client) Post(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodPost, path, query, out)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodPut, path, query, out)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodDelete, path, query, out)
}

// do builds and executes a request with retries for idempotent methods. When
// the API reports an expired/invalid bearer token and a refresh token is
// available, a fresh access token is obtained via POST /api/token and the
// request is replayed once.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, out any) error {
	u := *c.baseURL
	// path always starts with "/".
	u.Path = strings.TrimSuffix(u.Path, "/") + path
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	retryable := method == http.MethodGet || method == http.MethodPut || method == http.MethodDelete
	attempts := 1
	if retryable {
		attempts += c.maxRetries
	}

	refreshed := false
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			if err := sleepCtx(ctx, c.backoff(attempt)); err != nil {
				return err
			}
			c.debug("retrying %s %s (attempt %d/%d): %v", method, redactURL(&u), attempt+1, attempts, lastErr)
		}

		retryableErr, err := c.doOnce(ctx, method, &u, out)
		if err == nil {
			return nil
		}
		lastErr = err

		// Access token expired or was revoked: renew it once with the refresh
		// token (POST /api/token) and replay the request immediately. The
		// auth endpoints themselves never trigger another refresh cycle.
		if !refreshed && isAuthError(err) && c.RefreshTokenString() != "" && !isAuthPath(path) {
			if _, rerr := c.RefreshToken(ctx, c.RefreshTokenString()); rerr == nil {
				refreshed = true
				attempt-- // do not consume a retry slot
				continue
			}
		}

		if !retryable || !retryableErr {
			return err
		}
	}
	return fmt.Errorf("request %s %s failed after %d attempt(s): %w", method, redactURL(&u), attempts, lastErr)
}

// isAuthError reports whether err indicates an expired or invalid bearer
// token (HTTP 401/403).
func isAuthError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized || apiErr.StatusCode == http.StatusForbidden
	}
	return false
}

// isAuthPath reports whether path belongs to the authentication endpoints
// (/api/login, /api/token, /api/logout, /api/revoke) which must never
// re-enter the refresh flow.
func isAuthPath(path string) bool {
	switch path {
	case "/api/login", "/api/token", "/api/logout", "/api/revoke":
		return true
	}
	return false
}

// doOnce executes a single HTTP request. It returns whether the error is
// retryable alongside the error itself.
func (c *Client) doOnce(ctx context.Context, method string, u *url.URL, out any) (retryable bool, err error) {
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if tok := c.TokenString(); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	c.debug("%s %s", method, redactURL(u))
	resp, err := c.http.Do(req)
	if err != nil {
		// Network level errors are safe to retry for idempotent methods.
		return true, fmt.Errorf("performing request: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
	}()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20)) // 32 MiB cap
	if err != nil {
		return true, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := newAPIError(resp.StatusCode, body)
		return isRetryableStatus(resp.StatusCode), apiErr
	}

	// Some API failures are reported with HTTP 200 and a failure envelope.
	if apiErr := envelopeError(resp.StatusCode, body); apiErr != nil {
		return false, apiErr
	}

	if out == nil || len(body) == 0 {
		return false, nil
	}

	payload, err := unwrapJSON(body)
	if err != nil {
		return false, fmt.Errorf("decoding response from %s: %w", redactURL(u), err)
	}

	// Re-marshal the normalised payload into the caller's type. Using
	// map[string]any round-tripping keeps decoding lenient about envelope
	// shapes while still supporting typed structs.
	raw, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("decoding response from %s: %w", redactURL(u), err)
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return false, fmt.Errorf("decoding response from %s: %w", redactURL(u), err)
	}
	return false, nil
}

// backoff computes an exponential backoff with jitter for the given attempt.
func (c *Client) backoff(attempt int) time.Duration {
	d := c.retryWaitMin << (attempt - 1)
	if d > c.retryWaitMax || d <= 0 {
		d = c.retryWaitMax
	}
	// Full jitter: random duration in (0, d].
	return time.Duration(rand.Int63n(int64(d))) + time.Millisecond
}

func (c *Client) debug(format string, args ...any) {
	if c.Logf != nil {
		c.Logf(format, args...)
	}
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// Query builds url.Values from alternating key/value pairs. Empty values are
// skipped so optional parameters are simply omitted from the request.
func Query(kv ...string) url.Values {
	q := url.Values{}
	for i := 0; i+1 < len(kv); i += 2 {
		if kv[i+1] == "" {
			continue
		}
		q.Set(kv[i], kv[i+1])
	}
	return q
}

// QueryBool adds a boolean parameter only when v is set.
func QueryBool(q url.Values, key string, v *bool) {
	if v != nil {
		q.Set(key, strconv.FormatBool(*v))
	}
}

// JSONQuery marshals v to JSON and sets it as a query parameter. Used for the
// few parameters documented as objects/arrays (e.g. `tags`, `envs`,
// `launch_configuration`). Empty JSON documents are skipped.
func JSONQuery(q url.Values, key string, v any) error {
	if v == nil {
		return nil
	}
	if rm, ok := v.(json.RawMessage); ok {
		if len(bytes.TrimSpace(rm)) == 0 {
			return nil
		}
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("encoding %s: %w", key, err)
	}
	if string(raw) == "null" || string(raw) == "{}" || string(raw) == "[]" {
		return nil
	}
	q.Set(key, string(raw))
	return nil
}

// PathEscape escapes a single path segment.
func PathEscape(s string) string {
	return url.PathEscape(s)
}

// redactURL strips query values from a URL for safe logging.
func redactURL(u *url.URL) string {
	copied := *u
	if copied.RawQuery != "" {
		copied.RawQuery = "REDACTED"
	}
	return copied.String()
}

// IsNotFound reports whether err indicates a missing API object (HTTP 404).
func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}
