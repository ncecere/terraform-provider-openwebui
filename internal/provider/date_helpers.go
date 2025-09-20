package provider

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func formatDateValue(ts int64) types.String {
	if ts <= 0 {
		return types.StringNull()
	}

	return types.StringValue(time.Unix(ts, 0).UTC().Format("2006-01-02"))
}
