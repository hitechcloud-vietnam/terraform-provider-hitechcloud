// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// ---------------------------------------------------------------------------
// Coverage completions — endpoints present in the Postman collection that were
// not wrapped when the endpoint sweep was first done.
// ---------------------------------------------------------------------------

// ListAllDNS lists DNS entries across the whole account
// (GET /api/dns) — the account-wide counterpart of the per-service DNS API.
func (c *Client) ListAllDNS(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/dns", nil)
}

// ListServiceIPAddresses lists the IP addresses of a service using the
// singular "ip" endpoint (GET /api/service/{id}/ip) used by Network Services.
func (c *Client) ListServiceIPAddresses(ctx context.Context, serviceID string) ([]any, error) {
	return c.getList(ctx, "/api/service/"+PathEscape(serviceID)+"/ip", nil)
}

// GetServiceCluster returns the cluster assigned to a service
// (GET /api/service/{id}/cluster).
func (c *Client) GetServiceCluster(ctx context.Context, serviceID string) (map[string]any, error) {
	return c.getMap(ctx, "/api/service/"+PathEscape(serviceID)+"/cluster", nil)
}

// ListStatusesFiltered lists service statuses filtered by status value
// (GET /api/statuses?status=...). Pass an empty string for all statuses.
func (c *Client) ListStatusesFiltered(ctx context.Context, status string) ([]StatusEntry, error) {
	q := qmap(map[string]string{"status": status})
	var raw any
	if err := c.Get(ctx, "/api/statuses", q, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw)
	out := make([]StatusEntry, 0, len(list))
	for _, item := range list {
		m, ok := AsMap(item)
		if !ok {
			continue
		}
		out = append(out, StatusEntry{
			ID:     FirstString(m, "id", "service_id"),
			Name:   FirstString(m, "name", "title", "label"),
			Status: FirstString(m, "status", "state"),
			Type:   FirstString(m, "type", "group"),
		})
	}
	return out, nil
}
