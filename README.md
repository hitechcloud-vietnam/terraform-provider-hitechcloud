# Terraform Provider for HiTechCloud

This is the official [Terraform](https://www.terraform.io) provider for the
[HiTechCloud](https://hitechcloud.vn) platform. It is built on the
[Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
and covers the HiTechCloud User API: compute (VMs), AI Factory (GPU instances,
volumes, templates, clusters), DNS, domains, S3 storage, billing, account and
portal utilities.

- **Source address:** `hitechcloud-vietnam/hitechcloud`
- **Registry namespace:** `registry.terraform.io/hitechcloud-vietnam/hitechcloud`
- **Protocol:** Terraform Plugin Protocol 6

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25 (for building from source)

## Authentication

The provider authenticates against the HiTechCloud User API
(`POST /api/login?username=...&password=...`), which returns **two tokens**: an
access `token` and a `refresh_token`. All subsequent API calls use
`Authorization: Bearer <token>`; when the access token expires the provider
renews it automatically with `POST /api/token?refresh_token=...` and replays
the failed request. `POST /api/revoke` invalidates the refresh token.

You can either log in with `username`/`password` (or `HITECHCLOUD_USERNAME` /
`HITECHCLOUD_PASSWORD`) or supply a pre-issued token via the `token` attribute
(or `HITECHCLOUD_TOKEN`), optionally with `refresh_token`
(or `HITECHCLOUD_REFRESH_TOKEN`) for automatic renewal. Never hard-code
credentials in Terraform configuration; use environment variables or a secrets
manager.

```sh
export HITECHCLOUD_USERNAME="your-username"
export HITECHCLOUD_PASSWORD="your-password"
# Or use a pre-issued token instead:
# export HITECHCLOUD_TOKEN="your-api-token"
# export HITECHCLOUD_REFRESH_TOKEN="your-refresh-token"
# Optional: point at a non-production API endpoint
export HITECHCLOUD_ENDPOINT="https://api.hitechcloud.vn"
```

## Provider Configuration

```terraform
terraform {
  required_providers {
    hitechcloud = {
      source  = "hitechcloud-vietnam/hitechcloud"
      version = "~> 1.0"
    }
  }
}

provider "hitechcloud" {
  # username        = "..."            # or HITECHCLOUD_USERNAME (2-token login)
  # password        = "..."            # or HITECHCLOUD_PASSWORD (2-token login)
  # token           = "..."            # or HITECHCLOUD_TOKEN (pre-issued)
  # refresh_token   = "..."            # or HITECHCLOUD_REFRESH_TOKEN (auto-renew)
  # endpoint        = "https://api.hitechcloud.vn" # or HITECHCLOUD_ENDPOINT
  # request_timeout = "60s"            # or HITECHCLOUD_REQUEST_TIMEOUT
}
```

## Usage Example

```terraform
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
```

## Resources

| Resource | Description |
| --- | --- |
| `hitechcloud_dns_zone` | DNS zone of a hosted service |
| `hitechcloud_dns_record` | Record inside a DNS zone |
| `hitechcloud_domain_dns_record` | DNS record of a registered domain |
| `hitechcloud_domain_settings` | Domain nameservers / autorenew / registrar lock / ID protection |
| `hitechcloud_vm` | Virtual machine |
| `hitechcloud_vm_interface` | Network interface of a VM |
| `hitechcloud_vm_firewall_rule` | VM firewall rule |
| `hitechcloud_service_ip` | Additional IP of a service |
| `hitechcloud_ai_instance` | AI Factory GPU instance |
| `hitechcloud_ai_ssh_key` | AI Factory SSH key |
| `hitechcloud_ai_volume` | AI Factory volume |
| `hitechcloud_ai_template` | AI Factory template |
| `hitechcloud_ai_cluster` | AI Factory cluster |
| `hitechcloud_s3_bucket` | S3 bucket of a storage service |
| `hitechcloud_s3_subuser` | S3 sub-user (access key / secret key) |
| `hitechcloud_url_shortener_link` | Shortened URL (portal) |
| `hitechcloud_contact` | Account contact (create/update; API has no delete) |
| `hitechcloud_ticket` | Support ticket (destroy closes it) |
| `hitechcloud_rdns` | Reverse DNS (PTR) of a service IP |

## Data Sources

| Data Source | Description |
| --- | --- |
| `hitechcloud_account` | Account / billing profile |
| `hitechcloud_contacts` | Account contacts |
| `hitechcloud_services` | All services of the account |
| `hitechcloud_service` | One service by id |
| `hitechcloud_invoices` | Invoices |
| `hitechcloud_payment_methods` | Available payment methods |
| `hitechcloud_products` | Products of a category |
| `hitechcloud_certificates` | SSL certificates |
| `hitechcloud_domains` | Registered domains |
| `hitechcloud_domain` | One domain by id or name |
| `hitechcloud_domain_tlds` | Available TLDs |
| `hitechcloud_dns_zones` | DNS zones of a service (with records) |
| `hitechcloud_vms` | VMs of a service |
| `hitechcloud_ai_instance_types` | AI instance types |
| `hitechcloud_ai_ssh_keys` | AI SSH keys |
| `hitechcloud_ai_volumes` | AI volumes |
| `hitechcloud_ai_cluster_types` | AI cluster types |
| `hitechcloud_s3_buckets` | S3 buckets of a service |
| `hitechcloud_s3_subusers` | S3 sub-users of a service |
| `hitechcloud_balance` | Account balance |
| `hitechcloud_categories` | Product categories |
| `hitechcloud_payment_fees` | Payment method fees |
| `hitechcloud_whois` | WHOIS lookup for a domain |
| `hitechcloud_domain_dns_types` | DNS record types supported by a domain |
| `hitechcloud_domain_availability` | Domain availability check |
| `hitechcloud_url_shortener_links` | Shortened URLs |
| `hitechcloud_tickets` | Support tickets |
| `hitechcloud_ticket_departments` | Support departments |
| `hitechcloud_notifications` | Portal notifications |
| `hitechcloud_statuses` | Service status entries |

## Importing Existing Resources

Resources support `terraform import` with composite IDs:

```sh
terraform import hitechcloud_dns_zone.example service_id/zone_id
terraform import hitechcloud_dns_record.www service_id/zone_id/record_id
terraform import hitechcloud_vm.web service_id/vm_id
terraform import hitechcloud_ai_instance.gpu service_id/instance_id
terraform import hitechcloud_s3_bucket.photos service_id/bucket_name
```

See the resource documentation under [`docs/`](docs/) for the exact import ID
format of every resource.

## Building from Source

```sh
git clone https://github.com/hitechcloud-vietnam/terraform-provider-hitechcloud.git
cd terraform-provider-hitechcloud
make build        # build for the current platform
make install      # install into ~/.terraform.d/plugins for local use
```

## Testing

```sh
make test      # unit tests (client mock servers, provider schema)
make testacc   # acceptance tests (TF_ACC=1, uses mock API servers)
```

The client tests run against `httptest` mock servers — no real API access is
required. Acceptance tests are gated behind `TF_ACC` and also use mock API
servers, so they can run safely in CI.

## Releasing

Releases are built with [GoReleaser](https://goreleaser.com) for
`linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64` and
`windows/amd64`, packaged as zips with SHA256SUMS and a GPG signature
(`terraform-registry-manifest.json`, protocol `6.0`). Pushing a `v*` tag
triggers `.github/workflows/release.yml`.

## Release Signing

Every release ships `SHA256SUMS` plus a detached GPG signature
(`SHA256SUMS.sig`). The ASCII-armored public signing key (fingerprint
`1CA3844693B18FA2FEABBC622F780BEC94499FD2`) is **not** committed to this
repository — it is kept in the environment / secret store only. See
[`docs/guides/release-signing.md`](docs/guides/release-signing.md) for how to
verify releases and how to configure the `GPG_PRIVATE_KEY` / `GPG_PASSPHRASE`
repository secrets for your own signing key.

## License

Mozilla Public License 2.0. See [LICENSE](LICENSE).
