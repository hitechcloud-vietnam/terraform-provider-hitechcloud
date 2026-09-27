---
page_title: "hitechcloud_s3_buckets Data Source - hitechcloud"
subcategory: "Storage"
description: |-
  Lists the S3 buckets of a HiTechCloud storage service.
---

# hitechcloud_s3_buckets (Data Source)

Lists the S3 buckets of a HiTechCloud storage service
(`GET /api/service/{service_id}/s3/bucket`).

## Example Usage

```terraform
data "hitechcloud_s3_buckets" "all" {
  service_id = "1"
}

output "bucket_names" {
  value = [for b in data.hitechcloud_s3_buckets.all.buckets : b.name]
}
```

## Schema

### Required

- `service_id` (String) ID of the storage service.

### Read-Only

- `buckets` (Attributes List) S3 buckets of the service. (see
  [nested schema](#nestedatt--buckets))
- `id` (String) Data source identifier (same as `service_id`).

<a id="nestedatt--buckets"></a>

### Nested Schema for `buckets`

Read-Only:

- `name` (String) Bucket name.
- `objects` (Number) Number of objects in the bucket.
- `size_bytes` (Number) Total size in bytes.
