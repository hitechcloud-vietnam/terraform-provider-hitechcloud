---
page_title: "hitechcloud_dnssec_keys Data Source - hitechcloud"
subcategory: "Domains"
description: |-
  Lists the DNSSEC keys and available flags of a registered domain.
---

# hitechcloud_dnssec_keys (Data Source)

Lists the DNSSEC keys (DS records) and available flags of a registered domain
(`GET /api/domain/{id}/dnssec`, `GET /api/domain/{id}/dnssec/flags`).

## Example Usage

```terraform
data "hitechcloud_dnssec_keys" "example" {
  domain_id = "42"
}

output "key_tags" {
  value = [for k in data.hitechcloud_dnssec_keys.example.keys : k.key_tag]
}
```

## Schema

### Required

- `domain_id` (String) ID of the registered domain.

### Read-Only

- `available_flags` (List of String) Flags the API accepts for DNSSEC keys.
- `id` (String) The domain ID (same as `domain_id`).
- `keys` (Attributes List) DNSSEC keys configured for the domain. (see
  [nested schema](#nestedatt--keys))

<a id="nestedatt--keys"></a>

### Nested Schema for `keys`

Read-Only:

- `algorithm` (String) DNSSEC algorithm number.
- `digest` (String) Hex-encoded DS digest.
- `digest_type` (String) DS digest type.
- `flags` (String) DNSKEY flags.
- `key_tag` (String) DS key tag.
- `protocol` (String) DNSKEY protocol number.
- `public_key` (String) Base64-encoded public key material.
