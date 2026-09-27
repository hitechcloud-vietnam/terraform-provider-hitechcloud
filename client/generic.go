// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Generic lenient helpers shared by the endpoint groups below. The API ships
// no response examples, so read endpoints return the unwrapped payload as
// map[string]any / []any and callers pick the fields they need.

// qmap converts a plain map into url.Values, skipping empty values.
func qmap(m map[string]string) url.Values {
	q := url.Values{}
	for k, v := range m {
		if v != "" {
			q.Set(k, v)
		}
	}
	return q
}

// getMap GETs path and returns the unwrapped payload as an object.
func (c *Client) getMap(ctx context.Context, path string, q url.Values) (map[string]any, error) {
	var raw any
	if err := c.Get(ctx, path, q, &raw); err != nil {
		return nil, err
	}
	return AsMapAny(raw), nil
}

// getList GETs path and returns the collection items.
func (c *Client) getList(ctx context.Context, path string, q url.Values) ([]any, error) {
	var raw any
	if err := c.Get(ctx, path, q, &raw); err != nil {
		return nil, err
	}
	return ExtractList(raw), nil
}

// postMap POSTs to path and returns the unwrapped payload (may be empty).
func (c *Client) postMap(ctx context.Context, path string, q url.Values) (map[string]any, error) {
	var raw any
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return nil, err
	}
	return AsMapAny(raw), nil
}

// putMap PUTs to path and returns the unwrapped payload (may be empty).
func (c *Client) putMap(ctx context.Context, path string, q url.Values) (map[string]any, error) {
	var raw any
	if err := c.Put(ctx, path, q, &raw); err != nil {
		return nil, err
	}
	return AsMapAny(raw), nil
}

// deleteMap DELETEs path and returns the unwrapped payload (may be empty).
func (c *Client) deleteMap(ctx context.Context, path string, q url.Values) (map[string]any, error) {
	var raw any
	if err := c.Delete(ctx, path, q, &raw); err != nil {
		return nil, err
	}
	return AsMapAny(raw), nil
}

// GetRaw GETs path and returns the raw response body, for binary or PEM
// downloads (certificate files, ticket attachments, ...).
func (c *Client) GetRaw(ctx context.Context, path string, q url.Values) ([]byte, error) {
	u := *c.baseURL
	u.Path = strings.TrimSuffix(u.Path, "/") + path
	if len(q) > 0 {
		u.RawQuery = q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", c.userAgent)
	if tok := c.TokenString(); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, newAPIError(resp.StatusCode, body)
	}
	return body, nil
}
