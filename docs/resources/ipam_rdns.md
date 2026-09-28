---
page_title: "hitechcloud_ipam_rdns Resource - hitechcloud"
subcategory: "Network"
description: |-
  Manages the reverse DNS (PTR) record of an IPAM-managed address.
---

# hitechcloud_ipam_rdns (Resource)

Manages the reverse DNS (PTR) record of an IPAM-managed address
(`POST /api/service/{id}/htcipam/rdns`). Destroying the resource clears the PTR
record.

## Example Usage

```terraform
resource "hitechcloud_ipam_rdns" "gw" {
  service_id = "1"
  ip         = "203.0.113.20"
  hostname   = "gw.example.com"
}
```

## Schema

### Required

- `hostname` (String) Reverse DNS hostname (PTR) of the IP address; empty clears
  it.
- `ip` (String) IP address to configure, as text. Changing this forces a new
  resource.
- `service_id` (String) HiTechCloud service ID of the IPAM service. Changing this
  forces a new resource.

### Read-Only

- `id` (String) Composite identifier `service_id/ip`.

## Import

Import is supported using the composite ID `service_id/ip`:

```shell
terraform import hitechcloud_ipam_rdns.gw 1/203.0.113.20
```
