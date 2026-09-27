# Manage a DNS zone and its records.

terraform {
  required_providers {
    hitechcloud = {
      source = "hitechcloud-vietnam/hitechcloud"
    }
  }
}

provider "hitechcloud" {}

resource "hitechcloud_dns_zone" "example" {
  service_id = "1"
  name       = "example.com"
}

resource "hitechcloud_dns_record" "www" {
  service_id = "1"
  zone_id    = hitechcloud_dns_zone.example.zone_id
  name       = "www"
  type       = "A"
  content    = "203.0.113.10"
  ttl        = 3600
}

resource "hitechcloud_dns_record" "mx" {
  service_id = "1"
  zone_id    = hitechcloud_dns_zone.example.zone_id
  name       = "@"
  type       = "MX"
  content    = "mail.example.com"
  priority   = 10
}

output "zone_id" {
  value = hitechcloud_dns_zone.example.zone_id
}
