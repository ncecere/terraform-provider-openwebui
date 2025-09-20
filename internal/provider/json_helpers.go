package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// decodeOptionalJSON maps a Terraform string attribute containing JSON into a Go map.
func decodeOptionalJSON(value types.String, attribute path.Path, diags *diag.Diagnostics) map[string]any {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	raw := strings.TrimSpace(value.ValueString())
	if raw == "" {
		return nil
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		diags.AddAttributeError(
			attribute,
			"Invalid JSON value",
			fmt.Sprintf("Expected attribute %s to contain valid JSON object text: %v", attribute.String(), err),
		)
		return nil
	}

	return result
}

// encodeOptionalJSON converts a map into a Terraform string attribute containing canonical JSON.
func encodeOptionalJSON(data map[string]any) (types.String, error) {
	if data == nil {
		return types.StringNull(), nil
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return types.StringNull(), fmt.Errorf("marshal JSON: %w", err)
	}

	return types.StringValue(string(encoded)), nil
}
