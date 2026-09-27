// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// SSL certificates (folder "SSL Certificates") extras.

// DownloadCertificate downloads the certificate body (PEM)
// (GET /api/certificate/:id/crt).
func (c *Client) DownloadCertificate(ctx context.Context, id string) ([]byte, error) {
	return c.GetRaw(ctx, "/api/certificate/"+PathEscape(id)+"/crt", nil)
}

// ListCertificateProducts lists certificates available for purchase
// (GET /api/certificate/order).
func (c *Client) ListCertificateProducts(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/certificate/order", nil)
}

// OrderCertificate orders a new certificate (POST /api/certificate/order).
// params accepts: product_id, csr, years, pay_method, approver_email, admin,
// tech, billing, organization, software, data, aff_id.
func (c *Client) OrderCertificate(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/certificate/order", qmap(params))
}

// ListCertificateSoftware lists server software options for a certificate
// product (GET /api/certificate/order/:product_id/software).
func (c *Client) ListCertificateSoftware(ctx context.Context, productID string) ([]any, error) {
	return c.getList(ctx, "/api/certificate/order/"+PathEscape(productID)+"/software", nil)
}

// Service administration (folder "Services").

// ListServiceMethods lists the management methods available for a service
// (GET /api/service/:id/methods).
func (c *Client) ListServiceMethods(ctx context.Context, id string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(id)+"/methods", nil)
}

// CancelService requests cancellation of a service
// (POST /api/service/:id/cancel).
func (c *Client) CancelService(ctx context.Context, id, immediate, reason string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(id)+"/cancel", Query("immediate", immediate, "reason", reason))
	return err
}

// GetServiceLabel returns the service label (GET /api/service/:id/label).
func (c *Client) GetServiceLabel(ctx context.Context, id string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(id)+"/label", nil)
}

// ChangeServiceLabel renames the service (POST /api/service/:id/label).
func (c *Client) ChangeServiceLabel(ctx context.Context, id, label string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(id)+"/label", Query("label", label))
	return err
}

// RenewService manually renews a service (POST /api/service/:id/renew).
func (c *Client) RenewService(ctx context.Context, id string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(id)+"/renew", nil)
}

// GetServiceBillingCycle returns the billing cycle (GET /api/service/:id/cycle).
func (c *Client) GetServiceBillingCycle(ctx context.Context, id string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(id)+"/cycle", nil)
}

// ChangeServiceBillingCycle changes the billing cycle (POST /api/service/:id/cycle).
func (c *Client) ChangeServiceBillingCycle(ctx context.Context, id, cycle string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(id)+"/cycle", Query("cycle", cycle))
	return err
}

// Cart (folder "Cart").

// GetProductConfig returns the configuration options of a product
// (GET /api/order/:product_id).
func (c *Client) GetProductConfig(ctx context.Context, productID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/order/"+PathEscape(productID), nil)
}

// OrderService orders a new service (POST /api/order/:product_id). params
// accepts: domain, cycle, pay_method, custom, promocode, aff_id.
func (c *Client) OrderService(ctx context.Context, productID string, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/order/"+PathEscape(productID), qmap(params))
}

// OrderMultipleServices orders several services at once (POST /api/order).
// params accepts: pay_method, ignore_errors, items, aff_id.
func (c *Client) OrderMultipleServices(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/order", qmap(params))
}

// GetOrderQuote returns a quote for a set of items (POST /api/quote).
// params accepts: pay_method, output, items.
func (c *Client) GetOrderQuote(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/quote", qmap(params))
}
