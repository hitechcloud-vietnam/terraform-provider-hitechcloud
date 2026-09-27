# Provider configuration examples for the HiTechCloud Terraform provider.
# The API token is taken from the HITECHCLOUD_TOKEN environment variable.

terraform {
  required_providers {
    hitechcloud = {
      source  = "hitechcloud-vietnam/hitechcloud"
      version = "~> 0.1"
    }
  }
}

provider "hitechcloud" {
  # token           = "..."                     # or HITECHCLOUD_TOKEN
  # endpoint        = "https://api.hitechcloud.vn" # or HITECHCLOUD_ENDPOINT
  # request_timeout = "60s"                     # or HITECHCLOUD_REQUEST_TIMEOUT
}
