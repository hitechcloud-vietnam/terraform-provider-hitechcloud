---
page_title: "hitechcloud_account_logs Data Source - hitechcloud"
subcategory: "Account"
description: |-
  Reads the account activity log of the authenticated user.
---

# hitechcloud_account_logs (Data Source)

Reads the account activity log of the authenticated user (`GET /api/logs`).

## Example Usage

```terraform
data "hitechcloud_account_logs" "recent" {}

output "last_action" {
  value = try(data.hitechcloud_account_logs.recent.logs[0].action, null)
}
```

## Schema

### Read-Only

- `id` (String) Static identifier for this query (`account_logs`).
- `logs` (Attributes List) Account activity log entries, newest first. (see
  [nested schema](#nestedatt--logs))

<a id="nestedatt--logs"></a>

### Nested Schema for `logs`

Read-Only:

- `action` (String) Action performed.
- `created_at` (String) Timestamp of the entry.
- `id` (String) Log entry ID.
- `ip` (String) Source IP address.
- `user_agent` (String) Client user agent.
