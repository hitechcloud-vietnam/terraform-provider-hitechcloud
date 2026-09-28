---
page_title: "hitechcloud_vm_snapshot Resource - hitechcloud"
subcategory: "Compute"
description: |-
  Takes a snapshot of a HiTechCloud Proxmox virtual machine.
---

# hitechcloud_vm_snapshot (Resource)

Takes a snapshot of a HiTechCloud Proxmox virtual machine
(`POST /api/service/{id}/htcpve/snapshots`). Snapshots are immutable records:
changing any attribute forces a replacement.

The API does not expose a snapshot deletion endpoint — deleting the resource
removes it from Terraform state only and emits a warning.

## Example Usage

```terraform
resource "hitechcloud_vm_snapshot" "before_upgrade" {
  service_id  = "1"
  name        = "before-upgrade"
  description = "Snapshot taken before the OS upgrade"
}
```

## Schema

### Required

- `name` (String) Snapshot name; invalid characters are replaced by the API.
  Changing this forces a new resource.
- `service_id` (String) HiTechCloud service ID of the Proxmox VM. Changing this
  forces a new resource.

### Optional

- `description` (String) Free-text description of the snapshot. Changing this
  forces a new resource.

### Read-Only

- `created_at` (String) Creation timestamp reported by the API.
- `id` (String) Snapshot identifier returned by the API.

## Import

Import is supported using the composite ID `service_id/snapshot_id`:

```shell
terraform import hitechcloud_vm_snapshot.before_upgrade 1/snap-001
```
