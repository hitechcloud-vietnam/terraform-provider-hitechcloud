---
page_title: "hitechcloud_partner_payout Resource - hitechcloud"
subcategory: "Portal"
description: |-
  Requests a partner wallet payout.
---

# hitechcloud_partner_payout (Resource)

Requests a partner wallet payout (`POST /api/partner/payouts`). Payout requests
are immutable: changing any attribute forces a replacement.

The API has no payout cancellation endpoint — deleting the resource removes it
from Terraform state only and emits a warning.

## Example Usage

```terraform
resource "hitechcloud_partner_payout" "monthly" {
  amount = "5000000"
  method = "bank_transfer"
  note   = "September commission"
}
```

## Schema

### Required

- `amount` (String) Amount to pay out. Changing this forces a new resource.
- `method` (String) Payout method (e.g. bank transfer). Changing this forces a new
  resource.

### Optional

- `note` (String) Optional note for the payout request. Changing this forces a new
  resource.

### Read-Only

- `created_at` (String) Request timestamp reported by the API.
- `id` (String) Payout request identifier returned by the API.
- `status` (String) Current payout status reported by the API.

## Import

Import is supported using the payout request ID:

```shell
terraform import hitechcloud_partner_payout.monthly payout-1001
```
