---
page_title: "hitechcloud_ticket Resource - hitechcloud"
subcategory: "Support"
description: |-
  Manages a HiTechCloud support ticket.
---

# hitechcloud_ticket (Resource)

Manages a HiTechCloud support ticket (`POST /api/tickets`). Destroying the
resource closes the ticket (`PUT /api/tickets/{number}/close`); the API has no
ticket deletion. Replies are not managed by Terraform — the `body` is only sent
when the ticket is opened.

## Example Usage

```terraform
data "hitechcloud_ticket_departments" "all" {}

resource "hitechcloud_ticket" "request" {
  dept_id = data.hitechcloud_ticket_departments.all.departments[0].id
  subject = "Request: additional IP allocation"
  body    = "Please allocate an additional IP for service 1."
}
```

## Schema

### Required

- `body` (String) Initial ticket body. Changing this forces a new resource.
- `dept_id` (String) Support department identifier. Changing this forces a new
  resource.
- `subject` (String) Ticket subject. Changing this forces a new resource.

### Read-Only

- `department` (String) Support department name.
- `id` (String) Ticket number.
- `last_updated` (String) Last update timestamp.
- `number` (String) Ticket number as returned by the API.
- `status` (String) Ticket status.

## Import

Import is supported using the ticket number:

```shell
terraform import hitechcloud_ticket.request 42
```
