---
page_title: "hitechcloud_s3_subusers Data Source - hitechcloud"
subcategory: "Storage"
description: |-
  Lists the S3 sub-users of a HiTechCloud storage service.
---

# hitechcloud_s3_subusers (Data Source)

Lists the S3 sub-users of a HiTechCloud storage service
(`GET /api/service/{service_id}/s3/subuser`).

## Example Usage

```terraform
data "hitechcloud_s3_subusers" "all" {
  service_id = "1"
}

output "subuser_names" {
  value = [for s in data.hitechcloud_s3_subusers.all.subusers : s.name]
}
```

## Schema

### Required

- `service_id` (String) ID of the storage service.

### Read-Only

- `id` (String) Data source identifier (same as `service_id`).
- `subusers` (Attributes List) S3 sub-users of the service. (see
  [nested schema](#nestedatt--subusers))

<a id="nestedatt--subusers"></a>

### Nested Schema for `subusers`

Read-Only:

- `access` (String) Access level.
- `access_key` (String) Access key (S3 access key ID).
- `name` (String) Sub-user name.
- `uid` (String) Sub-user identifier.
