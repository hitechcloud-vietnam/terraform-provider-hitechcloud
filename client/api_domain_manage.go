// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// Domain management extras (folders "Domains", "DNS Manage"): EPP codes,
// sync, contacts, forwarding, ordering, documents and DNSSEC.

// GetDomainEPPCode returns the EPP/auth code of a domain
// (GET /api/domain/:id/epp).
func (c *Client) GetDomainEPPCode(ctx context.Context, domainID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/domain/"+PathEscape(domainID)+"/epp", nil)
}

// SyncDomain synchronises the domain with the registrar
// (GET /api/domain/:id/sync).
func (c *Client) SyncDomain(ctx context.Context, domainID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/domain/"+PathEscape(domainID)+"/sync", nil)
}

// GetDomainContactInfo returns the registrant contact of a domain
// (GET /api/domain/:id/contact).
func (c *Client) GetDomainContactInfo(ctx context.Context, domainID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/domain/"+PathEscape(domainID)+"/contact", nil)
}

// UpdateDomainContactInfo updates the registrant contact
// (PUT /api/domain/:id/contact). contactInfo is the contact payload
// understood by the API.
func (c *Client) UpdateDomainContactInfo(ctx context.Context, domainID, contactInfo string) error {
	_, err := c.putMap(ctx, "/api/domain/"+PathEscape(domainID)+"/contact", Query("contact_info", contactInfo))
	return err
}

// GetEmailForwarding returns the email forwarding rules
// (GET /api/domain/:id/emforwarding).
func (c *Client) GetEmailForwarding(ctx context.Context, domainID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/domain/"+PathEscape(domainID)+"/emforwarding", nil)
}

// UpdateEmailForwarding sets an email forwarding rule
// (PUT /api/domain/:id/emforwarding).
func (c *Client) UpdateEmailForwarding(ctx context.Context, domainID, from, to string) error {
	_, err := c.putMap(ctx, "/api/domain/"+PathEscape(domainID)+"/emforwarding", Query("from", from, "to", to))
	return err
}

// UpdateDomainForwarding sets web forwarding for the domain
// (PUT /api/domain/:id/forwarding). params accepts the forwarding target and
// type fields understood by the API.
func (c *Client) UpdateDomainForwarding(ctx context.Context, domainID string, params map[string]string) error {
	_, err := c.putMap(ctx, "/api/domain/"+PathEscape(domainID)+"/forwarding", qmap(params))
	return err
}

// GetDomainDocuments lists domain documents (GET /api/domain/:id/documents).
func (c *Client) GetDomainDocuments(ctx context.Context, domainID string) ([]any, error) {
	return c.getList(ctx, "/api/domain/"+PathEscape(domainID)+"/documents", nil)
}

// RenewDomain renews a domain (POST /api/domain/:id/renew).
func (c *Client) RenewDomain(ctx context.Context, domainID, years, payMethod string) (map[string]any, error) {
	return c.postMap(ctx, "/api/domain/"+PathEscape(domainID)+"/renew", Query("years", years, "pay_method", payMethod))
}

// OrderDomain registers a new domain (POST /api/domain/order). params
// accepts: name, years, action, tld_id, pay_method, epp, nameservers,
// registrant, admin, tech, billing, data, aff_id.
func (c *Client) OrderDomain(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/domain/order", qmap(params))
}

// GetTLDOrderForm returns the additional registration data form for a TLD
// (GET /api/domain/order/:id/form).
func (c *Client) GetTLDOrderForm(ctx context.Context, productID, tldID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/domain/order/"+PathEscape(productID)+"/form", Query("tld_id", tldID))
}

// WhoisLookup performs a WHOIS lookup with captcha token
// (GET /api/whoislookup/:domain).
func (c *Client) WhoisLookup(ctx context.Context, domain, captcha string) (map[string]string, error) {
	var raw any
	if err := c.Get(ctx, "/api/whoislookup/"+PathEscape(domain), Query("captcha", captcha), &raw); err != nil {
		return nil, err
	}
	return StringMap(AsMapAny(raw)), nil
}

// RegisterDomainNameservers registers custom nameservers at the registrar
// (POST /api/domain/:id/reg). params accepts the ns1..nsN host/ip fields.
func (c *Client) RegisterDomainNameservers(ctx context.Context, domainID string, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/domain/"+PathEscape(domainID)+"/reg", qmap(params))
}

// ListDNSSECFlags returns the DNSSEC flags supported for a domain
// (GET /api/domain/:id/dnssec/flags).
func (c *Client) ListDNSSECFlags(ctx context.Context, domainID string) ([]any, error) {
	return c.getList(ctx, "/api/domain/"+PathEscape(domainID)+"/dnssec/flags", nil)
}

// ListDNSSECKeys lists the DNSSEC keys of a domain
// (GET /api/domain/:id/dnssec).
func (c *Client) ListDNSSECKeys(ctx context.Context, domainID string) ([]any, error) {
	return c.getList(ctx, "/api/domain/"+PathEscape(domainID)+"/dnssec", nil)
}

// AddDNSSECKey adds a DNSSEC key (PUT /api/domain/:id/dnssec).
func (c *Client) AddDNSSECKey(ctx context.Context, domainID string, params map[string]string) error {
	_, err := c.putMap(ctx, "/api/domain/"+PathEscape(domainID)+"/dnssec", qmap(params))
	return err
}

// RemoveDNSSECKey deletes a DNSSEC key (DELETE /api/domain/:id/dnssec/:key).
func (c *Client) RemoveDNSSECKey(ctx context.Context, domainID, key string) error {
	_, err := c.deleteMap(ctx, "/api/domain/"+PathEscape(domainID)+"/dnssec/"+PathEscape(key), nil)
	return err
}
