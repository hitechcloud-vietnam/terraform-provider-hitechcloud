---
page_title: "hitechcloud_vms Data Source - hitechcloud"
subcategory: "Compute"
description: |-
  Lists the virtual machines of a HiTechCloud hosted service.
---

# hitechcloud_vms (Data Source)

Lists the virtual machines of a HiTechCloud hosted service
(`GET /api/service/{service_id}/vm`).

## Example Usage

```terraform
data "hitechcloud_vms" "all" {
  service_id = "1"
}

output "vm_labels" {
  value = [for v in data.hitechcloud_vms.all.vms : v.label]
}
```

## Schema

### Required

- `service_id` (String) ID of the hosted service.

### Read-Only

- `id` (String) Data source identifier (same as `service_id`).
- `vms` (Attributes List) VMs of the service. (see
  [nested schema](#nestedatt--vms))

<a id="nestedatt--vms"></a>

### Nested Schema for `vms`

Read-Only:

- `cpu` (Number) Number of virtual CPUs.
- `hostname` (String) OS hostname.
- `id` (String) VM identifier.
- `ips` (Set of String) IP addresses of the VM.
- `label` (String) Display label.
- `memory` (Number) Memory in MB.
- `status` (String) VM status.
- `template_id` (String) OS template identifier.
