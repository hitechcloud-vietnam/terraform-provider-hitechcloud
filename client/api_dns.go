// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// ---------------------------------------------------------------------------
// Service DNS zones and records
// (folder "DNS": /api/service/:service_id/dns/...)
// ---------------------------------------------------------------------------

// DNSZone is a DNS zone attached to a DNS service.
type DNSZone struct {
	ID      string
	Name    string
	Records []DNSRecord
}

// DNSRecord is a record inside a service DNS zone.
type DNSRecord struct {
	ID       string
	Name     string
	Type     string
	Content  string
	TTL      int64
	Priority int64
}

// ListDNSZones returns the DNS zones under a service
// (GET /api/service/:service_id/dns).
func (c *Client) ListDNSZones(ctx context.Context, serviceID string) ([]DNSZone, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/dns"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "zones", "dns", "dns_zones")
	out := make([]DNSZone, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseDNSZone(m))
		}
	}
	return out, nil
}

// GetDNSZone returns details of a DNS zone, including its records
// (GET /api/service/:service_id/dns/:zone_id).
func (c *Client) GetDNSZone(ctx context.Context, serviceID, zoneID string) (*DNSZone, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/dns/" + PathEscape(zoneID)
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "zones"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for DNS zone %s", zoneID)
	}
	zone := parseDNSZone(m)
	if zone.ID == "" {
		zone.ID = zoneID
	}
	return &zone, nil
}

// CreateDNSZone creates a new DNS zone (POST /api/service/:service_id/dns).
// It returns the new zone id; when the response does not carry an id the zone
// is located in the zone list by name.
func (c *Client) CreateDNSZone(ctx context.Context, serviceID, name string) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/dns"
	if err := c.Post(ctx, path, Query("name", name), &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "zone_id", "zoneid", "dns_id"); id != "" {
		return id, nil
	}
	// Fallback: locate the freshly created zone by name.
	zones, err := c.ListDNSZones(ctx, serviceID)
	if err != nil {
		return "", fmt.Errorf("DNS zone was created but its id could not be resolved: %w", err)
	}
	for _, z := range zones {
		if z.Name == name {
			return z.ID, nil
		}
	}
	return "", fmt.Errorf("DNS zone %q was created but could not be found in the zone list", name)
}

// DeleteDNSZone removes a DNS zone (DELETE /api/service/:service_id/dns/:zone_id).
func (c *Client) DeleteDNSZone(ctx context.Context, serviceID, zoneID string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/dns/" + PathEscape(zoneID)
	return c.Delete(ctx, path, nil, nil)
}

func parseDNSZone(m map[string]any) DNSZone {
	zone := DNSZone{
		ID:   FirstString(m, "id", "zone_id", "zoneid", "dns_id"),
		Name: FirstString(m, "name", "zone", "domain", "zone_name"),
	}
	for _, rec := range parseDNSRecords(ExtractList(m, "records", "rows")) {
		zone.Records = append(zone.Records, rec)
	}
	return zone
}

func parseDNSRecords(list []any) []DNSRecord {
	out := make([]DNSRecord, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, DNSRecord{
				ID:       FirstString(m, "id", "record_id", "recordid", "index"),
				Name:     FirstString(m, "name", "host", "hostname", "record_name"),
				Type:     FirstString(m, "type", "record_type", "rrtype"),
				Content:  FirstString(m, "content", "value", "target", "points_to", "data"),
				TTL:      FirstInt64(m, "ttl", "ttl_value"),
				Priority: FirstInt64(m, "priority", "prio", "priority_value"),
			})
		}
	}
	return out
}

// ListDNSRecords returns the records of a DNS zone.
func (c *Client) ListDNSRecords(ctx context.Context, serviceID, zoneID string) ([]DNSRecord, error) {
	zone, err := c.GetDNSZone(ctx, serviceID, zoneID)
	if err != nil {
		return nil, err
	}
	return zone.Records, nil
}

// FindDNSRecord locates a single record inside a zone by id.
func (c *Client) FindDNSRecord(ctx context.Context, serviceID, zoneID, recordID string) (*DNSRecord, error) {
	records, err := c.ListDNSRecords(ctx, serviceID, zoneID)
	if err != nil {
		return nil, err
	}
	for _, rec := range records {
		if rec.ID == recordID {
			r := rec
			return &r, nil
		}
	}
	return nil, &APIError{StatusCode: 404, Status: "Not Found", Message: fmt.Sprintf("DNS record %s not found", recordID)}
}

// CreateDNSRecord adds a record to a DNS zone
// (POST /api/service/:service_id/dns/:zone_id/records).
func (c *Client) CreateDNSRecord(ctx context.Context, serviceID, zoneID string, rec DNSRecord) (string, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/dns/" + PathEscape(zoneID) + "/records"
	q := Query(
		"name", rec.Name,
		"type", rec.Type,
		"content", rec.Content,
	)
	if rec.TTL > 0 {
		q.Set("ttl", fmt.Sprintf("%d", rec.TTL))
	}
	if rec.Priority > 0 {
		q.Set("priority", fmt.Sprintf("%d", rec.Priority))
	}
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "record_id", "recordid"); id != "" {
		return id, nil
	}
	// Fallback: match the new record in the zone listing.
	records, err := c.ListDNSRecords(ctx, serviceID, zoneID)
	if err != nil {
		return "", fmt.Errorf("DNS record was created but its id could not be resolved: %w", err)
	}
	for _, r := range records {
		if r.Name == rec.Name && r.Type == rec.Type && r.Content == rec.Content {
			return r.ID, nil
		}
	}
	return "", fmt.Errorf("DNS record %s/%s was created but could not be found in the zone", rec.Type, rec.Name)
}

// UpdateDNSRecord edits an existing record
// (PUT /api/service/:service_id/dns/:zone_id/records/:record_id).
func (c *Client) UpdateDNSRecord(ctx context.Context, serviceID, zoneID, recordID string, rec DNSRecord) error {
	path := "/api/service/" + PathEscape(serviceID) + "/dns/" + PathEscape(zoneID) + "/records/" + PathEscape(recordID)
	q := Query(
		"name", rec.Name,
		"type", rec.Type,
		"content", rec.Content,
	)
	if rec.TTL > 0 {
		q.Set("ttl", fmt.Sprintf("%d", rec.TTL))
	}
	if rec.Priority > 0 {
		q.Set("priority", fmt.Sprintf("%d", rec.Priority))
	}
	return c.Put(ctx, path, q, nil)
}

// DeleteDNSRecord removes a record
// (DELETE /api/service/:service_id/dns/:zone_id/records/:record_id).
func (c *Client) DeleteDNSRecord(ctx context.Context, serviceID, zoneID, recordID string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/dns/" + PathEscape(zoneID) + "/records/" + PathEscape(recordID)
	return c.Delete(ctx, path, nil, nil)
}

// ---------------------------------------------------------------------------
// Registered-domain DNS records
// (folder "DNS Manage": /api/domain/:id/dns/...)
// ---------------------------------------------------------------------------

// DomainDNSRecord is a DNS record of a registered domain.
type DomainDNSRecord struct {
	ID       string // record index used in /api/domain/:id/dns/:index
	Name     string
	Type     string
	Content  string
	Priority int64
}

// ListDomainDNSRecords returns DNS records of a registered domain
// (GET /api/domain/:id/dns).
func (c *Client) ListDomainDNSRecords(ctx context.Context, domainID string) ([]DomainDNSRecord, error) {
	var raw any
	path := "/api/domain/" + PathEscape(domainID) + "/dns"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "records", "dns", "dns_records")
	out := make([]DomainDNSRecord, 0, len(list))
	for _, item := range list {
		m, ok := AsMap(item)
		if !ok {
			continue
		}
		out = append(out, DomainDNSRecord{
			ID:       FirstString(m, "id", "record_id", "recordid", "index", "idx"),
			Name:     FirstString(m, "name", "host", "hostname", "record_name"),
			Type:     FirstString(m, "type", "record_type", "rrtype"),
			Content:  FirstString(m, "content", "value", "target", "data"),
			Priority: FirstInt64(m, "priority", "prio"),
		})
	}
	return out, nil
}

// GetDomainDNSRecord locates one record by its index/id.
func (c *Client) GetDomainDNSRecord(ctx context.Context, domainID, recordID string) (*DomainDNSRecord, error) {
	records, err := c.ListDomainDNSRecords(ctx, domainID)
	if err != nil {
		return nil, err
	}
	for _, rec := range records {
		if rec.ID == recordID {
			r := rec
			return &r, nil
		}
	}
	return nil, &APIError{StatusCode: 404, Status: "Not Found", Message: fmt.Sprintf("domain DNS record %s not found", recordID)}
}

// CreateDomainDNSRecord adds a record (POST /api/domain/:id/dns) and returns
// its index/id.
func (c *Client) CreateDomainDNSRecord(ctx context.Context, domainID string, rec DomainDNSRecord) (string, error) {
	var raw any
	path := "/api/domain/" + PathEscape(domainID) + "/dns"
	q := Query(
		"name", rec.Name,
		"type", rec.Type,
		"content", rec.Content,
	)
	if rec.Priority > 0 {
		q.Set("priority", fmt.Sprintf("%d", rec.Priority))
	}
	if err := c.Post(ctx, path, q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(raw, "id", "record_id", "recordid", "index", "idx"); id != "" {
		return id, nil
	}
	// Fallback: match the record in the listing.
	records, err := c.ListDomainDNSRecords(ctx, domainID)
	if err != nil {
		return "", fmt.Errorf("domain DNS record was created but its id could not be resolved: %w", err)
	}
	for _, r := range records {
		if r.Name == rec.Name && r.Type == rec.Type && r.Content == rec.Content {
			return r.ID, nil
		}
	}
	return "", fmt.Errorf("domain DNS record %s/%s was created but could not be found", rec.Type, rec.Name)
}

// UpdateDomainDNSRecord changes a record (PUT /api/domain/:id/dns/:index).
// The API accepts the record index both as path parameter `index` and query
// parameter `record_id`; both are sent for compatibility.
func (c *Client) UpdateDomainDNSRecord(ctx context.Context, domainID, recordID string, rec DomainDNSRecord) error {
	path := "/api/domain/" + PathEscape(domainID) + "/dns/" + PathEscape(recordID)
	q := Query(
		"record_id", recordID,
		"name", rec.Name,
		"type", rec.Type,
		"content", rec.Content,
	)
	if rec.Priority > 0 {
		q.Set("priority", fmt.Sprintf("%d", rec.Priority))
	}
	return c.Put(ctx, path, q, nil)
}

// DeleteDomainDNSRecord removes a record (DELETE /api/domain/:id/dns/:index).
func (c *Client) DeleteDomainDNSRecord(ctx context.Context, domainID, recordID string) error {
	path := "/api/domain/" + PathEscape(domainID) + "/dns/" + PathEscape(recordID)
	return c.Delete(ctx, path, Query("record_id", recordID), nil)
}

// ListDomainDNSTypes returns the record types supported for a domain
// (GET /api/domain/:id/dns/types).
func (c *Client) ListDomainDNSTypes(ctx context.Context, domainID string) ([]string, error) {
	var raw any
	path := "/api/domain/" + PathEscape(domainID) + "/dns/types"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	return StringList(raw), nil
}

// ---------------------------------------------------------------------------
// Domain settings (nameservers, auto-renew, registrar lock, ID protection)
// ---------------------------------------------------------------------------

// GetDomainNameservers returns the current nameservers (GET /api/domain/:id/ns).
func (c *Client) GetDomainNameservers(ctx context.Context, domainID string) ([]string, error) {
	var raw any
	path := "/api/domain/" + PathEscape(domainID) + "/ns"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	return StringList(First(raw, "nameservers", "ns", "nameserver", "list")), nil
}

// SetDomainNameservers changes nameservers (PUT /api/domain/:id/ns). Passing an
// empty list instructs the API to use the default nameservers.
func (c *Client) SetDomainNameservers(ctx context.Context, domainID string, nameservers []string) error {
	path := "/api/domain/" + PathEscape(domainID) + "/ns"
	q := url.Values{}
	for _, ns := range nameservers {
		q.Add("nameservers", ns)
	}
	return c.Put(ctx, path, q, nil)
}

// GetDomainAutorenew reads the auto-renew flag (GET /api/domain/:id/autorenew).
func (c *Client) GetDomainAutorenew(ctx context.Context, domainID string) (bool, error) {
	var raw any
	path := "/api/domain/" + PathEscape(domainID) + "/autorenew"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return false, err
	}
	return FirstBool(raw, "autorenew", "auto_renew", "status", "enabled"), nil
}

// SetDomainAutorenew toggles auto-renew (PUT /api/domain/:id/autorenew).
func (c *Client) SetDomainAutorenew(ctx context.Context, domainID string, enabled bool) error {
	path := "/api/domain/" + PathEscape(domainID) + "/autorenew"
	return c.Put(ctx, path, Query("autorenew", fmt.Sprintf("%t", enabled)), nil)
}

// GetDomainReglock reads the registrar lock (GET /api/domain/:id/reglock).
func (c *Client) GetDomainReglock(ctx context.Context, domainID string) (bool, error) {
	var raw any
	path := "/api/domain/" + PathEscape(domainID) + "/reglock"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return false, err
	}
	return FirstBool(raw, "reglock", "registrar_lock", "transfer_lock", "locked", "status", "enabled"), nil
}

// SetDomainReglock toggles the registrar lock (PUT /api/domain/:id/reglock).
func (c *Client) SetDomainReglock(ctx context.Context, domainID string, enabled bool) error {
	path := "/api/domain/" + PathEscape(domainID) + "/reglock"
	return c.Put(ctx, path, Query("switch", fmt.Sprintf("%t", enabled)), nil)
}

// SetDomainIDProtection toggles ID protection (PUT /api/domain/:id/idprotection).
// The API offers no dedicated read endpoint; the flag is read from domain
// details.
func (c *Client) SetDomainIDProtection(ctx context.Context, domainID string, enabled bool) error {
	path := "/api/domain/" + PathEscape(domainID) + "/idprotection"
	return c.Put(ctx, path, Query("switch", fmt.Sprintf("%t", enabled)), nil)
}
