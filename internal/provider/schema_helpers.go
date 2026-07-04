package provider

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const rateLimitBytesPerSecondKey = "rate_limit_bytes_per_second"

func optionalIntPointer(d *schema.ResourceData, key string) *int {
	if !configuredValueExists(d, key) {
		return nil
	}

	intValue := d.Get(key).(int)
	return &intValue
}

func configuredValueExists(d *schema.ResourceData, key string) bool {
	value, diags := d.GetRawConfigAt(cty.GetAttrPath(key))
	if diags.HasError() {
		return false
	}

	return value.IsKnown() && !value.IsNull()
}

func setOptionalInt(d *schema.ResourceData, key string, value *int) error {
	return d.Set(key, value)
}
