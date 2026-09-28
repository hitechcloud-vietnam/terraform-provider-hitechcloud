---
page_title: "hitechcloud_server_stock Data Source - hitechcloud"
subcategory: "Compute"
description: |-
  Reads the current dedicated server stock.
---

# hitechcloud_server_stock (Data Source)

Reads the current dedicated server stock (`GET /api/serverstock`).

## Example Usage

```terraform
data "hitechcloud_server_stock" "current" {}

output "stock_rows" {
  value = data.hitechcloud_server_stock.current.rows
}
```

## Schema

### Read-Only

- `id` (String) Static identifier for this query (`server_stock`).
- `rows` (Attributes List) Stock rows reported by the API. (see
  [nested schema](#nestedatt--rows))

<a id="nestedatt--rows"></a>

### Nested Schema for `rows`

Read-Only:

- `id` (String) Row identifier.
- `name` (String) Server or plan name.
- `status` (String) Availability status of the row.
- `extra` (Map of String) Remaining fields of the row.
