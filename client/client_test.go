// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"
)

// newTestClient returns a client pointed at the given mock server.
func newTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	cli, err := New(Config{
		Endpoint:     srv.URL,
		Token:        "test-token",
		Timeout:      5 * time.Second,
		MaxRetries:   3,
		RetryWaitMin: time.Millisecond,
		RetryWaitMax: 2 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	return cli
}

func TestLoginSetsBearerToken(t *testing.T) {
	var gotAuth, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if r.URL.Path == "/api/login" {
			_ = json.NewEncoder(w).Encode(map[string]any{"token": "abc123"})
			return
		}
		gotAuth = r.Header.Get("Authorization")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	res, err := cli.Login(context.Background(), "user@example.com", "secret")
	if err != nil {
		t.Fatalf("Login() error: %v", err)
	}
	if res.Token != "abc123" {
		t.Fatalf("Login() token = %q, want abc123", res.Token)
	}
	if gotPath != "/api/login" {
		t.Fatalf("login path = %q, want /api/login", gotPath)
	}

	// A subsequent request must carry the bearer token stored by Login
	// (the login response token replaces the configured one).
	if err := cli.Get(context.Background(), "/api/details", nil, nil); err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if gotAuth != "Bearer abc123" {
		t.Fatalf("Authorization = %q, want Bearer abc123", gotAuth)
	}
}

func TestLoginTokenFallbackKeys(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "nested", "user": map[string]any{"token": "deep"}})
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	res, err := cli.Login(context.Background(), "u", "p")
	if err != nil {
		t.Fatalf("Login() error: %v", err)
	}
	if res.Token == "" {
		t.Fatal("Login() token is empty, want fallback key resolution")
	}
}

func TestErrorNormalization(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"nope"}`))
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	err := cli.Get(context.Background(), "/api/missing", nil, nil)
	if err == nil {
		t.Fatal("Get() expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false, want true", err)
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
	if apiErr.Message != "nope" {
		t.Fatalf("Message = %q, want nope", apiErr.Message)
	}
}

func TestEnvelopeFailureOnHTTP200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":false,"error":"quota exceeded"}`))
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	err := cli.Get(context.Background(), "/api/anything", nil, nil)
	if err == nil {
		t.Fatal("Get() expected envelope error, got nil")
	}
	if want := "quota exceeded"; !containsStr(err.Error(), want) {
		t.Fatalf("error %q does not contain %q", err.Error(), want)
	}
}

func TestRetriesIdempotentMethods(t *testing.T) {
	var gets int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&gets, 1) <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	if err := cli.Get(context.Background(), "/api/flaky", nil, nil); err != nil {
		t.Fatalf("Get() error after retries: %v", err)
	}
	if n := atomic.LoadInt32(&gets); n != 3 {
		t.Fatalf("attempts = %d, want 3", n)
	}
}

func TestDoesNotRetryPostCreates(t *testing.T) {
	var posts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&posts, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	err := cli.Post(context.Background(), "/api/service/1/dns", Query("name", "example.com"), nil)
	if err == nil {
		t.Fatal("Post() expected error, got nil")
	}
	if n := atomic.LoadInt32(&posts); n != 1 {
		t.Fatalf("POST attempts = %d, want 1 (creates must not be retried)", n)
	}
}

func TestRetryOnServerErrorForDelete(t *testing.T) {
	var dels int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&dels, 1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	if err := cli.Delete(context.Background(), "/api/thing/1", nil, nil); err != nil {
		t.Fatalf("Delete() error after retry: %v", err)
	}
	if n := atomic.LoadInt32(&dels); n != 2 {
		t.Fatalf("DELETE attempts = %d, want 2", n)
	}
}

func TestEnvelopeUnwrapAndExtractList(t *testing.T) {
	raw := []byte(`{"data":{"items":[{"id":"1"},{"id":"2"}]}}`)
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	m, ok := AsMap(v)
	if !ok {
		t.Fatal("AsMap returned false")
	}
	if got := FirstString(m, "id"); got != "" {
		t.Fatalf("FirstString(id) = %q, want empty", got)
	}
	list := ExtractList(v, "items")
	if len(list) != 2 {
		t.Fatalf("ExtractList len = %d, want 2", len(list))
	}
	item, ok := AsMap(list[0])
	if !ok {
		t.Fatal("AsMap(list[0]) failed")
	}
	if got := FirstInt64(item, "id"); got != 1 {
		t.Fatalf("FirstInt64 = %d, want 1", got)
	}
}

func TestFirstFallbacks(t *testing.T) {
	m := map[string]any{
		"first_name": "Ada",
		"count":      float64(3),
		"active":     "true",
	}
	if got := FirstString(m, "name", "first_name"); got != "Ada" {
		t.Fatalf("FirstString = %q, want Ada", got)
	}
	if got := FirstInt64(m, "count"); got != 3 {
		t.Fatalf("FirstInt64 = %d, want 3", got)
	}
	if got := FirstBool(m, "active"); !got {
		t.Fatal("FirstBool = false, want true")
	}
}

func TestQueryHelpers(t *testing.T) {
	q := Query("a", "1", "b", "", "c", "x")
	if q.Get("a") != "1" || q.Get("c") != "x" {
		t.Fatalf("Query values = %v", q)
	}
	if _, ok := q["b"]; ok {
		t.Fatal("empty query value should be skipped")
	}

	if err := JSONQuery(q, "tags", []string{"x", "y"}); err != nil {
		t.Fatalf("JSONQuery: %v", err)
	}
	if q.Get("tags") != `["x","y"]` {
		t.Fatalf("tags = %q", q.Get("tags"))
	}
	// Empty documents are skipped entirely.
	if err := JSONQuery(q, "envs", map[string]string{}); err != nil {
		t.Fatalf("JSONQuery: %v", err)
	}
	if _, ok := q["envs"]; ok {
		t.Fatal("empty JSON object should be skipped")
	}
}

func TestFormatAndSplitID(t *testing.T) {
	id := FormatID("svc", "zone", "rec")
	if id != "svc/zone/rec" {
		t.Fatalf("FormatID = %q", id)
	}
	parts, err := SplitID(id, 3)
	if err != nil {
		t.Fatalf("SplitID error: %v", err)
	}
	if parts[0] != "svc" || parts[1] != "zone" || parts[2] != "rec" {
		t.Fatalf("SplitID parts = %v", parts)
	}
	if _, err := SplitID("only-one", 3); err == nil {
		t.Fatal("SplitID expected error for wrong part count")
	}
}

func TestStringListAndMap(t *testing.T) {
	if got := StringList([]any{"a", "b"}); len(got) != 2 {
		t.Fatalf("StringList = %v", got)
	}
	// Key=value maps are supported.
	if got := StringList(map[string]any{"A": "1"}); len(got) != 1 {
		t.Fatalf("StringList(map) = %v", got)
	}
	if got := StringMap(map[string]any{"K": "v"}); got["K"] != "v" {
		t.Fatalf("StringMap = %v", got)
	}
}

func TestGetAccountParsesDetails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/details" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":          "42",
				"email":       "user@example.com",
				"firstname":   "Thanh",
				"lastname":    "An",
				"companyname": "HiTechCloud",
			},
		})
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	acc, err := cli.GetAccount(context.Background())
	if err != nil {
		t.Fatalf("GetAccount() error: %v", err)
	}
	if acc.ID != "42" || acc.Email != "user@example.com" || acc.FirstName != "Thanh" || acc.LastName != "An" {
		t.Fatalf("account = %+v", acc)
	}
}

func TestCreateDNSZoneResolvesIDByFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/service/1/dns":
			// Creation response without an id.
			_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
		case r.Method == http.MethodGet && r.URL.Path == "/api/service/1/dns":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"zones": []any{
					map[string]any{"id": "z1", "name": "example.com"},
				},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	id, err := cli.CreateDNSZone(context.Background(), "1", "example.com")
	if err != nil {
		t.Fatalf("CreateDNSZone() error: %v", err)
	}
	if id != "z1" {
		t.Fatalf("id = %q, want z1 (fallback resolution by name)", id)
	}
}

func TestCreateAISendsJSONFieldsAndBools(t *testing.T) {
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotQuery = r.URL.Query()
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "i-1"})
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	autoDelete := true
	id, err := cli.CreateAIInstance(context.Background(), "1", AIInstanceCreate{
		Name:                "gpu-1",
		Cloud:               "shade",
		Region:              "hanoi",
		ShadeInstanceType:   "A100",
		Tags:                []string{"prod"},
		Envs:                map[string]string{"KEY": "VAL"},
		AutoDelete:          &autoDelete,
		LaunchConfiguration: `{"cmd":"train"}`,
	})
	if err != nil {
		t.Fatalf("CreateAIInstance() error: %v", err)
	}
	if id != "i-1" {
		t.Fatalf("id = %q", id)
	}
	if got := gotQuery.Get("tags"); got != `["prod"]` {
		t.Fatalf("tags query = %q", got)
	}
	if got := gotQuery.Get("envs"); got != `{"KEY":"VAL"}` {
		t.Fatalf("envs query = %q", got)
	}
	if got := gotQuery.Get("auto_delete"); got != "true" {
		t.Fatalf("auto_delete query = %q", got)
	}
	if got := gotQuery.Get("launch_configuration"); got != `{"cmd":"train"}` {
		t.Fatalf("launch_configuration query = %q", got)
	}
}

func TestCreateS3BucketUsesAuthoritativeName(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The API auto-prefixes the bucket name.
		_ = json.NewEncoder(w).Encode(map[string]any{"bucket": "123-photos"})
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	name, err := cli.CreateS3Bucket(context.Background(), "1", "photos")
	if err != nil {
		t.Fatalf("CreateS3Bucket() error: %v", err)
	}
	if name != "123-photos" {
		t.Fatalf("bucket = %q, want 123-photos (authoritative response field)", name)
	}
}

func TestWaitForVMDeletedTimesOut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "vm1", "status": "running"})
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := cli.WaitForVMDeleted(ctx, "1", "vm1", 50*time.Millisecond)
	if err == nil {
		t.Fatal("WaitForVMDeleted() expected timeout error")
	}
}

func TestAPIErrorTruncatesBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, string(make([]byte, 5000)))))
	}))
	defer srv.Close()

	cli := newTestClient(t, srv)
	err := cli.Get(context.Background(), "/api/bad", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if len(err.Error()) > 4096 {
		t.Fatalf("error message not truncated: %d bytes", len(err.Error()))
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
