---
page_title: "hitechcloud_rdns Resource - hitechcloud"
subcategory: "Compute"
description: |-
  Manages the reverse DNS (PTR) record of a HiTechCloud service IP.
---

# hitechcloud_rdns (Resource)

Manages the reverse DNS (PTR) record of a HiTechCloud service IP
(`POST /api/service/{service_id}/rdns`). Destroying the resource clears the
reverse DNS entry.

## Example Usage

```terraform
resource "hitechcloud_rdns" "web" {
  service_id = "1"
  ip         = "203.0.113.10"
  hostname   = "web.example.com"
}
```

## Schema

### Required

- `hostname` (String) Reverse DNS hostname (PTR) of the IP address.
- `ip` (String) IP address to configure. Changing this forces a new resource.
- `service_id` (String) ID of the hosted service owning the IP. Changing this
  forces a new resource.

### Read-Only

- `id` (String) Composite identifier `service_id/ip`.

## Import

Import is supported using the composite ID `service_id/ip`:

```shell
terraform import hitechcloud_rdns.web 1/203.0.113.10
```
