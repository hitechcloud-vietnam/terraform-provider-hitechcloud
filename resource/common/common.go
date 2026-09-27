// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

// Package common contains helpers shared by the Terraform resources and data
// sources under /resource.
package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SetStrings converts a []string into a types.Set of strings.
func SetStrings(values []string) types.Set {
	if len(values) == 0 {
		return types.SetNull(types.StringType)
	}
	set, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		// Cannot happen for []string -> Set[string]; keep the empty fallback.
		return types.SetNull(types.StringType)
	}
	return set
}

// StringsFromSet converts a types.Set of strings into a []string.
func StringsFromSet(ctx context.Context, s types.Set) ([]string, diag.Diagnostics) {
	var out []string
	if s.IsNull() || s.IsUnknown() {
		return out, nil
	}
	diags := s.ElementsAs(ctx, &out, false)
	return out, diags
}

// MapStrings converts a map[string]string into a types.Map of strings.
func MapStrings(values map[string]string) types.Map {
	if len(values) == 0 {
		return types.MapNull(types.StringType)
	}
	m, diags := types.MapValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		return types.MapNull(types.StringType)
	}
	return m
}

// StringsFromMap converts a types.Map of strings into map[string]string.
func StringsFromMap(ctx context.Context, m types.Map) (map[string]string, diag.Diagnostics) {
	out := map[string]string{}
	if m.IsNull() || m.IsUnknown() {
		return out, nil
	}
	diags := m.ElementsAs(ctx, &out, false)
	return out, diags
}

// SetStringsOrEmpty converts values into a set, preserving an explicitly
// configured empty set (avoiding null-vs-empty inconsistency errors).
func SetStringsOrEmpty(values []string, configured types.Set) types.Set {
	if len(values) > 0 {
		return SetStrings(values)
	}
	if !configured.IsNull() && !configured.IsUnknown() {
		return types.SetValueMust(types.StringType, []attr.Value{})
	}
	return types.SetNull(types.StringType)
}

// MapStringsOrEmpty is the map counterpart of SetStringsOrEmpty.
func MapStringsOrEmpty(values map[string]string, configured types.Map) types.Map {
	if len(values) > 0 {
		return MapStrings(values)
	}
	if !configured.IsNull() && !configured.IsUnknown() {
		return types.MapValueMust(types.StringType, map[string]attr.Value{})
	}
	return types.MapNull(types.StringType)
}

// StringOrNull returns a typed string, normalising empty values to null so the
// API-facing optional attributes behave predictably.
func StringOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// Int64OrNull normalises zero to null for optional integer attributes.
func Int64OrNull(n int64) types.Int64 {
	if n == 0 {
		return types.Int64Null()
	}
	return types.Int64Value(n)
}

// BoolPtr is a small helper for optional client request flags.
func BoolPtr(b bool) *bool { return &b }

// FormatIDError returns a consistent error message for malformed composite IDs.
func FormatIDError(id string, want int) error {
	return fmt.Errorf("unexpected ID %q: want %d slash-separated parts", id, want)
}
