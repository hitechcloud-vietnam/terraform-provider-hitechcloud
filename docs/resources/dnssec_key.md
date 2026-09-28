---
page_title: "hitechcloud_dnssec_key Resource - hitechcloud"
subcategory: "Domains"
description: |-
  Manages a DNSSEC key (DS record) of a registered domain.
---

# hitechcloud_dnssec_key (Resource)

Manages a DNSSEC key (DS record) of a registered domain
(`PUT /api/domain/{id}/dnssec`). Destroying the resource removes the key with
`DELETE /api/domain/{id}/dnssec/{key}`.

## Example Usage

```terraform
resource "hitechcloud_dnssec_key" "example" {
  domain_id   = "42"
  key_tag     = "2371"
  algorithm   = "13"
  digest_type = "2"
  digest      = "1A2B3C4D5E6F7890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890"
}
```

## Schema

### Required

- `algorithm` (String) DNSSEC algorithm number, for example `8` (RSA/SHA-256) or
  `13` (ECDSA P-256).
- `domain_id` (String) ID of the registered domain. Changing this forces a new
  resource.
- `key_tag` (String) DS key tag. Changing this forces a new resource.

### Optional

- `digest` (String) Hex-encoded DS digest of the DNSKEY record.
- `digest_type` (String) DS digest type, for example `2` (SHA-256).
- `flags` (String) DNSKEY flags, usually `257` (KSK) or `256` (ZSK).
- `protocol` (String) DNSKEY protocol number, usually `3`.
- `public_key` (String) Base64-encoded DNSKEY public key material.

### Read-Only

- `id` (String) Composite identifier `domain_id/key_tag`.

## Import

Import is supported using the composite ID `domain_id/key_tag`:

```shell
terraform import hitechcloud_dnssec_key.example 42/2371
```
