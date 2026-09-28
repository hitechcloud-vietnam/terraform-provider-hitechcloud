---
page_title: "hitechcloud_pmg_domain Resource - hitechcloud"
subcategory: "Network"
description: |-
  Protects a mail domain through a HiTechCloud PMG service and sets its mail transport.
---

# hitechcloud_pmg_domain (Resource)

Protects a mail domain through a HiTechCloud PMG (Proxmox Mail Gateway) service
and sets its target mail transport
(`POST /api/service/{id}/htcpmg/domains`,
`POST /api/service/{id}/htcpmg/transport`).

The API has no domain removal endpoint: destroying the resource removes it from
Terraform state only.

## Example Usage

```terraform
resource "hitechcloud_pmg_domain" "example" {
  service_id = "1"
  domain     = "example.com"
  host       = "mail.example.com"
  port       = "25"
}
```

## Schema

### Required

- `domain` (String) Mail domain to protect. Changing this forces a new resource.
- `service_id` (String) HiTechCloud service ID of the PMG service. Changing this
  forces a new resource.

### Optional

- `host` (String) Target mail server (hostname or IP) receiving filtered mail.
- `port` (String) Target SMTP port, defaults to `25`.

### Read-Only

- `id` (String) Composite identifier `service_id/domain`.

## Import

Import is supported using the composite ID `service_id/domain`:

```shell
terraform import hitechcloud_pmg_domain.example 1/example.com
```
