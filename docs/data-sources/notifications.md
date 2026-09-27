---
page_title: "hitechcloud_notifications Data Source - hitechcloud"
subcategory: "Support"
description: |-
  Lists the portal notifications of the HiTechCloud account.
---

# hitechcloud_notifications (Data Source)

Lists the portal notifications of the HiTechCloud account
(`GET /api/notifications`), optionally filtered by related object.

## Example Usage

```terraform
data "hitechcloud_notifications" "all" {}

output "notification_titles" {
  value = [for n in data.hitechcloud_notifications.all.notifications : n.title]
}
```

## Schema

### Optional

- `rel_id` (String) Filter by related object identifier.
- `rel_type` (String) Filter by related object type.

### Read-Only

- `id` (String) Data source identifier (`notifications`).
- `notifications` (Attributes List) Portal notifications. (see
  [nested schema](#nestedatt--notifications))

<a id="nestedatt--notifications"></a>

### Nested Schema for `notifications`

Read-Only:

- `date` (String) Creation timestamp.
- `id` (String) Notification identifier.
- `message` (String) Notification message.
- `status` (String) Notification status.
- `title` (String) Notification title.
