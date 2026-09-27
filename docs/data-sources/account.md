---
page_title: "hitechcloud_account Data Source - hitechcloud"
subcategory: "Account"
description: |-
  Retrieves the HiTechCloud account profile and billing details.
---

# hitechcloud_account (Data Source)

Retrieves the HiTechCloud account profile and billing details
(`GET /api/details`).

## Example Usage

```terraform
data "hitechcloud_account" "current" {}

output "account_email" {
  value = data.hitechcloud_account.current.email
}
```

## Schema

### Read-Only

- `address` (String) Street address.
- `bank_account` (String) Bank account number.
- `bank_account_name` (String) Bank account holder name.
- `bank_name` (String) Bank name.
- `birthday` (String) Birthday.
- `city` (String) City.
- `company_name` (String) Company name.
- `country` (String) Country.
- `currency` (String) Preferred currency.
- `email` (String) Account email.
- `first_name` (String) First name.
- `gender` (String) Gender.
- `id` (String) Account identifier.
- `last_name` (String) Last name.
- `national_id` (String) National ID / tax information.
- `phone_number` (String) Phone number.
- `postcode` (String) Postal code.
- `state` (String) State / province.
- `tax_id` (String) Tax identifier.
- `type` (String) Account type.
