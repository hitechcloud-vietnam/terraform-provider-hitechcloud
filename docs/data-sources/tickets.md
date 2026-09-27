---
page_title: "hitechcloud_tickets Data Source - hitechcloud"
subcategory: "Support"
description: |-
  Lists the support tickets of the HiTechCloud account.
---

# hitechcloud_tickets (Data Source)

Lists the support tickets of the HiTechCloud account (`GET /api/tickets`).

## Example Usage

```terraform
data "hitechcloud_tickets" "all" {}

output "open_tickets" {
  value = [for t in data.hitechcloud_tickets.all.tickets : t.number if t.status != "closed"]
}
```

## Schema

### Read-Only

- `id` (String) Data source identifier (`tickets`).
- `tickets` (Attributes List) Support tickets. (see
  [nested schema](#nestedatt--tickets))

<a id="nestedatt--tickets"></a>

### Nested Schema for `tickets`

Read-Only:

- `department` (String) Support department.
- `last_updated` (String) Last update timestamp.
- `number` (String) Ticket number.
- `status` (String) Ticket status.
- `subject` (String) Ticket subject.
