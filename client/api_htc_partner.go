// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// HiTechCloudPMG (folder "HiTechCloudPMG"): Proxmox Mail Gateway.

// GetPMGConfig returns the mail filtering configuration
// (GET /api/service/:id/htcpmg/config).
func (c *Client) GetPMGConfig(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpmg/config", nil)
}

// AddPMGDomain adds a mail domain (POST /api/service/:id/htcpmg/domains).
func (c *Client) AddPMGDomain(ctx context.Context, serviceID, domain string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpmg/domains", Query("domain", domain))
	return err
}

// SetPMGTransport sets the target mail server of a domain
// (POST /api/service/:id/htcpmg/transport).
func (c *Client) SetPMGTransport(ctx context.Context, serviceID, domain, host, port string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpmg/transport",
		Query("domain", domain, "host", host, "port", port))
	return err
}

// GetPMGStats returns mail statistics (GET /api/service/:id/htcpmg/stats).
func (c *Client) GetPMGStats(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpmg/stats", nil)
}

// HiTechCloudProxmox (folder "HiTechCloudProxmox"): Proxmox VE.

// GetPVEStatus returns machine status (GET /api/service/:id/htcpve/status).
func (c *Client) GetPVEStatus(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpve/status", nil)
}

// PVEPowerAction performs a power action on a VM
// (POST /api/service/:id/htcpve/power).
func (c *Client) PVEPowerAction(ctx context.Context, serviceID, action, vmid string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpve/power", Query("action", action, "vmid", vmid))
	return err
}

// ListPVEVMs lists machines (GET /api/service/:id/htcpve/vms).
func (c *Client) ListPVEVMs(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpve/vms", nil)
}

// ListPVEIPs lists addresses (GET /api/service/:id/htcpve/ips).
func (c *Client) ListPVEIPs(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpve/ips", nil)
}

// SetPVERDNS sets reverse DNS of a PVE IP (POST /api/service/:id/htcpve/rdns).
func (c *Client) SetPVERDNS(ctx context.Context, serviceID, ip, hostname string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpve/rdns", Query("ip", ip, "hostname", hostname))
	return err
}

// ListPVEBackups lists backups (GET /api/service/:id/htcpve/backups).
func (c *Client) ListPVEBackups(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpve/backups", nil)
}

// CreatePVEBackup creates a backup (POST /api/service/:id/htcpve/backups).
func (c *Client) CreatePVEBackup(ctx context.Context, serviceID, mode, notes string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpve/backups", Query("mode", mode, "notes", notes))
}

// ListPVESnapshots lists snapshots (GET /api/service/:id/htcpve/snapshots).
func (c *Client) ListPVESnapshots(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpve/snapshots", nil)
}

// CreatePVESnapshot creates a snapshot (POST /api/service/:id/htcpve/snapshots).
func (c *Client) CreatePVESnapshot(ctx context.Context, serviceID, name, description string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpve/snapshots",
		Query("name", name, "description", description))
}

// GetPVEUsage returns bandwidth usage (GET /api/service/:id/htcpve/usage).
func (c *Client) GetPVEUsage(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcpve/usage", nil)
}

// HiTechCloudIPAM (folder "HiTechCloudIPAM").

// ListIPAMIPs lists service IP addresses (GET /api/service/:id/htcipam/ips).
func (c *Client) ListIPAMIPs(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/htcipam/ips", nil)
}

// ListIPAMSubnets lists service subnets (GET /api/service/:id/htcipam/subnets).
func (c *Client) ListIPAMSubnets(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/htcipam/subnets", nil)
}

// GetIPAMRDNS reads reverse DNS (GET /api/service/:id/htcipam/rdns).
func (c *Client) GetIPAMRDNS(ctx context.Context, serviceID string) ([]RDNSRecord, error) {
	var raw any
	if err := c.Get(ctx, "/api/service/"+PathEscape(serviceID)+"/htcipam/rdns", nil, &raw); err != nil {
		return nil, err
	}
	return parseRDNS(raw), nil
}

// SetIPAMRDNS sets reverse DNS of an IPAM address
// (POST /api/service/:id/htcipam/rdns).
func (c *Client) SetIPAMRDNS(ctx context.Context, serviceID, ip, hostname string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/htcipam/rdns", Query("ip", ip, "hostname", hostname))
	return err
}

// Partner (folder "Partner"): partner program.

// GetPartnerProfile returns the partner profile (GET /api/partner).
func (c *Client) GetPartnerProfile(ctx context.Context) (map[string]any, error) {
	return c.getMap(ctx, "/api/partner", nil)
}

// ApplyPartner applies to the partner program (POST /api/partner/apply).
// params accepts: type, company_name, tax_code, contact_name, contact_email,
// contact_phone, website, intro.
func (c *Client) ApplyPartner(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/partner/apply", qmap(params))
}

// ListPartnerTiers lists partner tiers (GET /api/partner/tiers).
func (c *Client) ListPartnerTiers(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/partner/tiers", nil)
}

// GetPartnerPricing returns partner pricing for a category
// (GET /api/partner/pricing).
func (c *Client) GetPartnerPricing(ctx context.Context, category string) (map[string]any, error) {
	return c.getMap(ctx, "/api/partner/pricing", Query("category", category))
}

// ListPartnerCustomers lists partner customers (GET /api/partner/customers).
func (c *Client) ListPartnerCustomers(ctx context.Context, params map[string]string) ([]any, error) {
	return c.getList(ctx, "/api/partner/customers", qmap(params))
}

// ListPartnerEarnings lists partner earnings (GET /api/partner/earnings).
func (c *Client) ListPartnerEarnings(ctx context.Context, params map[string]string) ([]any, error) {
	return c.getList(ctx, "/api/partner/earnings", qmap(params))
}

// GetPartnerWallet returns the partner wallet (GET /api/partner/wallet).
func (c *Client) GetPartnerWallet(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.getMap(ctx, "/api/partner/wallet", qmap(params))
}

// ListPartnerPayouts lists payouts (GET /api/partner/payouts).
func (c *Client) ListPartnerPayouts(ctx context.Context, params map[string]string) ([]any, error) {
	return c.getList(ctx, "/api/partner/payouts", qmap(params))
}

// RequestPartnerPayout requests a payout (POST /api/partner/payouts).
func (c *Client) RequestPartnerPayout(ctx context.Context, amount, method, note string) (map[string]any, error) {
	return c.postMap(ctx, "/api/partner/payouts", Query("amount", amount, "method", method, "note", note))
}

// ListPartnerLeads lists leads (GET /api/partner/leads).
func (c *Client) ListPartnerLeads(ctx context.Context, params map[string]string) ([]any, error) {
	return c.getList(ctx, "/api/partner/leads", qmap(params))
}

// RegisterPartnerLead registers a lead (POST /api/partner/leads).
func (c *Client) RegisterPartnerLead(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/partner/leads", qmap(params))
}

// GetPartnerReferralLink returns the referral link (GET /api/partner/referral).
func (c *Client) GetPartnerReferralLink(ctx context.Context, to string) (map[string]any, error) {
	return c.getMap(ctx, "/api/partner/referral", Query("to", to))
}

// GetPartnerRateCard returns the partner rate card (GET /api/partner/rates).
func (c *Client) GetPartnerRateCard(ctx context.Context) (map[string]any, error) {
	return c.getMap(ctx, "/api/partner/rates", nil)
}
