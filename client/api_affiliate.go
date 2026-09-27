// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// Affiliate program (folder "Affiliate") and the advanced affiliate manager
// (folder "AffiliatesAdvanced").

// GetAffiliateSummary returns the affiliate summary (GET /api/affiliates/summary).
func (c *Client) GetAffiliateSummary(ctx context.Context) (map[string]any, error) {
	return c.getMap(ctx, "/api/affiliates/summary", nil)
}

// ListAffiliateCampaigns lists affiliate campaigns (GET /api/affiliates/campaigns).
func (c *Client) ListAffiliateCampaigns(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/affiliates/campaigns", nil)
}

// ListAffiliateCommissions lists affiliate commissions (GET /api/affiliates/commissions).
func (c *Client) ListAffiliateCommissions(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/affiliates/commissions", nil)
}

// ListAffiliatePayouts lists affiliate payouts (GET /api/affiliates/payouts).
func (c *Client) ListAffiliatePayouts(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/affiliates/payouts", nil)
}

// ListAffiliateVouchers lists affiliate vouchers (GET /api/affiliates/vouchers).
func (c *Client) ListAffiliateVouchers(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/affiliates/vouchers", nil)
}

// ListAffiliateCommissionPlans lists commission plans
// (GET /api/affiliates/commissionplans).
func (c *Client) ListAffiliateCommissionPlans(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/affiliates/commissionplans", nil)
}

// AffiliatesAdvanced ------------------------------------------------------------

// GetAffiliateAdvInfo returns the affiliate profile of clientID
// (GET /api/affiliates_adv/:client_id/info).
func (c *Client) GetAffiliateAdvInfo(ctx context.Context, clientID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/info", nil)
}

// GetAffiliateAdvStats returns affiliate statistics
// (GET /api/affiliates_adv/:client_id/stats).
func (c *Client) GetAffiliateAdvStats(ctx context.Context, clientID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/stats", nil)
}

// ListAffiliateAdvCommissionPlans lists commission plans
// (GET /api/affiliates_adv/:client_id/commission-plans). params accepts
// voucher_enabled.
func (c *Client) ListAffiliateAdvCommissionPlans(ctx context.Context, clientID string, params map[string]string) ([]any, error) {
	return c.getList(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/commission-plans", qmap(params))
}

// ListAffiliateAdvVouchers lists vouchers
// (GET /api/affiliates_adv/:client_id/vouchers).
func (c *Client) ListAffiliateAdvVouchers(ctx context.Context, clientID string) ([]any, error) {
	return c.getList(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/vouchers", nil)
}

// ListAffiliateAdvCommissions lists commissions
// (GET /api/affiliates_adv/:client_id/commissions). params accepts order_id,
// paid_status, date_from, date_to, page, perpage, orderby.
func (c *Client) ListAffiliateAdvCommissions(ctx context.Context, clientID string, params map[string]string) ([]any, error) {
	return c.getList(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/commissions", qmap(params))
}

// ListAffiliateAdvReferrals lists referred clients
// (GET /api/affiliates_adv/:client_id/referrals).
func (c *Client) ListAffiliateAdvReferrals(ctx context.Context, clientID string) ([]any, error) {
	return c.getList(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/referrals", nil)
}

// ListAffiliateAdvPayouts lists payout history
// (GET /api/affiliates_adv/:client_id/payouts).
func (c *Client) ListAffiliateAdvPayouts(ctx context.Context, clientID string) ([]any, error) {
	return c.getList(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/payouts", nil)
}

// ListAffiliateAdvCampaigns lists campaigns
// (GET /api/affiliates_adv/:client_id/campaigns).
func (c *Client) ListAffiliateAdvCampaigns(ctx context.Context, clientID string) ([]any, error) {
	return c.getList(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/campaigns", nil)
}

// ListAffiliateAdvAudit returns the affiliate activity log
// (GET /api/affiliates_adv/:client_id/audit).
func (c *Client) ListAffiliateAdvAudit(ctx context.Context, clientID string) ([]any, error) {
	return c.getList(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/audit", nil)
}

// ActivateAffiliateAdv activates the affiliate account
// (POST /api/affiliates_adv/:client_id/activate).
func (c *Client) ActivateAffiliateAdv(ctx context.Context, clientID string) error {
	_, err := c.postMap(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/activate", nil)
	return err
}

// SetAffiliateAdvCommissionPlan assigns a commission plan
// (POST /api/affiliates_adv/:client_id/commission-plan/:commission_id).
func (c *Client) SetAffiliateAdvCommissionPlan(ctx context.Context, clientID, commissionID string) error {
	_, err := c.postMap(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/commission-plan/"+PathEscape(commissionID), nil)
	return err
}

// CreateAffiliateAdvVoucher creates a voucher
// (POST /api/affiliates_adv/:client_id/vouchers/:plan_id). params accepts
// code, discount, cycle, expires, max_usage, audience.
func (c *Client) CreateAffiliateAdvVoucher(ctx context.Context, clientID, planID string, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/vouchers/"+PathEscape(planID), qmap(params))
}

// DeleteAffiliateAdvVoucher removes a voucher
// (DELETE /api/affiliates_adv/:client_id/vouchers/:voucher_id).
func (c *Client) DeleteAffiliateAdvVoucher(ctx context.Context, clientID, voucherID string) error {
	_, err := c.deleteMap(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/vouchers/"+PathEscape(voucherID), nil)
	return err
}

// SetAffiliateAdvLandingPage sets the affiliate landing page
// (POST /api/affiliates_adv/:client_id/landing-page).
func (c *Client) SetAffiliateAdvLandingPage(ctx context.Context, clientID, url string) error {
	_, err := c.postMap(ctx, "/api/affiliates_adv/"+PathEscape(clientID)+"/landing-page", Query("url", url))
	return err
}
