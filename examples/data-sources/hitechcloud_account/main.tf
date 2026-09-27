# Look up account and catalog data.

terraform {
  required_providers {
    hitechcloud = {
      source = "hitechcloud-vietnam/hitechcloud"
    }
  }
}

provider "hitechcloud" {}

data "hitechcloud_account" "current" {}

data "hitechcloud_services" "all" {}

data "hitechcloud_dns_zones" "all" {
  service_id = "1"
}

data "hitechcloud_ai_instance_types" "all" {
  service_id = "1"
}

output "account_email" {
  value = data.hitechcloud_account.current.email
}

output "service_names" {
  value = [for s in data.hitechcloud_services.all.services : s.name]
}
