---
page_title: "hitechcloud_pve_rdns Resource - hitechcloud"
subcategory: "Compute"
description: |-
  Manages the reverse DNS (PTR) record of a HiTechCloud Proxmox service IP.
---

# hitechcloud_pve_rdns (Resource)

Manages the reverse DNS (PTR) record of a HiTechCloud Proxmox service IP
(`POST /api/service/{id}/htcpve/rdns`). Destroying the resource clears the PTR
record.

## Example Usage

```terraform
resource "hitechcloud_pve_rdns" "node" {
  service_id = "1"
  ip         = "203.0.113.10"
  hostname   = "node.example.com"
}
```

## Schema

### Required

- `hostname` (String) Reverse DNS hostname (PTR) of the IP address; empty clears
  it.
- `ip` (String) IP address to configure. Changing this forces a new resource.
- `service_id` (String) HiTechCloud service ID of the Proxmox service. Changing
  this forces a new resource.

### Read-Only

- `id` (String) Composite identifier `service_id/ip`.

## Import

Import is supported using the composite ID `service_id/ip`:

```shell
terraform import hitechcloud_pve_rdns.node 1/203.0.113.10
```
