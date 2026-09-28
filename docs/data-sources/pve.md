---
page_title: "hitechcloud_pve Data Source - hitechcloud"
subcategory: "Compute"
description: |-
  Reads a HiTechCloud Proxmox (PVE) service: status, VMs, IPs, backups, snapshots and usage.
---

# hitechcloud_pve (Data Source)

Reads a HiTechCloud Proxmox (PVE) service: status, VM inventory, IPs, backups,
snapshots and usage (`GET /api/service/{id}/htcpve/status`, `/htcpve/vms`,
`/htcpve/ips`, `/htcpve/backups`, `/htcpve/snapshots`, `/htcpve/usage`).

## Example Usage

```terraform
data "hitechcloud_pve" "compute" {
  service_id = "1"
}

output "running_vms" {
  value = [for vm in data.hitechcloud_pve.compute.vms : vm.name if vm.status == "running"]
}
```

## Schema

### Required

- `service_id` (String) HiTechCloud service ID of the Proxmox service.

### Read-Only

- `backups` (Attributes List) Backups of the service. (see
  [nested schema](#nestedatt--backups))
- `id` (String) The service ID (same as `service_id`).
- `ips` (Attributes List) IP addresses assigned to the service. (see
  [nested schema](#nestedatt--ips))
- `snapshots` (Attributes List) Snapshots of the machine, excluding the
  synthetic "current" entry. (see [nested schema](#nestedatt--snapshots))
- `status` (Map of String) Node-level status values.
- `usage` (Map of String) Billing-relevant usage counters with units.
- `vms` (Attributes List) Virtual machines of the service. (see
  [nested schema](#nestedatt--vms))

<a id="nestedatt--backups"></a>

### Nested Schema for `backups`

Read-Only:

- `created_at` (String) Backup timestamp.
- `id` (String) Backup identifier.
- `mode` (String) Backup mode (snapshot/...).
- `notes` (String) Backup notes.
- `size` (Number) Backup size in bytes.
- `vmid` (String) Backed-up VM ID.

<a id="nestedatt--ips"></a>

### Nested Schema for `ips`

Read-Only:

- `hostname` (String) Reverse DNS hostname.
- `ip` (String) IP address.
- `type` (String) Address type (ipv4/ipv6).

<a id="nestedatt--snapshots"></a>

### Nested Schema for `snapshots`

Read-Only:

- `created_at` (String) Snapshot timestamp.
- `description` (String) Snapshot description.
- `id` (String) Snapshot identifier.
- `name` (String) Snapshot name.

<a id="nestedatt--vms"></a>

### Nested Schema for `vms`

Read-Only:

- `cpus` (Number) Number of vCPUs.
- `disk_gb` (Number) Disk size in GB.
- `memory_mb` (Number) Memory in MB.
- `name` (String) VM name.
- `status` (String) VM status (running/stopped/...).
- `vmid` (String) Proxmox VM ID.
