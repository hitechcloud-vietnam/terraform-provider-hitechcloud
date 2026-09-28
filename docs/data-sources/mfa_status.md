---
page_title: "hitechcloud_mfa_status Data Source - hitechcloud"
subcategory: "Account"
description: |-
  Reads the MFA state of a user: passkey status and credentials, email MFA status and codes.
---

# hitechcloud_mfa_status (Data Source)

Reads the MFA state of a user: passkey status and credentials, email MFA status
and codes
(`GET /api/passkeyv2/status/{user_type}/{user_id}`,
`/api/passkeyv2/credentials/{user_type}/{user_id}`,
`/api/email_mfa_v2/status/{user_type}/{user_id}`,
`/api/email_mfa_v2/list/{user_type}/{user_id}`).

## Example Usage

```terraform
data "hitechcloud_mfa_status" "me" {
  user_type = "client"
  user_id   = "100"
}
```

## Schema

### Required

- `user_id` (String) User identifier.
- `user_type` (String) User type, for example `client`.

### Read-Only

- `email_mfa_codes` (Attributes List) Email MFA codes / purposes on file. (see
  [nested schema](#nestedatt--email_mfa_codes))
- `email_mfa_status` (Map of String) Email MFA status fields.
- `id` (String) Composite identifier `user_type/user_id`.
- `passkey_credentials` (Attributes List) Registered passkey credentials. (see
  [nested schema](#nestedatt--passkey_credentials))
- `passkey_status` (Map of String) Passkey MFA status fields.

<a id="nestedatt--passkey_credentials"></a>

### Nested Schema for `passkey_credentials`

Read-Only:

- `id` (String) Credential identifier.
- `name` (String) Credential label.
- `created_at` (String) Registration timestamp.
- `last_used_at` (String) Last usage timestamp.

<a id="nestedatt--email_mfa_codes"></a>

### Nested Schema for `email_mfa_codes`

Read-Only:

- `id` (String) Code entry identifier.
- `purpose` (String) Purpose of the code.
- `status` (String) Code status.
- `created_at` (String) Creation timestamp.
