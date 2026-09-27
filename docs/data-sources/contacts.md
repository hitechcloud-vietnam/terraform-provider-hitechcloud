---
page_title: "hitechcloud_contacts Data Source - hitechcloud"
subcategory: "Account"
description: |-
  Lists the contacts of the HiTechCloud account.
---

# hitechcloud_contacts (Data Source)

Lists the contacts of the HiTechCloud account (`GET /api/contact`).

## Example Usage

```terraform
data "hitechcloud_contacts" "all" {}

output "contact_names" {
  value = [for c in data.hitechcloud_contacts.all.contacts : c.first_name]
}
```

## Schema

### Read-Only

- `contacts` (Attributes List) Contacts of the account. (see
  [nested schema](#nestedatt--contacts))
- `id` (String) Data source identifier (`contacts`).

<a id="nestedatt--contacts"></a>

### Nested Schema for `contacts`

Read-Only:

- `company_name` (String) Company name.
- `email` (String) Email address.
- `first_name` (String) First name.
- `id` (String) Contact identifier.
- `last_name` (String) Last name.
- `phone_number` (String) Phone number.
