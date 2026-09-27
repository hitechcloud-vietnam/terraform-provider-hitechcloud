---
page_title: "Release Notes"
subcategory: "Guides"
description: |-
  Release notes for the HiTechCloud Terraform provider.
---

# Release Notes

## 1.0.0

First stable release of the HiTechCloud Terraform provider, covering the
HiTechCloud User API (`https://api.hitechcloud.vn`):

- **19 resources**: DNS zones/records, domain DNS records and settings, VMs,
  VM interfaces, firewall rules, service IPs and rDNS, AI Factory instances,
  SSH keys, volumes, templates and clusters, S3 buckets and sub-users, account
  contacts, support tickets and URL shortener links.
- **30 data sources** across account, billing, domains, DNS, compute, AI
  Factory, storage, portal and support.
- Bearer-token authentication via `token` or `HITECHCLOUD_TOKEN`, configurable
  `endpoint` and `request_timeout`.
- Releases are built for `linux/amd64`, `linux/arm64`, `darwin/amd64`,
  `darwin/arm64` and `windows/amd64` with GPG-signed checksums (protocol 6.0).

See [CHANGELOG.md](https://github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/blob/main/CHANGELOG.md)
for the full list.
