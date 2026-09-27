# Manage S3 storage: a bucket and a sub-user with access keys.

terraform {
  required_providers {
    hitechcloud = {
      source = "hitechcloud-vietnam/hitechcloud"
    }
  }
}

provider "hitechcloud" {}

resource "hitechcloud_s3_bucket" "photos" {
  service_id = "1"
  name       = "photos"
  purge      = true
}

resource "hitechcloud_s3_subuser" "backup" {
  service_id = "1"
  name       = "backup-agent"
  access     = "readwrite"
}

output "s3_access_key" {
  value = hitechcloud_s3_subuser.backup.access_key
}

output "s3_secret_key" {
  value     = hitechcloud_s3_subuser.backup.secret_key
  sensitive = true
}
