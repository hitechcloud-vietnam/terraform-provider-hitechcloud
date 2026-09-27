---
page_title: "hitechcloud_ai_ssh_keys Data Source - hitechcloud"
subcategory: "AI Factory"
description: |-
  Lists the SSH keys of the HiTechCloud AI Factory.
---

# hitechcloud_ai_ssh_keys (Data Source)

Lists the SSH keys of the HiTechCloud AI Factory
(`GET /api/service/{service_id}/sshkeys`).

## Example Usage

```terraform
data "hitechcloud_ai_ssh_keys" "all" {
  service_id = "1"
}

output "key_ids" {
  value = [for k in data.hitechcloud_ai_ssh_keys.all.keys : k.key_id]
}
```

## Schema

### Required

- `service_id` (String) ID of the hosted service.

### Read-Only

- `id` (String) Data source identifier (same as `service_id`).
- `keys` (Attributes List) SSH keys. (see
  [nested schema](#nestedatt--keys))

<a id="nestedatt--keys"></a>

### Nested Schema for `keys`

Read-Only:

- `fingerprint` (String) Key fingerprint.
- `is_default` (Boolean) Whether the key is the account default.
- `key_id` (String) Key identifier.
- `name` (String) Key name.
- `public_key` (String) Public key material.
