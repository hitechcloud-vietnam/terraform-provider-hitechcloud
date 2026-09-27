#!/usr/bin/env bash
# Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0
#
# Local release-style build: cross-compiles the provider for all supported
# platforms and packages each artifact as a zip together with the Terraform
# Registry manifest, producing SHA256SUMS at the end. Useful for manual
# releases or verifying the release layout without GoReleaser.
set -euo pipefail

PROJECT="terraform-provider-hitechcloud"
VERSION="${1:-0.0.1-dev}"
DIST="dist"

OS_ARCH=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
)

rm -rf "${DIST}"
mkdir -p "${DIST}"

for entry in "${OS_ARCH[@]}"; do
  read -r os arch <<< "${entry}"
  ext=""
  if [ "${os}" = "windows" ]; then ext=".exe"; fi
  out="${DIST}/${PROJECT}_${VERSION}_${os}_${arch}"
  echo "Building ${os}/${arch} ..."
  # The binary must be named terraform-provider-<NAME>_v<VERSION> for the
  # Terraform Registry (same convention as GoReleaser).
  GOOS="${os}" GOARCH="${arch}" CGO_ENABLED=0 \
    go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "${out}/${PROJECT}_v${VERSION}${ext}" .
  cp LICENSE README.md CHANGELOG.md terraform-registry-manifest.json "${out}/" 2>/dev/null || true
  mv "${out}/terraform-registry-manifest.json" "${out}/${PROJECT}_${VERSION}_manifest.json" 2>/dev/null || true
  (cd "${out}" && zip -q -r "../${PROJECT}_${VERSION}_${os}_${arch}.zip" .)
  rm -rf "${out}"
done

(cd "${DIST}" && sha256sum ./*.zip > "${PROJECT}_${VERSION}_SHA256SUMS")
echo "Done. Artifacts in ${DIST}/"
