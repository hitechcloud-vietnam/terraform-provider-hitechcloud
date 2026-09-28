# hitechcloud_pbs example

Reads a Proxmox Backup Server (PBS) service.

```terraform
terraform {
  required_providers {
    hitechcloud = {
      source = "hitechcloud-vietnam/hitechcloud"
    }
  }
}

provider "hitechcloud" {}

data "hitechcloud_pbs" "backup" {
  service_id = "1"
}

output "pbs_usage" {
  value = data.hitechcloud_pbs.backup.usage
}

output "snapshot_count" {
  value = length(data.hitechcloud_pbs.backup.snapshots)
}
