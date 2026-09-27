---
page_title: "Release Signing (GPG)"
subcategory: "Guides"
description: |-
  How releases of the HiTechCloud Terraform provider are GPG-signed and how to configure the signing key.
---

# Release Signing

Every release of `terraform-provider-hitechcloud` publishes zips for
`linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64` and
`windows/amd64`, plus `SHA256SUMS` and a detached GPG signature
(`SHA256SUMS.sig`). The Terraform Registry requires signed checksums when
publishing providers.

## The signing key

The public signing key is an RSA-4096 key owned by `HiTechCloud Vietnam
<release@hitechcloud.vn>` (fingerprint
`1CA3844693B18FA2FEABBC622F780BEC94499FD2`). It is **not** committed to the
repository: the ASCII-armored key (`GPG_PUBLIC_KEY.asc`) is kept in the
environment / secret store only and is git-ignored here.

To verify a release, obtain the public key from the project's release
environment (or from the GitHub account that signed it), then:

```shell
gpg --import GPG_PUBLIC_KEY.asc
gpg --verify terraform-provider-hitechcloud_v1.0.0_SHA256SUMS.sig \
             terraform-provider-hitechcloud_v1.0.0_SHA256SUMS
sha256sum -c terraform-provider-hitechcloud_v1.0.0_SHA256SUMS
```

## Configuring signing for a fork

1. Generate a key (ASCII armored output is required by GitHub and the
   Registry):

   ```shell
   gpg --full-generate-key            # RSA 4096, no expiry or long expiry
   gpg --armor --export you@example.com            > GPG_PUBLIC_KEY.asc
   gpg --armor --export-secret-keys you@example.com > GPG_PRIVATE_KEY.asc
   ```

2. Add the **public** key to GitHub: *Settings → SSH and GPG keys → New GPG
   key →* paste the content of `GPG_PUBLIC_KEY.asc` ("ASCII Armor Required").
   Register the same key in your Terraform Registry account so releases can be
   published.

3. Add the **private** key to the repository secrets:
   - `GPG_PRIVATE_KEY` — the full ASCII-armored content of
     `GPG_PRIVATE_KEY.asc`.
   - `GPG_PASSPHRASE` — the key passphrase (empty when the key has none).

4. Push a `v*` tag. `.github/workflows/release.yml` imports the key and runs
   GoReleaser with signing enabled. When the secrets are not configured, the
   workflow still produces an unsigned draft release so artifacts are never
   blocked.

~> **Security:** never commit any key material. `GPG_PRIVATE_KEY.asc`,
`GPG_PUBLIC_KEY.asc` and revocation certificates are git-ignored in this
repository; keep them outside version control and store the keys only in the
CI secret store / environment.
