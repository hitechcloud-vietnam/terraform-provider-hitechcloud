// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vdns_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/provider"
)

// mockDNSServer is an in-memory implementation of the tiny slice of the
// HiTechCloud API exercised by the DNS zone and DNS record resources.
type mockDNSServer struct {
	mu      sync.Mutex
	zones   map[string]map[string]any // zone id -> zone doc
	records map[string]map[string]any // record id -> record doc
	nextID  int
}

func newMockDNSServer() *mockDNSServer {
	return &mockDNSServer{
		zones:   map[string]map[string]any{},
		records: map[string]map[string]any{},
	}
}

func (m *mockDNSServer) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.mu.Lock()
		defer m.mu.Unlock()

		q := r.URL.Query()
		path := r.URL.Path

		switch {
		case r.Method == http.MethodPost && path == "/api/service/1/dns":
			m.nextID++
			id := fmt.Sprintf("zone-%d", m.nextID)
			zone := map[string]any{"id": id, "name": q.Get("name"), "records": []any{}}
			m.zones[id] = zone
			writeJSON(w, map[string]any{"id": id})

		case r.Method == http.MethodGet && regexp.MustCompile(`^/api/service/1/dns/[^/]+$`).MatchString(path):
			id := path[len("/api/service/1/dns/"):]
			zone, ok := m.zones[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			// Collect records of this zone.
			recs := []any{}
			for _, rec := range m.records {
				if rec["zone_id"] == id {
					recs = append(recs, rec)
				}
			}
			out := map[string]any{"id": id, "name": zone["name"], "records": recs}
			writeJSON(w, out)

		case r.Method == http.MethodGet && path == "/api/service/1/dns":
			list := []any{}
			for _, zone := range m.zones {
				list = append(list, zone)
			}
			writeJSON(w, map[string]any{"zones": list})

		case r.Method == http.MethodDelete && regexp.MustCompile(`^/api/service/1/dns/[^/]+$`).MatchString(path):
			id := path[len("/api/service/1/dns/"):]
			delete(m.zones, id)
			writeJSON(w, map[string]any{"status": "ok"})

		case r.Method == http.MethodPost && regexp.MustCompile(`^/api/service/1/dns/[^/]+/records$`).MatchString(path):
			zoneID := path[len("/api/service/1/dns/") : len(path)-len("/records")]
			// Auto-create the zone so record tests can run standalone.
			if _, ok := m.zones[zoneID]; !ok {
				m.zones[zoneID] = map[string]any{"id": zoneID, "name": zoneID + ".test"}
			}
			m.nextID++
			id := fmt.Sprintf("rec-%d", m.nextID)
			rec := map[string]any{
				"id":       id,
				"zone_id":  zoneID,
				"name":     q.Get("name"),
				"type":     q.Get("type"),
				"content":  q.Get("content"),
				"ttl":      3600,
				"priority": 10,
			}
			m.records[id] = rec
			writeJSON(w, map[string]any{"id": id})

		case r.Method == http.MethodGet && regexp.MustCompile(`^/api/service/1/dns/[^/]+/records/[^/]+$`).MatchString(path):
			id := path[strings.LastIndex(path, "/")+1:]
			rec, ok := m.records[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeJSON(w, rec)

		case r.Method == http.MethodPut && regexp.MustCompile(`^/api/service/1/dns/[^/]+/records/[^/]+$`).MatchString(path):
			id := path[strings.LastIndex(path, "/")+1:]
			rec, ok := m.records[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			rec["content"] = q.Get("content")
			writeJSON(w, map[string]any{"status": "ok"})

		case r.Method == http.MethodDelete && regexp.MustCompile(`^/api/service/1/dns/[^/]+/records/[^/]+$`).MatchString(path):
			id := path[strings.LastIndex(path, "/")+1:]
			delete(m.records, id)
			writeJSON(w, map[string]any{"status": "ok"})

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"hitechcloud": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

// TestAccDNSZone_basic exercises a full create/read/update-untouched/import/destroy
// lifecycle of hitechcloud_dns_zone against a mock HiTechCloud API. Runs only
// when TF_ACC is set (CI provides the Terraform CLI).
func TestAccDNSZone_basic(t *testing.T) {
	mock := newMockDNSServer()
	srv := httptest.NewServer(mock.handler())
	defer srv.Close()

	os.Setenv("HITECHCLOUD_ENDPOINT", srv.URL)
	os.Setenv("HITECHCLOUD_TOKEN", "test-token")
	defer os.Unsetenv("HITECHCLOUD_ENDPOINT")
	defer os.Unsetenv("HITECHCLOUD_TOKEN")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "hitechcloud" {}

resource "hitechcloud_dns_zone" "test" {
  service_id = "1"
  name       = "example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hitechcloud_dns_zone.test", "name", "example.com"),
					resource.TestCheckResourceAttrSet("hitechcloud_dns_zone.test", "zone_id"),
				),
			},
			{
				ResourceName:      "hitechcloud_dns_zone.test",
				ImportState:       true,
				ImportStateId:     "1/zone-1",
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccDNSRecord_updates exercises hitechcloud_dns_record create + in-place
// update against the mock API.
func TestAccDNSRecord_updates(t *testing.T) {
	mock := newMockDNSServer()
	srv := httptest.NewServer(mock.handler())
	defer srv.Close()

	os.Setenv("HITECHCLOUD_ENDPOINT", srv.URL)
	os.Setenv("HITECHCLOUD_TOKEN", "test-token")
	defer os.Unsetenv("HITECHCLOUD_ENDPOINT")
	defer os.Unsetenv("HITECHCLOUD_TOKEN")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "hitechcloud" {}

resource "hitechcloud_dns_record" "test" {
  service_id = "1"
  zone_id    = "zone-1"
  name       = "www"
  type       = "A"
  content    = "1.2.3.4"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hitechcloud_dns_record.test", "content", "1.2.3.4"),
				),
			},
			{
				Config: `
provider "hitechcloud" {}

resource "hitechcloud_dns_record" "test" {
  service_id = "1"
  zone_id    = "zone-1"
  name       = "www"
  type       = "A"
  content    = "5.6.7.8"
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hitechcloud_dns_record.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hitechcloud_dns_record.test", "content", "5.6.7.8"),
				),
			},
		},
	})
}
