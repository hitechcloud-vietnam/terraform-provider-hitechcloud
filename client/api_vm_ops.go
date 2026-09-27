// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// VM lifecycle and management across the Virtualizor / Cloud GPU / Cloud
// Service / Cloud Instance / Cloud Virtual Machine / vCloudStack folders.
// Power operations on services without a vmid target the service's primary
// VM (the API exposes both shapes).

// SuspendVM suspends a VM (POST /api/service/:id/vms/:vmid/suspend).
func (c *Client) SuspendVM(ctx context.Context, serviceID, vmID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/suspend", nil)
	return err
}

// UnsuspendVM unsuspends a VM (POST /api/service/:id/vms/:vmid/unsuspend).
func (c *Client) UnsuspendVM(ctx context.Context, serviceID, vmID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/unsuspend", nil)
	return err
}

// ShutdownVM gracefully shuts a VM down (POST /api/service/:id/vms/:vmid/shutdown).
func (c *Client) ShutdownVM(ctx context.Context, serviceID, vmID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/shutdown", nil)
	return err
}

// ResetVM hard-resets a VM (POST /api/service/:id/vms/:vmid/reset).
func (c *Client) ResetVM(ctx context.Context, serviceID, vmID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/reset", nil)
	return err
}

// StopVM stops a VM (POST /api/service/:id/vms/:vmid/stop).
func (c *Client) StopVM(ctx context.Context, serviceID, vmID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/stop", nil)
	return err
}

// StartVM starts a VM (POST /api/service/:id/vms/:vmid/start).
func (c *Client) StartVM(ctx context.Context, serviceID, vmID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/start", nil)
	return err
}

// RebootVM reboots a VM (POST /api/service/:id/vms/:vmid/reboot). iso is
// optional and mounts an ISO on reboot.
func (c *Client) RebootVM(ctx context.Context, serviceID, vmID, iso string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/reboot", Query("iso", iso))
	return err
}

// RebootVMService reboots the service's VM without a vmid
// (PUT /api/service/:id/vms/reboot).
func (c *Client) RebootVMService(ctx context.Context, serviceID string) error {
	_, err := c.putMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/reboot", nil)
	return err
}

// StopVMService stops the service's VM without a vmid
// (PUT /api/service/:id/vms/stop).
func (c *Client) StopVMService(ctx context.Context, serviceID string) error {
	_, err := c.putMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/stop", nil)
	return err
}

// StartVMService starts the service's VM without a vmid
// (PUT /api/service/:id/vms/start).
func (c *Client) StartVMService(ctx context.Context, serviceID string) error {
	_, err := c.putMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/start", nil)
	return err
}

// ResetVMPassword resets the VM root password
// (POST /api/service/:id/vms/:vmid/resetpwd).
func (c *Client) ResetVMPassword(ctx context.Context, serviceID, vmID string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/resetpwd", nil)
}

// ListVMRebuildTemplates lists rebuild templates for a VM
// (GET /api/service/:id/vms/:vmid/rebuild).
func (c *Client) ListVMRebuildTemplates(ctx context.Context, serviceID, vmID, template string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/rebuild", Query("template", template))
}

// RebuildVM rebuilds a VM from a template
// (POST /api/service/:id/vms/:vmid/rebuild).
func (c *Client) RebuildVM(ctx context.Context, serviceID, vmID string, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/rebuild", qmap(params))
}

// ChangeVMSSHKey assigns an SSH key to the VM
// (POST /api/service/:id/vms/:vmid/addsshkey).
func (c *Client) ChangeVMSSHKey(ctx context.Context, serviceID, vmID string, params map[string]string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/addsshkey", qmap(params))
	return err
}

// ChangeVMHostname renames the VM (POST /api/service/:id/vms/:vmid/hostname).
func (c *Client) ChangeVMHostname(ctx context.Context, serviceID, vmID, hostname string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/hostname", Query("hostname", hostname))
	return err
}

// ListVMIPPools lists IP pools available to a VM
// (GET /api/service/:id/vms/:vmid/ippool).
func (c *Client) ListVMIPPools(ctx context.Context, serviceID, vmID, iface, bridge string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/ippool", Query("iface", iface, "bridge", bridge))
}

// AllocateVMIPs allocates IPs from a pool
// (POST /api/service/:id/vms/:vmid/ippool/:pool).
func (c *Client) AllocateVMIPs(ctx context.Context, serviceID, vmID, pool string, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/ippool/"+PathEscape(pool), qmap(params))
}

// ListServiceNetworks lists the networks available to a service
// (GET /api/service/:id/networks).
func (c *Client) ListServiceNetworks(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/networks", nil)
}

// ListVMIPs lists the IPs of a VM (GET /api/service/:id/vms/:vmid/ips).
func (c *Client) ListVMIPs(ctx context.Context, serviceID, vmID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/ips", nil)
}

// AssignVMIP assigns an IP to a VM (POST /api/service/:id/vms/:vmid/ips).
func (c *Client) AssignVMIP(ctx context.Context, serviceID, vmID, ipID, interfaceID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/ips",
		Query("ipid", ipID, "interfaceid", interfaceID))
	return err
}

// RemoveVMIP removes an IP from a VM (DELETE /api/service/:id/vms/:vmid/ips/:ipid).
func (c *Client) RemoveVMIP(ctx context.Context, serviceID, vmID, ipID, interfaceID string) error {
	_, err := c.deleteMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/ips/"+PathEscape(ipID),
		Query("interfaceid", interfaceID))
	return err
}

// ListAvailableVMIPs lists IPs available to a VM interface
// (GET /api/service/:id/vms/:vmid/interfaces/:iface/ips).
func (c *Client) ListAvailableVMIPs(ctx context.Context, serviceID, vmID, iface string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/interfaces/"+PathEscape(iface)+"/ips", nil)
}

// GetVMRDNS returns the rDNS of the VM's IPs
// (GET /api/service/:id/vms/:vmid/rdns).
func (c *Client) GetVMRDNS(ctx context.Context, serviceID, vmID string) ([]RDNSRecord, error) {
	var raw any
	if err := c.Get(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/rdns", nil, &raw); err != nil {
		return nil, err
	}
	return parseRDNS(raw), nil
}

// SetVMRDNS updates the rDNS of a VM IP
// (POST /api/service/:id/vms/:vmid/rdns, IP as key and hostname as value).
func (c *Client) SetVMRDNS(ctx context.Context, serviceID, vmID, ip, hostname string) error {
	q := Query()
	q.Set(ip, hostname)
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/rdns", q)
	return err
}

// RebuildVMNetwork rebuilds the VM network configuration
// (POST /api/service/:id/vms/:vmid/rebuild_network).
func (c *Client) RebuildVMNetwork(ctx context.Context, serviceID, vmID string, params map[string]string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/rebuild_network", qmap(params))
	return err
}

// VM usage graphs (folder "Cloud Instance" / "Cloud Virtual Machine").
// kind is one of cpu, net, disk, memory.

// GetVMUsageGraph returns a usage graph payload
// (GET /api/service/:id/vms/:vmid/usage/:kind).
func (c *Client) GetVMUsageGraph(ctx context.Context, serviceID, vmID, kind string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/usage/"+PathEscape(kind), nil)
}

// GetVMUsage returns overall VM usage (GET /api/service/:id/vms/:vmid/usage).
func (c *Client) GetVMUsage(ctx context.Context, serviceID, vmID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/usage", nil)
}

// ListVMDisks lists the disks attached to a VM
// (GET /api/service/:id/vms/:vmid/storage).
func (c *Client) ListVMDisks(ctx context.Context, serviceID, vmID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/storage", nil)
}

// ResizeVMDisk resizes a VM disk
// (PUT /api/service/:id/vms/:vmid/storage/:diskid).
func (c *Client) ResizeVMDisk(ctx context.Context, serviceID, vmID, diskID, size string) error {
	_, err := c.putMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/storage/"+PathEscape(diskID),
		Query("size", size))
	return err
}

// ListServiceImages lists ISO images of a service
// (GET /api/service/:id/images).
func (c *Client) ListServiceImages(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/images", nil)
}

// AddServiceImage registers an ISO image (POST /api/service/:id/images).
// params accepts: label, file_url, min_memory, version, os, distro,
// virtualization.
func (c *Client) AddServiceImage(ctx context.Context, serviceID string, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/images", qmap(params))
}

// ListVMImages lists the ISO images mountable on a VM
// (GET /api/service/:id/vms/:vmid/images).
func (c *Client) ListVMImages(ctx context.Context, serviceID, vmID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/images", nil)
}

// MountVMImage mounts an ISO image (POST /api/service/:id/vms/:vmid/images).
func (c *Client) MountVMImage(ctx context.Context, serviceID, vmID, image string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/images", Query("image", image))
	return err
}

// SetVMBootOrder sets the VM boot order (POST /api/service/:id/vms/:vmid/boot).
func (c *Client) SetVMBootOrder(ctx context.Context, serviceID, vmID, order string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/boot", Query("order", order))
	return err
}

// SetVMPXE toggles PXE boot (POST /api/service/:id/vms/:vmid/tuntap).
func (c *Client) SetVMPXE(ctx context.Context, serviceID, vmID, state string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/tuntap", Query("state", state))
	return err
}

// Upgrade service (folder "Cloud Virtual Machine").

// ListServiceUpgradeOptions lists upgrade options (GET /api/service/:id/upgrade).
func (c *Client) ListServiceUpgradeOptions(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/upgrade", nil)
}

// RequestServiceUpgrade requests a service upgrade
// (POST /api/service/:id/upgrade).
func (c *Client) RequestServiceUpgrade(ctx context.Context, serviceID string, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/upgrade", qmap(params))
}

// GetServiceResources returns service resource usage
// (GET /api/service/:id/resources).
func (c *Client) GetServiceResources(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/resources", nil)
}

// ListServiceVMTemplates lists rebuild templates for a VM via the service
// (GET /api/service/:id/templates/:vmid).
func (c *Client) ListServiceVMTemplates(ctx context.Context, serviceID, vmID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/templates/"+PathEscape(vmID), nil)
}

// vCloudStack (folder "vCloudStack Public Cloud").

// RescueVM boots a VM into rescue mode (POST /api/service/:id/vms/:vmid/rescue).
func (c *Client) RescueVM(ctx context.Context, serviceID, vmID string) (map[string]any, error) {
	return c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/rescue", nil)
}

// UnrescueVM leaves rescue mode (POST /api/service/:id/vms/:vmid/unrescue).
func (c *Client) UnrescueVM(ctx context.Context, serviceID, vmID string) error {
	_, err := c.postMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/unrescue", nil)
	return err
}

// GetVMConsole returns the console URL/token
// (GET /api/service/:id/vms/:vmid/console).
func (c *Client) GetVMConsole(ctx context.Context, serviceID, vmID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/vms/"+PathEscape(vmID)+"/console", nil)
}
