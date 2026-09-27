// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// Bare Metal & Colocation / Collocation Services / Hosting Services
// folders: reinstall, diagnostics, rescue, power, IPs, VLANs, PDU,
// bandwidth and server stock.

// ListReinstallTemplates lists OS templates and recipes for reinstall
// (GET /api/service/:id/reinstall/templates).
func (c *Client) ListReinstallTemplates(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/reinstall/templates", nil)
}

// GetReinstallDetails returns reinstall details (GET /api/service/:id/reinstall).
func (c *Client) GetReinstallDetails(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/reinstall", nil)
}

// ReinstallServer reinstalls the server (POST /api/service/:id/reinstall).
// params accepts: profile, rootpassword, adminuser, userpassword,
// packageselection, extra1.
func (c *Client) ReinstallServer(ctx context.Context, serviceID string, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/reinstall", qmap(params))
}

// CancelDiagnostics cancels a running diagnostics job
// (POST /api/service/:id/diag/cancel).
func (c *Client) CancelDiagnostics(ctx context.Context, serviceID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/diag/cancel", nil)
	return err
}

// ListDiagnosticsTemplates lists diagnostics templates
// (GET /api/service/:id/diag/templates).
func (c *Client) ListDiagnosticsTemplates(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/diag/templates", nil)
}

// GetDiagnosticsStatus returns diagnostics status (GET /api/service/:id/diag).
func (c *Client) GetDiagnosticsStatus(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/diag", nil)
}

// RunDiagnostics starts diagnostics (POST /api/service/:id/diag).
func (c *Client) RunDiagnostics(ctx context.Context, serviceID string, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/diag", qmap(params))
}

// ListRescueTemplates lists rescue templates (GET /api/service/:id/rescue/templates).
func (c *Client) ListRescueTemplates(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/rescue/templates", nil)
}

// GetRescueStatus returns rescue status (GET /api/service/:id/rescue).
func (c *Client) GetRescueStatus(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/rescue", nil)
}

// RescueServer boots the server into rescue mode (POST /api/service/:id/rescue).
func (c *Client) RescueServer(ctx context.Context, serviceID, template, password string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/rescue", Query("template", template, "password", password))
}

// CancelRescue cancels rescue mode (POST /api/service/:id/rescue/cancel).
func (c *Client) CancelRescue(ctx context.Context, serviceID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/rescue/cancel", nil)
	return err
}

// GetServerInfo returns hardware/server info (GET /api/service/:id/info).
func (c *Client) GetServerInfo(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/info", nil)
}

// GetServerPowerStatus returns the power status (GET /api/service/:id/status).
func (c *Client) GetServerPowerStatus(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/status", nil)
}

// UpdateServerHostname renames the server (POST /api/service/:id/hostname).
func (c *Client) UpdateServerHostname(ctx context.Context, serviceID, hostname string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/hostname", Query("hostname", hostname))
	return err
}

// ListServerIPs lists IPs of a bare-metal service (GET /api/service/:id/ips).
func (c *Client) ListServerIPs(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/ips", nil)
}

// AddServerIP adds IPs to a bare-metal service (POST /api/service/:id/ips).
func (c *Client) AddServerIP(ctx context.Context, serviceID, vlan, domain, num string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/ips", Query("vlan", vlan, "domain", domain, "num", num))
}

// GetServerIPDetails returns one IP (GET /api/service/:id/ips/:ip).
func (c *Client) GetServerIPDetails(ctx context.Context, serviceID, ip string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/ips/"+PathEscape(ip), nil)
}

// EditServerIP updates an IP (POST /api/service/:id/ips/:ip).
func (c *Client) EditServerIP(ctx context.Context, serviceID, ip, domain string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/ips/"+PathEscape(ip), Query("domain", domain))
	return err
}

// DeleteServerIP removes an IP (DELETE /api/service/:id/ips/:ip).
func (c *Client) DeleteServerIP(ctx context.Context, serviceID, ip string) error {
	_, err := c.deleteMap(ctx, "/api/service/"+PathEscape(serviceID)+"/ips/"+PathEscape(ip), nil)
	return err
}

// ListServerVLANs lists VLANs (GET /api/service/:id/vlans).
func (c *Client) ListServerVLANs(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/vlans", nil)
}

// RebootServer reboots the server (POST /api/service/:id/reboot).
func (c *Client) RebootServer(ctx context.Context, serviceID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/reboot", nil)
	return err
}

// PowerOffServer powers the server off (POST /api/service/:id/poweroff).
func (c *Client) PowerOffServer(ctx context.Context, serviceID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/poweroff", nil)
	return err
}

// PowerOnServer powers the server on (POST /api/service/:id/poweron).
func (c *Client) PowerOnServer(ctx context.Context, serviceID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/poweron", nil)
	return err
}

// ResetServer resets the server (POST /api/service/:id/reset).
func (c *Client) ResetServer(ctx context.Context, serviceID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/reset", nil)
	return err
}

// ListServerStock lists servers in stock (GET /api/serverstock).
func (c *Client) ListServerStock(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/serverstock", nil)
}

// GetBandwidthGraphs returns bandwidth graphs (GET /api/service/:id/bandwidth-graphs).
func (c *Client) GetBandwidthGraphs(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/bandwidth-graphs", nil)
}

// GetBandwidthUsage returns bandwidth usage (GET /api/service/:id/bandwidth).
func (c *Client) GetBandwidthUsage(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/bandwidth", nil)
}

// ListPDUPorts lists PDU ports (GET /api/service/:id/pdu).
func (c *Client) ListPDUPorts(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/pdu", nil)
}

// GetPDUPortState returns one PDU port state (GET /api/service/:id/pdu/:port).
func (c *Client) GetPDUPortState(ctx context.Context, serviceID, port string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/pdu/"+PathEscape(port), nil)
}

// SetPDUPortPower sets a PDU port power state
// (POST /api/service/:id/pdu/:port).
func (c *Client) SetPDUPortPower(ctx context.Context, serviceID, port, power string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/pdu/"+PathEscape(port), Query("power", power))
	return err
}

// LocationV2 (folder "LocationV2").

// ListCountries lists countries (GET /api/location_v2/countries).
func (c *Client) ListCountries(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/location_v2/countries", nil)
}

// ListStates lists states of a country (GET /api/location_v2/:code/state).
func (c *Client) ListStates(ctx context.Context, code string) ([]any, error) {
	return c.getList(ctx, "/api/location_v2/"+PathEscape(code)+"/state", nil)
}

// ListCities lists cities of a state (GET /api/location_v2/state/:id/city).
func (c *Client) ListCities(ctx context.Context, stateID string) ([]any, error) {
	return c.getList(ctx, "/api/location_v2/state/"+PathEscape(stateID)+"/city", nil)
}

// ListCitiesByStateName lists cities by state name
// (GET /api/location_v2/state/city/:name).
func (c *Client) ListCitiesByStateName(ctx context.Context, name string) ([]any, error) {
	return c.getList(ctx, "/api/location_v2/state/city/"+PathEscape(name), nil)
}
