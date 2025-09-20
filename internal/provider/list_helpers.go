package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// expandStringList converts a Terraform list attribute into a Go slice of strings.
func expandStringList(ctx context.Context, value types.List, attribute path.Path, diags *diag.Diagnostics) []string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	var result []string
	if err := value.ElementsAs(ctx, &result, false); err != nil {
		diags.AddAttributeError(
			attribute,
			"Invalid string list value",
			fmt.Sprintf("Unable to decode attribute %s into a list of strings: %v", attribute.String(), err),
		)
		return nil
	}

	return result
}

// flattenStringSlice converts a slice of strings into a Terraform List value.
func flattenStringSlice(ctx context.Context, values []string) (types.List, diag.Diagnostics) {
	if len(values) == 0 {
		return types.ListNull(types.StringType), nil
	}

	list, diags := types.ListValueFrom(ctx, types.StringType, values)
	return list, diags
}
