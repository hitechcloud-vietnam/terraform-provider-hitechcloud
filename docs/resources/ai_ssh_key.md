---
page_title: "hitechcloud_ai_ssh_key Resource - hitechcloud"
subcategory: "AI Factory"
description: |-
  Manages an SSH key of the HiTechCloud AI Factory.
---

# hitechcloud_ai_ssh_key (Resource)

Manages an SSH key of the HiTechCloud AI Factory
(`/api/service/{service_id}/ai/sshkey`). Setting `default = true` marks the key
as the account default after create/update.

## Example Usage

```terraform
resource "hitechcloud_ai_ssh_key" "main" {
  service_id = "1"
  name       = "laptop"
  public_key = file("~/.ssh/id_ed25519.pub")
  default    = true
}
```

## Schema

### Required

- `name` (String) Key name.
- `public_key` (String) OpenSSH public key material. Changing this forces a new
  resource.
- `service_id` (String) ID of the hosted service. Changing this forces a new
  resource.

### Optional

- `default` (Boolean) Make this the default SSH key of the account.

### Read-Only

- `id` (String) Composite identifier `service_id/key_id`.
- `is_default` (Boolean) Whether the key is currently the account default.
- `key_id` (String) Identifier of the key as returned by the API.

## Import

Import is supported using the composite ID `service_id/key_id`:

```shell
terraform import hitechcloud_ai_ssh_key.main 1/key-123
```
