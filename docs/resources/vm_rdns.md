---
page_title: "hitechcloud_vm_rdns Resource - hitechcloud"
subcategory: "Compute"
description: |-
  Manages the reverse DNS (PTR) record of a VM IP address.
---

# hitechcloud_vm_rdns (Resource)

Manages the reverse DNS (PTR) record of a VM IP address
(`POST /api/service/{id}/vms/{vmid}/rdns`). Destroying the resource clears the
PTR record.

## Example Usage

```terraform
resource "hitechcloud_vm_rdns" "web" {
  service_id = "1"
  vm_id      = "101"
  ip         = "203.0.113.30"
  hostname   = "web.example.com"
}
```

## Schema

### Required

- `hostname` (String) Reverse DNS hostname (PTR) of the IP address; empty clears
  it.
- `ip` (String) IP address to configure. Changing this forces a new resource.
- `service_id` (String) ID of the hosted service owning the VM. Changing this
  forces a new resource.
- `vm_id` (String) VM identifier. Changing this forces a new resource.

### Read-Only

- `id` (String) Composite identifier `service_id/vm_id/ip`.

## Import

Import is supported using the composite ID `service_id/vm_id/ip`:

```shell
terraform import hitechcloud_vm_rdns.web 1/101/203.0.113.30
```
