---
page_title: "hitechcloud_s3_bucket Resource - hitechcloud"
subcategory: "Storage"
description: |-
  Manages an S3 bucket of a HiTechCloud storage service.
---

# hitechcloud_s3_bucket (Resource)

Manages an S3 bucket of a HiTechCloud storage service
(`/api/service/{service_id}/s3/bucket`). The API may auto-prefix the requested
name; the effective name is exposed as the computed `bucket` attribute.

## Example Usage

```terraform
resource "hitechcloud_s3_bucket" "photos" {
  service_id = "1"
  name       = "photos"

  # Remove all objects before destroying the bucket.
  purge = true
}
```

## Schema

### Required

- `name` (String) Requested bucket name. Changing this forces a new resource.
- `service_id` (String) ID of the storage service. Changing this forces a new
  resource.

### Optional

- `purge` (Boolean) Delete all objects before destroying the bucket. Defaults
  to `false`.

### Read-Only

- `bucket` (String) Effective bucket name assigned by the API.
- `id` (String) Composite identifier `service_id/bucket_name`.

## Import

Import is supported using the composite ID `service_id/bucket_name`:

```shell
terraform import hitechcloud_s3_bucket.photos 1/photos
```
