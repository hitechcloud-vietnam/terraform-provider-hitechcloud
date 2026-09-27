// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// VNeIDEKYC (folder "VNeIDEKYC"): electronic KYC sessions for persons and
// organizations.

// StartEKYCSession starts a personal eKYC session
// (GET /api/vneidekyc/ekyc/session).
func (c *Client) StartEKYCSession(ctx context.Context, clientID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/vneidekyc/ekyc/session", Query("client_id", clientID))
}

// GetEKYCSession returns one eKYC session
// (GET /api/vneidekyc/ekyc/session/:session_hash).
func (c *Client) GetEKYCSession(ctx context.Context, sessionHash string) (map[string]any, error) {
	return c.getMap(ctx, "/api/vneidekyc/ekyc/session/"+PathEscape(sessionHash), nil)
}

// ListEKYCSessions lists the eKYC sessions of a client
// (GET /api/vneidekyc/ekyc/list/:client_id).
func (c *Client) ListEKYCSessions(ctx context.Context, clientID, limit string) ([]any, error) {
	return c.getList(ctx, "/api/vneidekyc/ekyc/list/"+PathEscape(clientID), Query("limit", limit))
}

// UploadEKYCFile uploads an eKYC document
// (POST /api/vneidekyc/ekyc/:session_hash/upload).
func (c *Client) UploadEKYCFile(ctx context.Context, sessionHash, fileRole, file string) (map[string]any, error) {
	return c.postMap(ctx, "/api/vneidekyc/ekyc/"+PathEscape(sessionHash)+"/upload",
		Query("file_role", fileRole, "file", file))
}

// SubmitEKYC submits an eKYC session for review
// (POST /api/vneidekyc/ekyc/:session_hash/submit).
func (c *Client) SubmitEKYC(ctx context.Context, sessionHash string, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/vneidekyc/ekyc/"+PathEscape(sessionHash)+"/submit", qmap(params))
}

// CancelEKYCSession cancels an eKYC session
// (POST /api/vneidekyc/ekyc/:session_hash/cancel).
func (c *Client) CancelEKYCSession(ctx context.Context, sessionHash string) error {
	_, err := c.postMap(ctx, "/api/vneidekyc/ekyc/"+PathEscape(sessionHash)+"/cancel", nil)
	return err
}

// AcceptEKYCSession accepts an eKYC session (admin)
// (POST /api/vneidekyc/ekyc/:session_hash/accept).
func (c *Client) AcceptEKYCSession(ctx context.Context, sessionHash, reviewNote string) error {
	_, err := c.postMap(ctx, "/api/vneidekyc/ekyc/"+PathEscape(sessionHash)+"/accept", Query("review_note", reviewNote))
	return err
}

// RejectEKYCSession rejects an eKYC session (admin)
// (POST /api/vneidekyc/ekyc/:session_hash/reject).
func (c *Client) RejectEKYCSession(ctx context.Context, sessionHash, rejectionReason, reviewNote string) error {
	_, err := c.postMap(ctx, "/api/vneidekyc/ekyc/"+PathEscape(sessionHash)+"/reject",
		Query("rejection_reason", rejectionReason, "review_note", reviewNote))
	return err
}

// LookupOrganizationTaxCode looks up an organization by tax code
// (POST /api/vneidekyc/org/lookup).
func (c *Client) LookupOrganizationTaxCode(ctx context.Context, taxCode string) (map[string]any, error) {
	return c.postMap(ctx, "/api/vneidekyc/org/lookup", Query("tax_code", taxCode))
}

// StartOrganizationVerification starts organization verification
// (POST /api/vneidekyc/org/start).
func (c *Client) StartOrganizationVerification(ctx context.Context, clientID, taxCode string) (map[string]any, error) {
	return c.postMap(ctx, "/api/vneidekyc/org/start", Query("client_id", clientID, "tax_code", taxCode))
}

// GetOrganization returns an organization verification record
// (GET /api/vneidekyc/org/:org_hash).
func (c *Client) GetOrganization(ctx context.Context, orgHash string) (map[string]any, error) {
	return c.getMap(ctx, "/api/vneidekyc/org/"+PathEscape(orgHash), nil)
}

// AcceptOrganization accepts an organization verification (admin)
// (POST /api/vneidekyc/org/:org_hash/accept).
func (c *Client) AcceptOrganization(ctx context.Context, orgHash string) error {
	_, err := c.postMap(ctx, "/api/vneidekyc/org/"+PathEscape(orgHash)+"/accept", nil)
	return err
}

// RejectOrganization rejects an organization verification (admin)
// (POST /api/vneidekyc/org/:org_hash/reject).
func (c *Client) RejectOrganization(ctx context.Context, orgHash, rejectionReason string) error {
	_, err := c.postMap(ctx, "/api/vneidekyc/org/"+PathEscape(orgHash)+"/reject", Query("rejection_reason", rejectionReason))
	return err
}

// UploadSignedPDF uploads a signed PDF (POST /api/vneidekyc/pdf/upload).
func (c *Client) UploadSignedPDF(ctx context.Context, clientID, docType, orgHash, file string) (map[string]any, error) {
	return c.postMap(ctx, "/api/vneidekyc/pdf/upload",
		Query("client_id", clientID, "doc_type", docType, "org_hash", orgHash, "file", file))
}

// GetSignedPDF returns a signed PDF record (GET /api/vneidekyc/pdf/:pdf_hash).
func (c *Client) GetSignedPDF(ctx context.Context, pdfHash string) (map[string]any, error) {
	return c.getMap(ctx, "/api/vneidekyc/pdf/"+PathEscape(pdfHash), nil)
}

// GetClientVerificationStatus returns the verification status of a client
// (GET /api/vneidekyc/client/:client_id/status).
func (c *Client) GetClientVerificationStatus(ctx context.Context, clientID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/vneidekyc/client/"+PathEscape(clientID)+"/status", nil)
}

// WillExpired (folder "WillExpired"): expiring services and renewals.

// ListWillExpired lists expiring items (GET /api/willexpired). params accepts:
// days, type, status, search, sort, dir, page, per_page.
func (c *Client) ListWillExpired(ctx context.Context, params map[string]string) ([]any, error) {
	return c.getList(ctx, "/api/willexpired", qmap(params))
}

// GetWillExpiredSummary returns the expiring summary
// (GET /api/willexpired/summary).
func (c *Client) GetWillExpiredSummary(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.getMap(ctx, "/api/willexpired/summary", qmap(params))
}

// GetWillExpiredConfig returns module configuration
// (GET /api/willexpired/config).
func (c *Client) GetWillExpiredConfig(ctx context.Context) (map[string]any, error) {
	return c.getMap(ctx, "/api/willexpired/config", nil)
}

// ListWillExpiredInvoices lists open renewal invoices
// (GET /api/willexpired/invoices).
func (c *Client) ListWillExpiredInvoices(ctx context.Context, params map[string]string) ([]any, error) {
	return c.getList(ctx, "/api/willexpired/invoices", qmap(params))
}

// ListWillExpiredRequests lists the renewal request log
// (GET /api/willexpired/requests).
func (c *Client) ListWillExpiredRequests(ctx context.Context, limit string) ([]any, error) {
	return c.getList(ctx, "/api/willexpired/requests", Query("limit", limit))
}

// ExportWillExpired exports the expiring list (GET /api/willexpired/export).
func (c *Client) ExportWillExpired(ctx context.Context, params map[string]string) ([]byte, error) {
	return c.GetRaw(ctx, "/api/willexpired/export", qmap(params))
}

// GetWillExpiredItem returns one item (GET /api/willexpired/:type/:id).
func (c *Client) GetWillExpiredItem(ctx context.Context, itemType, id string) (map[string]any, error) {
	return c.getMap(ctx, "/api/willexpired/"+PathEscape(itemType)+"/"+PathEscape(id), nil)
}

// RenewWillExpiredItem requests renewal (POST /api/willexpired/:type/:id/renew).
func (c *Client) RenewWillExpiredItem(ctx context.Context, itemType, id string) (map[string]any, error) {
	return c.postMap(ctx, "/api/willexpired/"+PathEscape(itemType)+"/"+PathEscape(id)+"/renew", nil)
}

// GetWillExpiredAutoRenew returns the auto-renew flag
// (GET /api/willexpired/:type/:id/autorenew).
func (c *Client) GetWillExpiredAutoRenew(ctx context.Context, itemType, id string) (map[string]any, error) {
	return c.getMap(ctx, "/api/willexpired/"+PathEscape(itemType)+"/"+PathEscape(id)+"/autorenew", nil)
}

// SetWillExpiredAutoRenew sets the auto-renew flag
// (PUT /api/willexpired/:type/:id/autorenew).
func (c *Client) SetWillExpiredAutoRenew(ctx context.Context, itemType, id, autorenew string) error {
	_, err := c.putMap(ctx, "/api/willexpired/"+PathEscape(itemType)+"/"+PathEscape(id)+"/autorenew",
		Query("autorenew", autorenew))
	return err
}
