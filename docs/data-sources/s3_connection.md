---
page_title: "hitechcloud_s3_connection Data Source - hitechcloud"
subcategory: "Storage"
description: |-
  Reads an S3 object-storage service: connection info, credentials and usage metrics.
---

# hitechcloud_s3_connection (Data Source)

Reads an S3 object-storage service: connection info, credentials and usage
metrics (`GET /api/service/{id}/s3`, `/s3/credentials`, `/s3/usage`,
`/s3/metrics`).

## Example Usage

```terraform
data "hitechcloud_s3_connection" "primary" {
  service_id = "1"
}

output "s3_endpoint" {
  value = data.hitechcloud_s3_connection.primary.endpoint
}
```

## Schema

### Required

- `service_id` (String) HiTechCloud service ID (`hb_accounts.id`) of the S3
  service.

### Read-Only

- `access_key` (String, Sensitive) S3 access key.
- `endpoint` (String) S3 endpoint URL.
- `id` (String) The service ID (same as `service_id`).
- `metrics` (Map of String) Raw metric values keyed by metric name.
- `region` (String) S3 region.
- `secret_key` (String, Sensitive) S3 secret key.
- `usage` (Map of String) Billing-relevant usage counters with units.
