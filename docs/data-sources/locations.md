---
page_title: "hitechcloud_locations Data Source - hitechcloud"
subcategory: "Account"
description: |-
  Reads supported countries and states.
---

# hitechcloud_locations (Data Source)

Reads supported countries and their states
(`GET /api/countries`, `GET /api/states`). When `country_id` is set, only the
states of that country are returned; otherwise all states are listed.

## Example Usage

```terraform
data "hitechcloud_locations" "vn" {
  country_id = "VN"
}
```

## Schema

### Optional

- `country_id` (String) Restrict the state list to this country.

### Read-Only

- `countries` (Attributes List) Supported countries. (see
  [nested schema](#nestedatt--countries))
- `id` (String) Static identifier for this query (`locations`).
- `states` (Attributes List) Supported states of the requested countries. (see
  [nested schema](#nestedatt--states))

<a id="nestedatt--countries"></a>

### Nested Schema for `countries`

Read-Only:

- `id` (String) Country identifier.
- `name` (String) Country name.
- `extra` (Map of String) Remaining fields of the row.

<a id="nestedatt--states"></a>

### Nested Schema for `states`

Read-Only:

- `id` (String) State identifier.
- `name` (String) State name.
- `country_id` (String) Country the state belongs to.
- `extra` (Map of String) Remaining fields of the row.
