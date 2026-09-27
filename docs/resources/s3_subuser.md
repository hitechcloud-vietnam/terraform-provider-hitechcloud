---
page_title: "hitechcloud_s3_subuser Resource - hitechcloud"
subcategory: "Storage"
description: |-
  Manages an S3 sub-user (access key) of a HiTechCloud storage service.
---

# hitechcloud_s3_subuser (Resource)

Manages an S3 sub-user of a HiTechCloud storage service
(`/api/service/{service_id}/s3/subusers`). The `secret_key` is returned by the
API only once at creation and is stored in state; it cannot be recovered later.

## Example Usage

```terraform
resource "hitechcloud_s3_subuser" "backup" {
  service_id = "1"
  name       = "backup-agent"
  access     = "readwrite"
}

output "s3_secret" {
  value     = hitechcloud_s3_subuser.backup.secret_key
  sensitive = true
}
```

## Schema

### Required

- `name` (String) Sub-user name. Changing this forces a new resource.
- `service_id` (String) ID of the storage service. Changing this forces a new
  resource.

### Optional

- `access` (String) Access level (`read`, `write`, `readwrite`, `full`).
  Defaults to the API default. Changing this forces a new resource.

### Read-Only

- `access_key` (String) Access key (S3 access key ID).
- `id` (String) Composite identifier `service_id/name`.
- `secret_key` (String, Sensitive) Secret key. Only available immediately
  after creation.

## Import

Import is supported using the composite ID `service_id/name`. Note that the
`secret_key` cannot be imported and will be empty for imported resources:

```shell
terraform import hitechcloud_s3_subuser.backup 1/backup-agent
```
