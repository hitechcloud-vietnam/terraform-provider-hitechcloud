---
page_title: "hitechcloud_willexpired Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Lists services and domains approaching their expiry date.
---

# hitechcloud_willexpired (Data Source)

Lists services and domains approaching their expiry date
(`GET /api/willexpired`, `/willexpired/summary`, `/willexpired/invoices`).

## Example Usage

```terraform
data "hitechcloud_willexpired" "soon" {
  item_type = "domain"
}

output "expiring_domains" {
  value = [for item in data.hitechcloud_willexpired.soon.items : item.name]
}
```

## Schema

### Optional

- `item_type` (String) Filter by item type: `service` or `domain`. Empty lists
  both.
- `status` (String) Filter by expiry status.

### Read-Only

- `id` (String) Static identifier for this query.
- `items` (Attributes List) Items approaching expiry. (see
  [nested schema](#nestedatt--items))
- `summary` (Map of String) Aggregate counters of the expiring items.
- `upcoming_invoices` (Attributes List) Invoices generated for expiring items.
  (see [nested schema](#nestedatt--upcoming_invoices))

<a id="nestedatt--items"></a>

### Nested Schema for `items`

Read-Only:

- `autorenew` (Bool) Whether auto-renew is enabled.
- `expires_at` (String) Expiry timestamp.
- `id` (String) Item ID (service ID or domain name).
- `name` (String) Item display name.
- `status` (String) Expiry status.
- `type` (String) Item type (service/domain).

<a id="nestedatt--upcoming_invoices"></a>

### Nested Schema for `upcoming_invoices`

Read-Only:

- `amount` (String) Invoice amount.
- `due_at` (String) Invoice due date.
- `id` (String) Invoice ID.
- `item_id` (String) Related item ID.
