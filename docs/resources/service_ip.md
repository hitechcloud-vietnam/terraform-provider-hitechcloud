---
page_title: "hitechcloud_service_ip Resource - hitechcloud"
subcategory: "Compute"
description: |-
  Manages an additional IP address of a HiTechCloud service.
---

# hitechcloud_service_ip (Resource)

Manages an additional IP address of a HiTechCloud service
(`/api/service/{service_id}/ip`).

## Example Usage

```terraform
resource "hitechcloud_service_ip" "extra" {
  service_id = "1"
  vlan       = "100"
  domain     = "example.com"
}
```

## Schema

### Required

- `service_id` (String) ID of the hosted service. Changing this forces a new
  resource.

### Optional

- `domain` (String) Reverse DNS domain of the IP.
- `vlan` (String) VLAN of the IP. Defaults to the API default. Changing this
  forces a new resource.

### Read-Only

- `id` (String) Composite identifier `service_id/ip_id`.
- `ip` (String) The IP address.
- `ip_id` (String) Identifier of the IP as returned by the API.

## Import

Import is supported using the composite ID `service_id/ip_id`:

```shell
terraform import hitechcloud_service_ip.extra 1/ip-123
```
