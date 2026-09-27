// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// ---------------------------------------------------------------------------
// Virtual machines (folder "Cloud Instance": /api/service/:id/vms/...)
// ---------------------------------------------------------------------------

// VM describes a virtual machine (Cloud Instance).
type VM struct {
	ID         string
	Label      string
	Status     string
	Hostname   string
	Note       string
	TemplateID string
	Memory     int64 // MB
	CPU        int64 // cores
	CPUShare   int64 // percent
	Disk       int64 // GB
	Swap       int64 // GB
	IPs        []string
}

// VMCreate carries the parameters of POST /api/service/:id/vms.
type VMCreate struct {
	Label       string
	TemplateID  string
	Password    string
	Memory      int64
	CPU         int64
	CPUShare    int64
	Disk        int64
	Swap        int64
	Note        string
	LicenseKey  string
	LicenseType string
}

// ListVMs lists virtual servers under a service (GET /api/service/:id/vms).
func (c *Client) ListVMs(ctx context.Context, serviceID string) ([]VM, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/vms"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	return parseVMs(ExtractList(raw, "vms", "servers", "virtual_servers", "machines")), nil
}

// GetVM returns details of one VM (GET /api/service/:id/vms/:vmid).
func (c *Client) GetVM(ctx context.Context, serviceID, vmID string) (*VM, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/vms/" + PathEscape(vmID)
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "vms"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for VM %s", vmID)
	}
	vm := parseVM(m)
	if vm.ID == "" {
		vm.ID = vmID
	}
	return &vm, nil
}

// CreateVM provisions a new VM (POST /api/service/:id/vms) and resolves its id.
func (c *Client) CreateVM(ctx context.Context, serviceID string, in VMCreate) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/vms"
	q := Query(
		"label", in.Label,
		"template_id", in.TemplateID,
		"password", in.Password,
		"note", in.Note,
		"license_key", in.LicenseKey,
		"license_type", in.LicenseType,
	)
	if in.Memory > 0 {
		q.Set("memory", fmt.Sprintf("%d", in.Memory))
	}
	if in.CPU > 0 {
		q.Set("cpu", fmt.Sprintf("%d", in.CPU))
	}
	if in.CPUShare > 0 {
		q.Set("cpu_share", fmt.Sprintf("%d", in.CPUShare))
	}
	if in.Disk > 0 {
		q.Set("disk", fmt.Sprintf("%d", in.Disk))
	}
	if in.Swap > 0 {
		q.Set("swap", fmt.Sprintf("%d", in.Swap))
	}
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "vm_id", "vmid", "virtual_id", "server_id"); id != "" {
		return id, nil
	}
	// Fallback: locate the VM by label.
	vms, err := c.ListVMs(ctx, serviceID)
	if err != nil {
		return "", fmt.Errorf("VM was created but its id could not be resolved: %w", err)
	}
	for _, vm := range vms {
		if in.Label != "" && vm.Label == in.Label {
			return vm.ID, nil
		}
	}
	return "", fmt.Errorf("VM %q was created but could not be found in the VM list", in.Label)
}

// UpdateVM resizes a VM (PUT /api/service/:id/vms/:vmid).
func (c *Client) UpdateVM(ctx context.Context, serviceID, vmID string, memory, cpu, cpuShare int64) error {
	path := "/api/service/" + PathEscape(serviceID) + "/vms/" + PathEscape(vmID)
	q := url.Values{}
	if memory > 0 {
		q.Set("memory", fmt.Sprintf("%d", memory))
	}
	if cpu > 0 {
		q.Set("cpu", fmt.Sprintf("%d", cpu))
	}
	if cpuShare > 0 {
		q.Set("cpu_share", fmt.Sprintf("%d", cpuShare))
	}
	return c.Put(ctx, path, q, nil)
}

// DeleteVM removes a VM (DELETE /api/service/:id/vms/:vmid).
func (c *Client) DeleteVM(ctx context.Context, serviceID, vmID string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/vms/" + PathEscape(vmID)
	return c.Delete(ctx, path, nil, nil)
}

// WaitForVMDeleted polls until the VM is gone (404) or the timeout elapses.
// Virtual server removal is not guaranteed to be synchronous, so Terraform
// deletes wait for a terminal state to avoid racing dependent resources.
func (c *Client) WaitForVMDeleted(ctx context.Context, serviceID, vmID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		_, err := c.GetVM(ctx, serviceID, vmID)
		if IsNotFound(err) {
			return nil
		}
		if err == nil && time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for VM %s deletion", vmID)
		}
		if err != nil && !IsNotFound(err) && time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for VM %s deletion: %w", vmID, err)
		}
		if err := sleepCtx(ctx, 5*time.Second); err != nil {
			return err
		}
	}
}

func parseVMs(list []any) []VM {
	out := make([]VM, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseVM(m))
		}
	}
	return out
}

func parseVM(m map[string]any) VM {
	return VM{
		ID:         FirstString(m, "id", "vm_id", "vmid", "virtual_id", "server_id"),
		Label:      FirstString(m, "label", "name", "title"),
		Status:     FirstString(m, "status", "state"),
		Hostname:   FirstString(m, "hostname", "host_name", "host"),
		Note:       FirstString(m, "note", "notes", "description"),
		TemplateID: asString(First(m, "template_id", "template", "os_template")),
		Memory:     FirstInt64(m, "memory", "ram", "memory_mb"),
		CPU:        FirstInt64(m, "cpu", "cpus", "cpu_cores", "cores"),
		CPUShare:   FirstInt64(m, "cpu_share", "cpushare", "cpu_shares"),
		Disk:       FirstInt64(m, "disk", "disk_gb", "disk_space", "hdd"),
		Swap:       FirstInt64(m, "swap", "swap_gb"),
		IPs:        StringList(First(m, "ips", "ip_addresses", "ip", "addresses")),
	}
}

// ---------------------------------------------------------------------------
// VM network interfaces
// (folder "Cloud Service": /api/service/:id/vms/:vmid/interfaces/...)
// ---------------------------------------------------------------------------

// VMInterface is a virtual machine network interface.
type VMInterface struct {
	ID       string
	Bridge   string
	Firewall bool
	IPv4IDs  []string
	IPv6IDs  []string
}

// ListVMInterfaces lists interfaces of a VM
// (GET /api/service/:id/vms/:vmid/interfaces).
func (c *Client) ListVMInterfaces(ctx context.Context, serviceID, vmID string) ([]VMInterface, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/vms/" + PathEscape(vmID) + "/interfaces"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "interfaces", "nics", "network_interfaces")
	out := make([]VMInterface, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseVMInterface(m))
		}
	}
	return out, nil
}

// GetVMInterface returns one interface
// (GET /api/service/:id/vms/:vmid/interfaces/:iface).
func (c *Client) GetVMInterface(ctx context.Context, serviceID, vmID, iface string) (*VMInterface, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/vms/" + PathEscape(vmID) + "/interfaces/" + PathEscape(iface)
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "interfaces"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for interface %s", iface)
	}
	nic := parseVMInterface(m)
	if nic.ID == "" {
		nic.ID = iface
	}
	return &nic, nil
}

// CreateVMInterface adds a network interface
// (POST /api/service/:id/vms/:vmid/interfaces).
func (c *Client) CreateVMInterface(ctx context.Context, serviceID, vmID string, in VMInterface) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/vms/" + PathEscape(vmID) + "/interfaces"
	q := url.Values{}
	if in.Bridge != "" {
		q.Set("bridge", in.Bridge)
	}
	q.Set("firewall", fmt.Sprintf("%t", in.Firewall))
	for _, id := range in.IPv4IDs {
		q.Add("ipv4", id)
	}
	for _, id := range in.IPv6IDs {
		q.Add("ipv6", id)
	}
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "iface", "interface", "interface_id", "name"); id != "" {
		return id, nil
	}
	// Fallback: locate by bridge.
	ifaces, err := c.ListVMInterfaces(ctx, serviceID, vmID)
	if err != nil {
		return "", fmt.Errorf("interface was created but its id could not be resolved: %w", err)
	}
	for _, nic := range ifaces {
		if nic.Bridge == in.Bridge {
			return nic.ID, nil
		}
	}
	return "", fmt.Errorf("network interface on bridge %q was created but could not be found", in.Bridge)
}

// UpdateVMInterface changes firewall/IP assignment
// (PUT /api/service/:id/vms/:vmid/interfaces/:iface).
func (c *Client) UpdateVMInterface(ctx context.Context, serviceID, vmID, iface string, in VMInterface) error {
	path := "/api/service/" + PathEscape(serviceID) + "/vms/" + PathEscape(vmID) + "/interfaces/" + PathEscape(iface)
	q := url.Values{}
	q.Set("firewall", fmt.Sprintf("%t", in.Firewall))
	for _, id := range in.IPv4IDs {
		q.Add("ipv4", id)
	}
	for _, id := range in.IPv6IDs {
		q.Add("ipv6", id)
	}
	return c.Put(ctx, path, q, nil)
}

// DeleteVMInterface removes an interface
// (DELETE /api/service/:id/vms/:vmid/interfaces/:iface).
func (c *Client) DeleteVMInterface(ctx context.Context, serviceID, vmID, iface string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/vms/" + PathEscape(vmID) + "/interfaces/" + PathEscape(iface)
	return c.Delete(ctx, path, nil, nil)
}

func parseVMInterface(m map[string]any) VMInterface {
	return VMInterface{
		ID:       FirstString(m, "id", "iface", "interface", "interface_id", "name"),
		Bridge:   FirstString(m, "bridge", "bridge_name", "network"),
		Firewall: FirstBool(m, "firewall", "firewall_enabled"),
		IPv4IDs:  StringList(First(m, "ipv4", "ip4", "ipv4_ids", "ipv4_id")),
		IPv6IDs:  StringList(First(m, "ipv6", "ip6", "ipv6_ids", "ipv6_id")),
	}
}

// ---------------------------------------------------------------------------
// VM firewall rules (folder "Cloud GPU": /api/service/:id/vms/firewall)
// ---------------------------------------------------------------------------

// FirewallRule is a virtual server firewall rule.
type FirewallRule struct {
	Position     int64
	Action       string // accept | drop
	Type         string // source | destination | both
	Comment      string
	Protocol     string // tcp | udp | icmp
	AddressStart string
	AddressEnd   string
	PortStart    int64
	PortEnd      int64
}

// ListFirewallRules lists firewall rules (GET /api/service/:id/vms/firewall).
func (c *Client) ListFirewallRules(ctx context.Context, serviceID string) ([]FirewallRule, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/vms/firewall"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "rules", "firewall", "firewall_rules")
	out := make([]FirewallRule, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseFirewallRule(m))
		}
	}
	return out, nil
}

// CreateFirewallRule adds a rule (POST /api/service/:id/vms/firewall) and
// resolves the rule position for later deletion.
func (c *Client) CreateFirewallRule(ctx context.Context, serviceID string, rule FirewallRule) (int64, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/vms/firewall"
	q := Query(
		"action", rule.Action,
		"type", rule.Type,
		"comment", rule.Comment,
		"protocol", rule.Protocol,
		"address_start", rule.AddressStart,
		"address_end", rule.AddressEnd,
	)
	if rule.PortStart > 0 {
		q.Set("port_start", fmt.Sprintf("%d", rule.PortStart))
	}
	if rule.PortEnd > 0 {
		q.Set("port_end", fmt.Sprintf("%d", rule.PortEnd))
	}
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return 0, err
	}
	if pos := FirstInt64(raw, "position", "id", "index", "rule_id"); pos > 0 {
		return pos, nil
	}
	// Fallback: match the rule in the listing.
	rules, err := c.ListFirewallRules(ctx, serviceID)
	if err != nil {
		return 0, fmt.Errorf("firewall rule was created but its position could not be resolved: %w", err)
	}
	for _, r := range rules {
		if firewallRuleEqual(r, rule) {
			return r.Position, nil
		}
	}
	return 0, fmt.Errorf("firewall rule %s/%s was created but could not be found", rule.Protocol, rule.Action)
}

// DeleteFirewallRule removes a rule by position
// (DELETE /api/service/:id/vms/firewall/:position).
func (c *Client) DeleteFirewallRule(ctx context.Context, serviceID string, position int64) error {
	path := "/api/service/" + PathEscape(serviceID) + "/vms/firewall/" + fmt.Sprintf("%d", position)
	return c.Delete(ctx, path, nil, nil)
}

func firewallRuleEqual(a, b FirewallRule) bool {
	return a.Action == b.Action &&
		a.Type == b.Type &&
		a.Protocol == b.Protocol &&
		a.AddressStart == b.AddressStart &&
		a.AddressEnd == b.AddressEnd &&
		a.PortStart == b.PortStart &&
		a.PortEnd == b.PortEnd
}

func parseFirewallRule(m map[string]any) FirewallRule {
	return FirewallRule{
		Position:     FirstInt64(m, "position", "id", "index", "rule_id"),
		Action:       FirstString(m, "action"),
		Type:         FirstString(m, "type", "direction"),
		Comment:      FirstString(m, "comment", "description"),
		Protocol:     FirstString(m, "protocol", "proto"),
		AddressStart: FirstString(m, "address_start", "addressstart", "src_start", "ip_start"),
		AddressEnd:   FirstString(m, "address_end", "addressend", "src_end", "ip_end"),
		PortStart:    FirstInt64(m, "port_start", "portstart", "from_port"),
		PortEnd:      FirstInt64(m, "port_end", "portend", "to_port"),
	}
}

// ---------------------------------------------------------------------------
// Bare metal / colocation service IPs
// (folder "Bare Metal & Colocation": /api/service/:id/ips/...)
// ---------------------------------------------------------------------------

// ServiceIP is an IP address attached to a bare metal service.
type ServiceIP struct {
	ID     string
	IP     string
	Domain string
	VLAN   string
}

// ListServiceIPs lists server IPs (GET /api/service/:id/ips).
func (c *Client) ListServiceIPs(ctx context.Context, serviceID string) ([]ServiceIP, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/ips"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "ips", "ip", "ip_addresses", "addresses")
	out := make([]ServiceIP, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseServiceIP(m))
		}
	}
	return out, nil
}

// GetServiceIP returns IP details (GET /api/service/:id/ips/:ip).
func (c *Client) GetServiceIP(ctx context.Context, serviceID, ipID string) (*ServiceIP, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/ips/" + PathEscape(ipID)
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "ips"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for service IP %s", ipID)
	}
	ip := parseServiceIP(m)
	if ip.ID == "" {
		ip.ID = ipID
	}
	return &ip, nil
}

// CreateServiceIP adds an IP to a server (POST /api/service/:id/ips) and
// resolves the new IP record id. The `num` parameter is intentionally not
// exposed: one Terraform resource manages exactly one IP address.
func (c *Client) CreateServiceIP(ctx context.Context, serviceID, vlan, domain string) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/ips"
	if err := c.Post(ctx, path, Query("vlan", vlan, "domain", domain), &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "ip_id", "ipid"); id != "" {
		return id, nil
	}
	// Fallback: match by domain in the listing.
	ips, err := c.ListServiceIPs(ctx, serviceID)
	if err != nil {
		return "", fmt.Errorf("service IP was created but its id could not be resolved: %w", err)
	}
	for _, ip := range ips {
		if domain != "" && ip.Domain == domain {
			return ip.ID, nil
		}
	}
	return "", fmt.Errorf("service IP was created but could not be found in the IP list")
}

// UpdateServiceIP updates the IP domain (POST /api/service/:id/ips/:ip).
func (c *Client) UpdateServiceIP(ctx context.Context, serviceID, ipID, domain string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/ips/" + PathEscape(ipID)
	return c.Post(ctx, path, Query("domain", domain), nil)
}

// DeleteServiceIP removes an IP (DELETE /api/service/:id/ips/:ip).
func (c *Client) DeleteServiceIP(ctx context.Context, serviceID, ipID string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/ips/" + PathEscape(ipID)
	return c.Delete(ctx, path, nil, nil)
}

func parseServiceIP(m map[string]any) ServiceIP {
	return ServiceIP{
		ID:     FirstString(m, "id", "ip_id", "ipid"),
		IP:     FirstString(m, "ip", "ip_address", "ipaddress", "address"),
		Domain: FirstString(m, "domain", "ptr", "rdns", "hostname"),
		VLAN:   FirstString(m, "vlan", "vlan_id", "vlanid"),
	}
}
