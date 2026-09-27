// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// The HiTechCloud API ships without response examples, therefore response
// decoding must tolerate multiple plausible shapes:
//
//   - bare objects:      {"id": "1", "name": "web"}
//   - bare arrays:       [{"id": "1"}, ...]
//   - wrapped envelopes: {"success": true, "data": {...}} / {"result": {...}}
//   - keyed collections: {"domains": [...]} / {"records": [...]}
//
// The helpers below normalise those shapes so callers can rely on plain Go
// maps and slices.

// unwrapJSON decodes raw JSON and recursively unwraps the common envelope
// keys used by HostBill-style APIs.
func unwrapJSON(raw []byte) (any, error) {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	return unwrapValue(v), nil
}

var envelopeKeys = []string{"data", "result", "results", "details", "response", "content", "object"}

func unwrapValue(v any) any {
	m, ok := v.(map[string]any)
	if !ok {
		return v
	}
	for _, key := range envelopeKeys {
		if inner, ok := m[key]; ok && inner != nil {
			switch inner.(type) {
			case map[string]any, []any:
				return unwrapValue(inner)
			}
		}
	}
	return m
}

// AsMap coerces v into a map[string]any.
func AsMap(v any) (map[string]any, bool) {
	m, ok := v.(map[string]any)
	return m, ok
}

// AsSlice coerces v into a []any.
func AsSlice(v any) ([]any, bool) {
	s, ok := v.([]any)
	return s, ok
}

// Dig walks a (possibly nested) value using dotted path notation, e.g.
// Dig(v, "data.user.id"). Missing or non-object intermediates yield nil.
func Dig(v any, path string) any {
	current := v
	for _, part := range strings.Split(path, ".") {
		m, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current, ok = m[part]
		if !ok {
			return nil
		}
	}
	return current
}

// First digs several candidate paths and returns the first non-nil value.
func First(v any, paths ...string) any {
	for _, p := range paths {
		if got := Dig(v, p); got != nil {
			return got
		}
	}
	return nil
}

// FirstString digs several candidate keys and returns the first non-empty
// string representation.
func FirstString(v any, keys ...string) string {
	for _, key := range keys {
		if s := asString(Dig(v, key)); s != "" {
			return s
		}
	}
	return ""
}

// FirstInt64 digs several candidate keys and returns the first value that can
// be interpreted as an integer.
func FirstInt64(v any, keys ...string) int64 {
	for _, key := range keys {
		if n, ok := asInt64(Dig(v, key)); ok {
			return n
		}
	}
	return 0
}

// FirstBool digs several candidate keys and returns the first value that can
// be interpreted as a boolean.
func FirstBool(v any, keys ...string) bool {
	for _, key := range keys {
		if b, ok := asBool(Dig(v, key)); ok {
			return b
		}
	}
	return false
}

// FirstFloat64 digs several candidate keys and returns the first value that
// can be interpreted as a float.
func FirstFloat64(v any, keys ...string) float64 {
	for _, key := range keys {
		if f, ok := asFloat64(Dig(v, key)); ok {
			return f
		}
	}
	return 0
}

// ExtractList finds a list inside a decoded payload. It understands bare
// arrays as well as common collection keys.
func ExtractList(v any, keys ...string) []any {
	switch t := v.(type) {
	case []any:
		return t
	case map[string]any:
		candidates := append(append([]string{}, keys...),
			"items", "list", "results", "rows", "entries", "elements")
		candidates = append(candidates, envelopeKeys...)
		for _, key := range candidates {
			if inner, ok := t[key]; ok {
				if s, ok := inner.([]any); ok {
					return s
				}
			}
		}
		// Recurse into envelope objects, e.g. {"data": {"items": [...]}}.
		for _, key := range envelopeKeys {
			if inner, ok := t[key]; ok {
				if m, ok := inner.(map[string]any); ok {
					if s := ExtractList(m, keys...); s != nil {
						return s
					}
				}
			}
		}
	}
	return nil
}

// FindByID returns the first element of list whose id matches id, probing the
// given id keys (falling back to "id").
func FindByID(list []any, id string, idKeys ...string) (map[string]any, bool) {
	keys := append([]string{"id"}, idKeys...)
	for _, item := range list {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		for _, key := range keys {
			if asString(m[key]) == id {
				return m, true
			}
		}
	}
	return nil, false
}

// StringList converts a decoded array into []string, tolerating scalar values
// and map entries ("key=value") which some APIs return for tag-like fields.
func StringList(v any) []string {
	switch t := v.(type) {
	case []any:
		out := make([]string, 0, len(t))
		for _, item := range t {
			switch elem := item.(type) {
			case map[string]any:
				// e.g. {"key": "env", "value": "prod"} or {"name": "..."}
				if s := asString(elem["value"]); s != "" {
					out = append(out, asString(elem["key"])+"="+s)
				} else {
					out = append(out, asString(elem["name"]))
				}
			default:
				if s := asString(item); s != "" {
					out = append(out, s)
				}
			}
		}
		return out
	case map[string]any:
		out := make([]string, 0, len(t))
		for k, val := range t {
			out = append(out, k+"="+asString(val))
		}
		return out
	case nil:
		return nil
	default:
		if s := asString(v); s != "" {
			return []string{s}
		}
		return nil
	}
}

// StringMap converts a decoded object into map[string]string. Arrays of
// key/value objects are also understood.
func StringMap(v any) map[string]string {
	out := map[string]string{}
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			out[k] = asString(val)
		}
	case []any:
		for _, item := range t {
			if m, ok := item.(map[string]any); ok {
				key := FirstString(m, "key", "name")
				if key != "" {
					out[key] = asString(First(m, "value", "val"))
				}
			}
		}
	}
	return out
}

// asString renders a JSON value as string.
func asString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		return strconv.FormatBool(t)
	case float64:
		// Avoid "1.000000" for integral values.
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case json.Number:
		return t.String()
	default:
		raw, err := json.Marshal(t)
		if err != nil {
			return ""
		}
		return string(raw)
	}
}

func asInt64(v any) (int64, bool) {
	switch t := v.(type) {
	case float64:
		return int64(t), true
	case int:
		return int64(t), true
	case int64:
		return t, true
	case json.Number:
		n, err := t.Int64()
		return n, err == nil
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(t), 10, 64)
		return n, err == nil
	case bool:
		if t {
			return 1, true
		}
		return 0, true
	default:
		return 0, false
	}
}

func asFloat64(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case json.Number:
		f, err := t.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func asBool(v any) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case float64:
		return t != 0, true
	case int:
		return t != 0, true
	case int64:
		return t != 0, true
	case string:
		s := strings.TrimSpace(strings.ToLower(t))
		switch s {
		case "true", "1", "yes", "on", "enabled":
			return true, true
		case "false", "0", "no", "off", "disabled":
			return false, true
		}
		return false, false
	default:
		return false, false
	}
}

// FormatID joins composite identifier parts.
func FormatID(parts ...string) string {
	return strings.Join(parts, "/")
}

// SplitID splits a composite identifier into exactly n parts.
func SplitID(id string, n int) ([]string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != n {
		return nil, fmt.Errorf("unexpected ID format %q: want %d slash-separated parts", id, n)
	}
	for _, p := range parts {
		if p == "" {
			return nil, fmt.Errorf("unexpected ID format %q: empty part", id)
		}
	}
	return parts, nil
}
