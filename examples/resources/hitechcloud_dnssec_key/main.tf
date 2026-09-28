# hitechcloud_dnssec_key example

Manages a DNSSEC (DS) key of a registered domain.

```terraform
terraform {
  required_providers {
    hitechcloud = {
      source = "hitechcloud-vietnam/hitechcloud"
    }
  }
}

provider "hitechcloud" {}

resource "hitechcloud_dnssec_key" "example" {
  domain_id   = "42"
  key_tag     = "2371"
  algorithm   = "13"
  digest_type = "2"
  digest      = "1A2B3C4D5E6F7890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890"
}
