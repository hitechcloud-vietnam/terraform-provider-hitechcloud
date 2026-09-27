// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// APIError normalises error responses coming from the HiTechCloud API.
type APIError struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	// Raw is the (truncated) raw response body, useful for debugging.
	Raw string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("HiTechCloud API error: %s (HTTP %d)", e.Status, e.StatusCode)
	}
	return fmt.Sprintf("HiTechCloud API error: %s (HTTP %d)", e.Message, e.StatusCode)
}

// newAPIError builds an APIError from an HTTP response.
func newAPIError(statusCode int, body []byte) *APIError {
	msg := extractErrorMessage(body)
	if msg == "" {
		msg = http.StatusText(statusCode)
	}
	return &APIError{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Message:    msg,
		Raw:        truncate(string(body), 2048),
	}
}

// envelopeError detects failure envelopes delivered with HTTP 2xx status, e.g.
// {"success": false, "message": "..."} or {"error": "..."}.
func envelopeError(statusCode int, body []byte) *APIError {
	if len(body) == 0 {
		return nil
	}
	var raw any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return nil
	}

	if success, ok := m["success"].(bool); ok && !success {
		return &APIError{
			StatusCode: statusCode,
			Status:     http.StatusText(statusCode),
			Message:    extractErrorMessage(body),
			Raw:        truncate(string(body), 2048),
		}
	}
	if errVal, ok := m["error"]; ok && errVal != nil {
		if s := asString(errVal); s != "" && !strings.EqualFold(s, "null") && !strings.EqualFold(s, "false") {
			return &APIError{
				StatusCode: statusCode,
				Status:     http.StatusText(statusCode),
				Message:    s,
				Raw:        truncate(string(body), 2048),
			}
		}
	}
	return nil
}

// extractErrorMessage pulls a human readable message out of a JSON error body.
func extractErrorMessage(body []byte) string {
	var raw any
	if err := json.Unmarshal(body, &raw); err != nil {
		return truncate(strings.TrimSpace(string(body)), 512)
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return truncate(strings.TrimSpace(string(body)), 512)
	}
	for _, key := range []string{"message", "error_description", "error", "detail", "errors", "msg"} {
		if v, ok := m[key]; ok && v != nil {
			switch t := v.(type) {
			case string:
				if t != "" {
					return t
				}
			case []any:
				parts := make([]string, 0, len(t))
				for _, item := range t {
					parts = append(parts, asString(item))
				}
				if joined := strings.Join(parts, "; "); joined != "" {
					return joined
				}
			default:
				if s := asString(v); s != "" {
					return s
				}
			}
		}
	}
	return ""
}

// isRetryableStatus reports whether an HTTP status is worth retrying.
func isRetryableStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusRequestTimeout,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return statusCode >= 500
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
