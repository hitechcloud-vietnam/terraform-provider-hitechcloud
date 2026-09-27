---
page_title: "hitechcloud_ticket_departments Data Source - hitechcloud"
subcategory: "Support"
description: |-
  Lists the support departments of the HiTechCloud account.
---

# hitechcloud_ticket_departments (Data Source)

Lists the support departments of the HiTechCloud account
(`GET /api/ticket/departments`).

## Example Usage

```terraform
data "hitechcloud_ticket_departments" "all" {}

output "department_names" {
  value = [for d in data.hitechcloud_ticket_departments.all.departments : d.name]
}
```

## Schema

### Read-Only

- `departments` (Attributes List) Support departments. (see
  [nested schema](#nestedatt--departments))
- `id` (String) Data source identifier (`ticket_departments`).

<a id="nestedatt--departments"></a>

### Nested Schema for `departments`

Read-Only:

- `id` (String) Department identifier.
- `name` (String) Department name.
