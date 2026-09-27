---
page_title: "hitechcloud_contact Resource - hitechcloud"
subcategory: "Account"
description: |-
  Manages a contact of the HiTechCloud account.
---

# hitechcloud_contact (Resource)

Manages a contact of the HiTechCloud account (`POST/PUT /api/contact`). The API
has no contact deletion endpoint: destroying the resource only removes it from
state and emits a warning — the contact remains in the account.

## Example Usage

```terraform
resource "hitechcloud_contact" "ops" {
  email       = "ops@example.com"
  password    = var.contact_password
  first_name  = "Operations"
  last_name   = "Team"
  company_name = "Example Corp"
  phone_number = "+84-24-0000-0000"
}
```

## Schema

### Required

- `email` (String) Email address of the contact.

### Optional

- `address` (String) Street address.
- `bank_account` (String) Bank account number.
- `bank_account_name` (String) Bank account holder name.
- `bank_name` (String) Bank name.
- `birthday` (String) Birthday.
- `city` (String) City.
- `company_name` (String) Company name.
- `country` (String) Country.
- `first_name` (String) First name.
- `gender` (String) Gender.
- `last_name` (String) Last name.
- `national_id` (String) National ID.
- `password` (String, Sensitive) Password for the contact login. Only used at
  creation; changing this forces a new resource.
- `phone_number` (String) Phone number.
- `postcode` (String) Postal code.
- `privileges` (String) Privilege set of the contact.
- `state` (String) State / province.
- `tax_id` (String) Tax identifier.
- `type` (String) Contact type.

### Read-Only

- `contact_id` (String) Identifier of the contact as returned by the API.
- `id` (String) Composite identifier (same as `contact_id`).

## Import

Import is supported using the contact ID:

```shell
terraform import hitechcloud_contact.ops 7
```
