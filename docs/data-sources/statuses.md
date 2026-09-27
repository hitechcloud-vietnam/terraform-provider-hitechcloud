---
page_title: "hitechcloud_statuses Data Source - hitechcloud"
subcategory: "Support"
description: |-
  Lists the service status entries of the HiTechCloud account.
---

# hitechcloud_statuses (Data Source)

Lists the service status entries of the HiTechCloud account
(`GET /api/statuses`).

## Example Usage

```terraform
data "hitechcloud_statuses" "all" {}

output "service_statuses" {
  value = { for s in data.hitechcloud_statuses.all.statuses : s.name => s.status }
}
```

## Schema

### Read-Only

- `id` (String) Data source identifier (`statuses`).
- `statuses` (Attributes List) Service status entries. (see
  [nested schema](#nestedatt--statuses))

<a id="nestedatt--statuses"></a>

### Nested Schema for `statuses`

Read-Only:

- `id` (String) Service identifier.
- `name` (String) Service name.
- `status` (String) Current status.
- `type` (String) Service type / group.
