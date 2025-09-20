package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/nickcecere/terraform-provider-openwebui/internal/client"
)

func resolveGroupNamesToIDs(ctx context.Context, apiClient *client.Client, names []string, attribute path.Path, diags *diag.Diagnostics) []string {
	if len(names) == 0 {
		return nil
	}

	groups, err := apiClient.ListGroups(ctx)
	if err != nil {
		diags.AddAttributeError(
			attribute,
			"Unable to list groups",
			fmt.Sprintf("Failed to retrieve groups from Open WebUI: %v", err),
		)
		return nil
	}

	byName := make(map[string]string, len(groups))
	byID := make(map[string]string, len(groups))
	for _, g := range groups {
		lowerName := strings.ToLower(g.Name)
		byName[lowerName] = g.ID
		byID[strings.ToLower(g.ID)] = g.ID
	}

	var ids []string
	for _, raw := range names {
		identifier := strings.TrimSpace(raw)
		if identifier == "" {
			continue
		}

		key := strings.ToLower(identifier)
		if id, ok := byName[key]; ok {
			ids = append(ids, id)
			continue
		}
		if id, ok := byID[key]; ok {
			ids = append(ids, id)
			continue
		}

		group, err := apiClient.GetGroup(ctx, identifier)
		if err == nil {
			ids = append(ids, group.ID)
			continue
		}

		diags.AddAttributeError(
			attribute,
			"Unknown group reference",
			fmt.Sprintf("No Open WebUI group was found for %q.", identifier),
		)
	}

	return uniqueStrings(ids)
}

func fetchGroupNamesForIDs(ctx context.Context, apiClient *client.Client, ids []string) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(ids) == 0 {
		return nil, diags
	}

	groups, err := apiClient.ListGroups(ctx)
	if err != nil {
		diags.AddError(
			"Unable to list groups",
			fmt.Sprintf("Failed to retrieve groups from Open WebUI: %v", err),
		)
		return nil, diags
	}

	byID := make(map[string]string, len(groups))
	for _, g := range groups {
		byID[g.ID] = g.Name
	}

	var names []string
	for _, id := range ids {
		if name, ok := byID[id]; ok {
			names = append(names, name)
			continue
		}

		group, err := apiClient.GetGroup(ctx, id)
		if err != nil {
			if err == client.ErrNotFound {
				continue
			}
			diags.AddError(
				"Fetch group failed",
				fmt.Sprintf("Failed to retrieve group %s: %v", id, err),
			)
			continue
		}

		names = append(names, group.Name)
	}

	return uniqueStrings(names), diags
}

func buildAccessControl(readIDs, writeIDs []string) map[string]any {
	if len(readIDs) == 0 && len(writeIDs) == 0 {
		return nil
	}

	mergedRead := uniqueStrings(append(readIDs, writeIDs...))

	control := make(map[string]any)
	if len(mergedRead) > 0 {
		control["read"] = map[string]any{
			"group_ids": mergedRead,
			"user_ids":  []string{},
		}
	}

	if len(writeIDs) > 0 {
		control["write"] = map[string]any{
			"group_ids": writeIDs,
			"user_ids":  []string{},
		}
	}

	return control
}

func extractGroupIDsFromAccessControl(access map[string]any, section string) []string {
	if access == nil {
		return nil
	}

	raw, ok := access[section]
	if !ok || raw == nil {
		return nil
	}

	sectionMap, ok := raw.(map[string]any)
	if !ok {
		return nil
	}

	idsRaw, ok := sectionMap["group_ids"]
	if !ok || idsRaw == nil {
		return nil
	}

	switch v := idsRaw.(type) {
	case []any:
		var ids []string
		for _, item := range v {
			if str, ok := item.(string); ok && str != "" {
				ids = append(ids, str)
			}
		}
		return ids
	case []string:
		return v
	default:
		return nil
	}
}
