---
page_title: "hitechcloud_pve_backup Resource - hitechcloud"
subcategory: "Compute"
description: |-
  Takes a backup of a HiTechCloud Proxmox service.
---

# hitechcloud_pve_backup (Resource)

Takes a backup of a HiTechCloud Proxmox service
(`POST /api/service/{id}/htcpve/backups`). Backups are immutable records:
changing any attribute forces a replacement.

The API has no backup deletion endpoint — deleting the resource removes it from
Terraform state only and emits a warning.

## Example Usage

```terraform
resource "hitechcloud_pve_backup" "nightly" {
  service_id = "1"
  mode       = "snapshot"
  notes      = "nightly backup"
}
```

## Schema

### Required

- `service_id` (String) HiTechCloud service ID of the Proxmox service. Changing
  this forces a new resource.

### Optional

- `mode` (String) Backup mode, for example `snapshot`. Changing this forces a new
  resource.
- `notes` (String) Free-text notes for the backup. Changing this forces a new
  resource.

### Read-Only

- `created_at` (String) Creation timestamp reported by the API.
- `id` (String) Backup identifier returned by the API.

## Import

Import is supported using the composite ID `service_id/backup_id`:

```shell
terraform import hitechcloud_pve_backup.nightly 1/backup-001
```
