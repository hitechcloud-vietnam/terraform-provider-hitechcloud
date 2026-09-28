---
page_title: "hitechcloud_pbs Data Source - hitechcloud"
subcategory: "Storage"
description: |-
  Reads a Proxmox Backup Server (PBS) service: connection info, credentials, usage and snapshots.
---

# hitechcloud_pbs (Data Source)

Reads a Proxmox Backup Server (PBS) service: connection info, credentials, usage
metrics and snapshots (`GET /api/service/{id}/pbs`, `/pbs/credentials`,
`/pbs/usage`, `/pbs/metrics`, `/pbs/snapshots`, `/pbs/groups`).

## Example Usage

```terraform
data "hitechcloud_pbs" "backup" {
  service_id = "1"
}

output "snapshot_count" {
  value = length(data.hitechcloud_pbs.backup.snapshots)
}
```

## Schema

### Required

- `service_id` (String) HiTechCloud service ID (`hb_accounts.id`) of the PBS
  service.

### Read-Only

- `endpoint` (String) PBS endpoint / hostname.
- `groups` (List of String) Backup groups present in the namespace.
- `id` (String) The service ID (same as `service_id`).
- `metrics` (Map of String) Raw metric values keyed by metric name.
- `namespace` (String) PBS namespace of the customer.
- `snapshots` (Attributes List) Backups in the customer's namespace (newest
  first). (see [nested schema](#nestedatt--snapshots))
- `usage` (Map of String) Billing-relevant usage counters (`backup_space`,
  `snapshots`, `backup_groups`) with units.
- `username` (String) PBS access username.

<a id="nestedatt--snapshots"></a>

### Nested Schema for `snapshots`

Read-Only:

- `backup_id` (String) Backup ID inside the group.
- `created_at` (String) Backup timestamp.
- `group` (String) Backup group.
- `id` (String) Snapshot identifier.
- `size` (Number) Snapshot size in bytes; null when PBS does not report a size.
