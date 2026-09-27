// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
)

// ---------------------------------------------------------------------------
// Ceph S3 object storage (folder "Ceph S3": /api/service/:id/s3/...)
// ---------------------------------------------------------------------------

// S3Bucket is a bucket of the S3 service.
type S3Bucket struct {
	Name      string
	Created   string
	SizeBytes int64
	Objects   int64
}

// ListS3Buckets lists buckets (GET /api/service/:id/s3/buckets).
func (c *Client) ListS3Buckets(ctx context.Context, serviceID string) ([]S3Bucket, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/s3/buckets"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "buckets")
	out := make([]S3Bucket, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseS3Bucket(m))
		} else if s := asString(item); s != "" {
			out = append(out, S3Bucket{Name: s})
		}
	}
	return out, nil
}

// GetS3Bucket fetches one bucket by name from the bucket listing. The API has
// no single-bucket read endpoint.
func (c *Client) GetS3Bucket(ctx context.Context, serviceID, bucket string) (*S3Bucket, error) {
	buckets, err := c.ListS3Buckets(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	for _, b := range buckets {
		if b.Name == bucket {
			found := b
			return &found, nil
		}
	}
	return nil, &APIError{StatusCode: 404, Status: "Not Found", Message: fmt.Sprintf("bucket %s not found", bucket)}
}

// CreateS3Bucket creates a bucket (POST /api/service/:id/s3/buckets).
//
// The API normalises the name to lower case and may prepend a service prefix
// when "Auto-prefix Buckets" is enabled, so the authoritative name is read
// from the response (`bucket` field) as the API documentation instructs.
func (c *Client) CreateS3Bucket(ctx context.Context, serviceID, name string) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/s3/buckets"
	if err := c.Post(ctx, path, Query("name", name), &raw); err != nil {
		return "", err
	}
	if bucket := FirstString(raw, "bucket", "name", "bucket_name"); bucket != "" {
		return bucket, nil
	}
	// Fallback: locate the bucket in the listing.
	buckets, err := c.ListS3Buckets(ctx, serviceID)
	if err != nil {
		return "", fmt.Errorf("bucket was created but its name could not be resolved: %w", err)
	}
	for _, b := range buckets {
		if b.Name == name {
			return b.Name, nil
		}
	}
	return "", fmt.Errorf("bucket %q was created but could not be found in the bucket list", name)
}

// DeleteS3Bucket removes a bucket (DELETE /api/service/:id/s3/buckets/:bucket).
// With purge=true the contents are deleted as well; the API refuses to remove
// non-empty buckets otherwise.
func (c *Client) DeleteS3Bucket(ctx context.Context, serviceID, bucket string, purge bool) error {
	path := "/api/service/" + PathEscape(serviceID) + "/s3/buckets/" + PathEscape(bucket)
	q := Query("purge", "0")
	if purge {
		q.Set("purge", "1")
	}
	return c.Delete(ctx, path, q, nil)
}

func parseS3Bucket(m map[string]any) S3Bucket {
	return S3Bucket{
		Name:      FirstString(m, "bucket", "name", "bucket_name"),
		Created:   FirstString(m, "created", "created_at", "creation_date"),
		SizeBytes: FirstInt64(m, "size_bytes", "size", "bytes", "size_actual"),
		Objects:   FirstInt64(m, "objects", "object_count", "num_objects"),
	}
}

// S3Subuser is a sub-user with its own S3 key pair.
type S3Subuser struct {
	Name      string
	Access    string
	AccessKey string
	SecretKey string // returned by the API exactly once, at creation
}

// ListS3Subusers lists sub-users (GET /api/service/:id/s3/subusers). Secret
// keys are never included in listings by design of the API.
func (c *Client) ListS3Subusers(ctx context.Context, serviceID string) ([]S3Subuser, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/s3/subusers"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "subusers", "sub_users", "users")
	out := make([]S3Subuser, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseS3Subuser(m))
		}
	}
	return out, nil
}

// GetS3Subuser fetches one sub-user from the listing.
func (c *Client) GetS3Subuser(ctx context.Context, serviceID, name string) (*S3Subuser, error) {
	users, err := c.ListS3Subusers(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Name == name {
			found := u
			return &found, nil
		}
	}
	return nil, &APIError{StatusCode: 404, Status: "Not Found", Message: fmt.Sprintf("sub-user %s not found", name)}
}

// CreateS3Subuser creates a sub-user (POST /api/service/:id/s3/subusers).
// access must be one of read | write | readwrite | full. The returned secret
// key is only ever provided in this response.
func (c *Client) CreateS3Subuser(ctx context.Context, serviceID, name, access string) (*S3Subuser, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/s3/subusers"
	if err := c.Post(ctx, path, Query("name", name, "access", access), &raw); err != nil {
		return nil, err
	}
	sub := S3Subuser{Name: name, Access: access}
	if m, ok := AsMap(raw); ok {
		sub = parseS3Subuser(m)
	}
	if sub.Name == "" {
		sub.Name = name
	}
	if sub.Access == "" {
		sub.Access = access
	}
	return &sub, nil
}

// DeleteS3Subuser removes a sub-user and its keys
// (DELETE /api/service/:id/s3/subusers/:subuser). The API accepts both `name`
// and `uid:name` forms.
func (c *Client) DeleteS3Subuser(ctx context.Context, serviceID, subuser string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/s3/subusers/" + PathEscape(subuser)
	return c.Delete(ctx, path, nil, nil)
}

func parseS3Subuser(m map[string]any) S3Subuser {
	return S3Subuser{
		Name:      FirstString(m, "name", "subuser", "user", "uid_name"),
		Access:    FirstString(m, "access", "permissions", "perm", "access_level"),
		AccessKey: FirstString(m, "access_key", "accesskey", "access_key_id"),
		SecretKey: FirstString(m, "secret_key", "secretkey", "secret", "secret_access_key"),
	}
}
