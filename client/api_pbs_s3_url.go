// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// Proxmox Backup Server (folder "Proxmox Backup").

// GetPBSConnectionInfo returns the PBS connection info
// (GET /api/service/:id/pbs). Contains no secrets.
func (c *Client) GetPBSConnectionInfo(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/pbs", nil)
}

// GetPBSCredentials returns the PBS credentials including secrets
// (GET /api/service/:id/pbs/credentials).
func (c *Client) GetPBSCredentials(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/pbs/credentials", nil)
}

// GetPBSUsage returns backup usage (GET /api/service/:id/pbs/usage).
func (c *Client) GetPBSUsage(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/pbs/usage", nil)
}

// GetPBSMetrics returns backup metrics (GET /api/service/:id/pbs/metrics).
func (c *Client) GetPBSMetrics(ctx context.Context, serviceID, keys string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/pbs/metrics", Query("keys", keys))
}

// ListPBSSnapshots lists PBS snapshots (GET /api/service/:id/pbs/snapshots).
func (c *Client) ListPBSSnapshots(ctx context.Context, serviceID, limit string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/pbs/snapshots", Query("limit", limit))
}

// ListPBSGroups lists backup groups (GET /api/service/:id/pbs/groups).
func (c *Client) ListPBSGroups(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/pbs/groups", nil)
}

// ChangePBSPassword changes the PBS user password
// (POST /api/service/:id/pbs/password).
func (c *Client) ChangePBSPassword(ctx context.Context, serviceID, password string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/pbs/password", Query("password", password))
	return err
}

// RotatePBSToken rotates the PBS API token and returns the new secret
// (POST /api/service/:id/pbs/token).
func (c *Client) RotatePBSToken(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/pbs/token", nil)
}

// RevokePBSToken revokes the PBS API token
// (DELETE /api/service/:id/pbs/token).
func (c *Client) RevokePBSToken(ctx context.Context, serviceID string) error {
	_, err := c.deleteMap(ctx, "/api/service/"+PathEscape(serviceID)+"/pbs/token", nil)
	return err
}

// Ceph S3 extras (folder "Ceph S3").

// GetS3ConnectionInfo returns the S3 connection info
// (GET /api/service/:id/s3). Contains no secrets.
func (c *Client) GetS3ConnectionInfo(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/s3", nil)
}

// GetS3Credentials returns the S3 credentials including secrets
// (GET /api/service/:id/s3/credentials).
func (c *Client) GetS3Credentials(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/s3/credentials", nil)
}

// GetS3Usage returns S3 usage (GET /api/service/:id/s3/usage).
func (c *Client) GetS3Usage(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/s3/usage", nil)
}

// GetS3Metrics returns S3 metrics (GET /api/service/:id/s3/metrics).
func (c *Client) GetS3Metrics(ctx context.Context, serviceID, keys string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/s3/metrics", Query("keys", keys))
}

// RotateS3SecretKey rotates the S3 secret key
// (POST /api/service/:id/s3/key).
func (c *Client) RotateS3SecretKey(ctx context.Context, serviceID, secretKey string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/s3/key", Query("secret_key", secretKey))
}

// URL Shortener extras (folder "URL Shortener").

// GetURLShortenerConfig returns the shortener configuration
// (GET /api/url-shortener/config). No secrets are returned.
func (c *Client) GetURLShortenerConfig(ctx context.Context) (map[string]any, error) {
	return c.getMap(ctx, "/api/url-shortener/config", nil)
}

// GetURLShortenerStats returns link statistics
// (GET /api/url-shortener/stats).
func (c *Client) GetURLShortenerStats(ctx context.Context) (map[string]any, error) {
	return c.getMap(ctx, "/api/url-shortener/stats", nil)
}
