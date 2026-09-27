// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// HiTechCloud AI Factory extras (folder "HiTechCloud AI Factory"): service
// provisioning mode, JSON examples, schema and the service-level instance.

// GetAIProvisioningMode returns the provisioning mode of the AI service
// (GET /api/service/:id/hitechcloud/mode).
func (c *Client) GetAIProvisioningMode(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/hitechcloud/mode", nil)
}

// GetAIJSONExamples returns JSON payload examples
// (GET /api/service/:id/hitechcloud/examples).
func (c *Client) GetAIJSONExamples(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/hitechcloud/examples", nil)
}

// GetAISchema returns the API schema of the AI service
// (GET /api/service/:id/hitechcloud/schema).
func (c *Client) GetAISchema(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/hitechcloud/schema", nil)
}

// GetAIServiceInstance returns the service-level AI instance
// (GET /api/service/:id/instance).
func (c *Client) GetAIServiceInstance(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/instance", nil)
}

// SyncAIServiceInstance synchronises the service instance
// (POST /api/service/:id/instance/sync).
func (c *Client) SyncAIServiceInstance(ctx context.Context, serviceID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/instance/sync", nil)
	return err
}

// RestartAIServiceInstance restarts the service instance
// (POST /api/service/:id/instance/restart).
func (c *Client) RestartAIServiceInstance(ctx context.Context, serviceID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/instance/restart", nil)
	return err
}

// UpdateAIServiceInstance updates the service instance
// (POST /api/service/:id/instance/update). params accepts: name, auto_delete,
// alert, tags.
func (c *Client) UpdateAIServiceInstance(ctx context.Context, serviceID string, params map[string]string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/instance/update", qmap(params))
	return err
}

// RestartAIInstance restarts one AI instance
// (POST /api/service/:id/instances/:instance_id/restart).
func (c *Client) RestartAIInstance(ctx context.Context, serviceID, instanceID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/instances/"+PathEscape(instanceID)+"/restart", nil)
	return err
}

// ListAIFeaturedTemplates lists featured templates
// (GET /api/service/:id/templates/featured).
func (c *Client) ListAIFeaturedTemplates(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/templates/featured", nil)
}

// PasskeyV2 / EmailMfaV2 (folders "PasskeyV2", "EmailMfaV2"): multi-factor
// authentication management.

// GetPasskeyMFAStatus returns MFA status
// (GET /api/passkeyv2/status/:user_type/:user_id).
func (c *Client) GetPasskeyMFAStatus(ctx context.Context, userType, userID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/passkeyv2/status/"+PathEscape(userType)+"/"+PathEscape(userID), nil)
}

// ListPasskeyCredentials lists passkey credentials
// (GET /api/passkeyv2/credentials/:user_type/:user_id).
func (c *Client) ListPasskeyCredentials(ctx context.Context, userType, userID string) ([]any, error) {
	return c.getList(ctx, "/api/passkeyv2/credentials/"+PathEscape(userType)+"/"+PathEscape(userID), nil)
}

// DeletePasskeyCredential removes a passkey credential
// (POST /api/passkeyv2/credentials/:user_type/:user_id/delete).
func (c *Client) DeletePasskeyCredential(ctx context.Context, userType, userID, credentialID string) error {
	_, err := c.postMap(ctx, "/api/passkeyv2/credentials/"+PathEscape(userType)+"/"+PathEscape(userID)+"/delete",
		Query("credential_id", credentialID))
	return err
}

// SendPasskeyEmailOTP sends an email OTP (POST /api/passkeyv2/email_otp/send).
func (c *Client) SendPasskeyEmailOTP(ctx context.Context, userType, userID string) error {
	_, err := c.postMap(ctx, "/api/passkeyv2/email_otp/send", Query("user_type", userType, "user_id", userID))
	return err
}

// VerifyPasskeyEmailOTP verifies an email OTP
// (POST /api/passkeyv2/email_otp/verify).
func (c *Client) VerifyPasskeyEmailOTP(ctx context.Context, userType, userID, code string) (map[string]any, error) {
	return c.postMap(ctx, "/api/passkeyv2/email_otp/verify",
		Query("user_type", userType, "user_id", userID, "code", code))
}

// DisablePasskeyMFA disables passkey MFA
// (POST /api/passkeyv2/disable/:user_type/:user_id).
func (c *Client) DisablePasskeyMFA(ctx context.Context, userType, userID string) error {
	_, err := c.postMap(ctx, "/api/passkeyv2/disable/"+PathEscape(userType)+"/"+PathEscape(userID), nil)
	return err
}

// GetEmailMFAStatus returns email MFA status
// (GET /api/email_mfa_v2/status/:user_type/:user_id).
func (c *Client) GetEmailMFAStatus(ctx context.Context, userType, userID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/email_mfa_v2/status/"+PathEscape(userType)+"/"+PathEscape(userID), nil)
}

// SendEmailMFACode sends a verification code (POST /api/email_mfa_v2/send).
func (c *Client) SendEmailMFACode(ctx context.Context, userType, userID, purpose string) error {
	_, err := c.postMap(ctx, "/api/email_mfa_v2/send", Query("user_type", userType, "user_id", userID, "purpose", purpose))
	return err
}

// VerifyEmailMFACode verifies a code (POST /api/email_mfa_v2/verify).
func (c *Client) VerifyEmailMFACode(ctx context.Context, userType, userID, code, purpose string) (map[string]any, error) {
	return c.postMap(ctx, "/api/email_mfa_v2/verify",
		Query("user_type", userType, "user_id", userID, "code", code, "purpose", purpose))
}

// ListEmailMFACodes lists active codes (GET /api/email_mfa_v2/list/:user_type/:user_id).
func (c *Client) ListEmailMFACodes(ctx context.Context, userType, userID string) ([]any, error) {
	return c.getList(ctx, "/api/email_mfa_v2/list/"+PathEscape(userType)+"/"+PathEscape(userID), nil)
}

// RevokeAllEmailMFACodes revokes every active code
// (POST /api/email_mfa_v2/revokeall).
func (c *Client) RevokeAllEmailMFACodes(ctx context.Context, userType, userID string) error {
	_, err := c.postMap(ctx, "/api/email_mfa_v2/revokeall", Query("user_type", userType, "user_id", userID))
	return err
}

// DisableEmailMFA disables email MFA (POST /api/email_mfa_v2/disable).
func (c *Client) DisableEmailMFA(ctx context.Context, userType, userID string) error {
	_, err := c.postMap(ctx, "/api/email_mfa_v2/disable", Query("user_type", userType, "user_id", userID))
	return err
}
