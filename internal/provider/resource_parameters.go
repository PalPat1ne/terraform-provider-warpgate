package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/warp-tech/terraform-provider-warpgate/internal/client"
)

func resourceParameters() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParametersCreate,
		ReadContext:   resourceParametersRead,
		UpdateContext: resourceParametersUpdate,
		DeleteContext: resourceParametersDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"allow_own_credential_management": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Allow users to manage their own credentials",
			},
			"rate_limit_bytes_per_second": optionalIntParameter(
				"Global bandwidth limit",
				validation.IntAtLeast(0),
			),
			"ssh_client_auth_publickey": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable SSH public key authentication",
			},
			"ssh_client_auth_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable SSH password authentication",
			},
			"ssh_client_auth_keyboard_interactive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable SSH keyboard interactive authentication",
			},
			"password_login_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "How the password login form is presented on the gateway login page.",
				ValidateFunc: validation.StringInSlice([]string{"Enabled", "Minimized", "Disabled"}, false),
			},
			"ticket_self_service_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable ticket self-service.",
			},
			"ticket_auto_approve_existing_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Automatically approve ticket requests when the requester already has access.",
			},
			"ticket_max_duration_seconds": optionalIntParameter(
				"Maximum ticket duration in seconds.",
				validation.IntAtLeast(0),
			),
			"ticket_max_uses": optionalIntParameter(
				"Maximum number of uses for tickets.",
				validation.IntBetween(0, 32767),
			),
			"ticket_require_description": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Require a description for ticket requests.",
			},
			"ticket_request_show_all_targets": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Show all targets when requesting tickets.",
			},
			"target_click_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Action to take when clicking a target.",
				ValidateFunc: validation.StringInSlice([]string{"Connect", "ShowInstructions"}, false),
			},
			"show_session_menu": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "When enabled, Warpgate injects a session menu into HTTP sessions, allowing users to log out or return to the home page.",
			},
			"password_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Password policy rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_length": optionalIntParameter(
							"Minimum number of characters, or 0 for no requirement.",
							validation.IntAtLeast(0),
						),
						"require_uppercase": optionalComputedBoolParameter("Require at least one uppercase character."),
						"require_lowercase": optionalComputedBoolParameter("Require at least one lowercase character."),
						"require_digits":    optionalComputedBoolParameter("Require at least one digit."),
						"require_special":   optionalComputedBoolParameter("Require at least one special character."),
					},
				},
			},
			"max_api_token_duration_seconds": optionalIntParameter(
				"Maximum API token duration in seconds.",
				validation.IntAtLeast(0),
			),
			"record_scp": optionalComputedBoolParameter("Record SCP sessions."),
			"login_protection_enabled": optionalComputedBoolParameter(
				"Enable login protection.",
			),
			"login_protection_retention_seconds": optionalIntParameter(
				"How long login protection records are retained, in seconds.",
				validation.IntAtLeast(0),
			),
			"lp_ip_max_attempts": optionalIntParameter(
				"Maximum failed login attempts per IP address.",
				validation.IntAtLeast(0),
			),
			"lp_ip_time_window_seconds": optionalIntParameter(
				"Time window for failed login attempts per IP address, in seconds.",
				validation.IntAtLeast(0),
			),
			"lp_ip_base_block_duration_seconds": optionalIntParameter(
				"Base IP block duration in seconds.",
				validation.IntAtLeast(0),
			),
			"lp_ip_block_duration_multiplier": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Computed:     true,
				Description:  "Multiplier applied to repeated IP block durations.",
				ValidateFunc: validation.FloatAtLeast(0),
			},
			"lp_ip_max_block_duration_seconds": optionalIntParameter(
				"Maximum IP block duration in seconds.",
				validation.IntAtLeast(0),
			),
			"lp_ip_cooldown_reset_seconds": optionalIntParameter(
				"Cooldown period before the IP block escalation resets, in seconds.",
				validation.IntAtLeast(0),
			),
			"lp_user_max_attempts": optionalIntParameter(
				"Maximum failed login attempts per user.",
				validation.IntAtLeast(0),
			),
			"lp_user_time_window_seconds": optionalIntParameter(
				"Time window for failed login attempts per user, in seconds.",
				validation.IntAtLeast(0),
			),
			"lp_user_auto_unlock": optionalComputedBoolParameter(
				"Automatically unlock users after the lockout duration.",
			),
			"lp_user_lockout_duration_seconds": optionalIntParameter(
				"User lockout duration in seconds.",
				validation.IntAtLeast(0),
			),
			"lp_user_exempt_admins": optionalComputedBoolParameter(
				"Exempt administrators from user login protection lockouts.",
			),
			"ssh_banner": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Banner shown to SSH clients before authentication.",
			},
			"web_ssh_enabled": optionalComputedBoolParameter(
				"Enable web-based SSH sessions.",
			),
			"analytics_consent": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Whether the instance reports anonymous usage analytics.",
				ValidateFunc: validation.StringInSlice([]string{"Undecided", "Off", "On"}, false),
			},
			"analytics_normal": optionalComputedBoolParameter(
				"Enable the normal analytics payload level.",
			),
		},
	}
}

func optionalComputedBoolParameter(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: description,
	}
}

func optionalIntParameter(description string, validateFunc schema.SchemaValidateFunc) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		Description:  description,
		ValidateFunc: validateFunc,
	}
}

// resourceParametersCreate handles the creation of Warpgate parameters (singleton resource)
func resourceParametersCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerMeta := meta.(*providerMeta)
	c := providerMeta.client

	req := expandParametersUpdateRequest(d)

	_, err := c.UpdateParameters(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create parameters: %w", err))
	}

	d.SetId("parameters")

	return resourceParametersRead(ctx, d, meta)
}

// resourceParametersRead retrieves the parameters from Warpgate and updates the Terraform state
func resourceParametersRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerMeta := meta.(*providerMeta)
	c := providerMeta.client

	var diags diag.Diagnostics

	params, err := c.GetParameters(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read parameters: %w", err))
	}

	if params == nil {
		d.SetId("")
		return diags
	}

	d.SetId("parameters")

	for _, field := range []struct {
		name  string
		value any
	}{
		{"allow_own_credential_management", params.AllowOwnCredentialManagement},
		{"rate_limit_bytes_per_second", params.RateLimitBytesPerSecond},
		{"ssh_client_auth_publickey", params.SSHClientAuthPublickey},
		{"ssh_client_auth_password", params.SSHClientAuthPassword},
		{"ssh_client_auth_keyboard_interactive", params.SSHClientAuthKeyboardInteractive},
		{"password_login_mode", params.PasswordLoginMode},
		{"ticket_self_service_enabled", params.TicketSelfServiceEnabled},
		{"ticket_auto_approve_existing_access", params.TicketAutoApproveExistingAccess},
		{"ticket_max_duration_seconds", int(params.TicketMaxDurationSeconds)},
		{"ticket_max_uses", params.TicketMaxUses},
		{"ticket_require_description", params.TicketRequireDescription},
		{"ticket_request_show_all_targets", params.TicketRequestShowAllTargets},
		{"target_click_action", params.TargetClickAction},
		{"show_session_menu", params.ShowSessionMenu},
		{"password_policy", flattenPasswordPolicy(params.PasswordPolicy)},
		{"max_api_token_duration_seconds", int(params.MaxAPITokenDurationSeconds)},
		{"record_scp", params.RecordSCP},
		{"login_protection_enabled", params.LoginProtectionEnabled},
		{"login_protection_retention_seconds", params.LoginProtectionRetentionSeconds},
		{"lp_ip_max_attempts", params.LPIPMaxAttempts},
		{"lp_ip_time_window_seconds", params.LPIPTimeWindowSeconds},
		{"lp_ip_base_block_duration_seconds", params.LPIPBaseBlockDurationSeconds},
		{"lp_ip_block_duration_multiplier", params.LPIPBlockDurationMultiplier},
		{"lp_ip_max_block_duration_seconds", params.LPIPMaxBlockDurationSeconds},
		{"lp_ip_cooldown_reset_seconds", params.LPIPCooldownResetSeconds},
		{"lp_user_max_attempts", params.LPUserMaxAttempts},
		{"lp_user_time_window_seconds", params.LPUserTimeWindowSeconds},
		{"lp_user_auto_unlock", params.LPUserAutoUnlock},
		{"lp_user_lockout_duration_seconds", params.LPUserLockoutDurationSeconds},
		{"lp_user_exempt_admins", params.LPUserExemptAdmins},
		{"ssh_banner", params.SSHBanner},
		{"web_ssh_enabled", params.WebSSHEnabled},
		{"analytics_consent", params.AnalyticsConsent},
		{"analytics_normal", params.AnalyticsNormal},
	} {
		if err := d.Set(field.name, field.value); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set %s: %w", field.name, err))
		}
	}

	return diags
}

// resourceParametersUpdate handles the update of Warpgate parameters
func resourceParametersUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerMeta := meta.(*providerMeta)
	c := providerMeta.client

	req := expandParametersUpdateRequest(d)

	_, err := c.UpdateParameters(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update parameters: %w", err))
	}

	return resourceParametersRead(ctx, d, meta)
}

func expandParametersUpdateRequest(d *schema.ResourceData) *client.ParametersUpdateRequest {
	req := &client.ParametersUpdateRequest{
		AllowOwnCredentialManagement:     d.Get("allow_own_credential_management").(bool),
		RateLimitBytesPerSecond:          optionalIntPointer(d, "rate_limit_bytes_per_second"),
		SSHClientAuthPublickey:           optionalBoolPointer(d, "ssh_client_auth_publickey"),
		SSHClientAuthPassword:            optionalBoolPointer(d, "ssh_client_auth_password"),
		SSHClientAuthKeyboardInteractive: optionalBoolPointer(d, "ssh_client_auth_keyboard_interactive"),
		TicketSelfServiceEnabled:         optionalBoolPointer(d, "ticket_self_service_enabled"),
		TicketAutoApproveExistingAccess:  optionalBoolPointer(d, "ticket_auto_approve_existing_access"),
		TicketMaxDurationSeconds:         optionalInt64Pointer(d, "ticket_max_duration_seconds"),
		TicketMaxUses:                    optionalIntPointer(d, "ticket_max_uses"),
		TicketRequireDescription:         optionalBoolPointer(d, "ticket_require_description"),
		TicketRequestShowAllTargets:      optionalBoolPointer(d, "ticket_request_show_all_targets"),
		TargetClickAction:                optionalStringPointer(d, "target_click_action"),
		ShowSessionMenu:                  optionalBoolPointer(d, "show_session_menu"),
		PasswordPolicy:                   expandPasswordPolicy(d),
		MaxAPITokenDurationSeconds:       optionalInt64Pointer(d, "max_api_token_duration_seconds"),
		RecordSCP:                        optionalBoolPointer(d, "record_scp"),
		LoginProtectionEnabled:           optionalBoolPointer(d, "login_protection_enabled"),
		LoginProtectionRetentionSeconds:  optionalIntPointer(d, "login_protection_retention_seconds"),
		LPIPMaxAttempts:                  optionalIntPointer(d, "lp_ip_max_attempts"),
		LPIPTimeWindowSeconds:            optionalIntPointer(d, "lp_ip_time_window_seconds"),
		LPIPBaseBlockDurationSeconds:     optionalIntPointer(d, "lp_ip_base_block_duration_seconds"),
		LPIPBlockDurationMultiplier:      optionalFloat64Pointer(d, "lp_ip_block_duration_multiplier"),
		LPIPMaxBlockDurationSeconds:      optionalIntPointer(d, "lp_ip_max_block_duration_seconds"),
		LPIPCooldownResetSeconds:         optionalIntPointer(d, "lp_ip_cooldown_reset_seconds"),
		LPUserMaxAttempts:                optionalIntPointer(d, "lp_user_max_attempts"),
		LPUserTimeWindowSeconds:          optionalIntPointer(d, "lp_user_time_window_seconds"),
		LPUserAutoUnlock:                 optionalBoolPointer(d, "lp_user_auto_unlock"),
		LPUserLockoutDurationSeconds:     optionalIntPointer(d, "lp_user_lockout_duration_seconds"),
		LPUserExemptAdmins:               optionalBoolPointer(d, "lp_user_exempt_admins"),
		SSHBanner:                        optionalStringPointer(d, "ssh_banner"),
		WebSSHEnabled:                    optionalBoolPointer(d, "web_ssh_enabled"),
		AnalyticsConsent:                 optionalStringPointer(d, "analytics_consent"),
		AnalyticsNormal:                  optionalBoolPointer(d, "analytics_normal"),
	}

	if passwordLoginMode := optionalStringPointer(d, "password_login_mode"); passwordLoginMode != nil {
		req.PasswordLoginMode = passwordLoginMode
	}

	return req
}

func flattenPasswordPolicy(policy client.PasswordPolicy) []any {
	return []any{
		map[string]any{
			"min_length":        policy.MinLength,
			"require_uppercase": policy.RequireUppercase,
			"require_lowercase": policy.RequireLowercase,
			"require_digits":    policy.RequireDigits,
			"require_special":   policy.RequireSpecial,
		},
	}
}

func expandPasswordPolicy(d *schema.ResourceData) *client.PasswordPolicy {
	if !configuredValueExists(d, "password_policy") {
		return nil
	}

	rawPolicies := d.Get("password_policy").([]any)
	if len(rawPolicies) == 0 || rawPolicies[0] == nil {
		return nil
	}

	rawPolicy := rawPolicies[0].(map[string]any)

	return &client.PasswordPolicy{
		MinLength:        rawPolicy["min_length"].(int),
		RequireUppercase: rawPolicy["require_uppercase"].(bool),
		RequireLowercase: rawPolicy["require_lowercase"].(bool),
		RequireDigits:    rawPolicy["require_digits"].(bool),
		RequireSpecial:   rawPolicy["require_special"].(bool),
	}
}

// resourceParametersDelete handles the deletion of Warpgate parameters
func resourceParametersDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	// For a global parameters resource, we don't actually delete it from the API.
	d.SetId("")

	return diags
}
