// Package client provides types and functions for interacting with Warpgate API
package client

import (
	"context"
	"fmt"
	"net/http"
)

// ParameterValues represents the global parameters retrieved from Warpgate
type ParameterValues struct {
	AllowOwnCredentialManagement     bool           `json:"allow_own_credential_management"`
	RateLimitBytesPerSecond          int            `json:"rate_limit_bytes_per_second,omitempty"`
	SSHClientAuthPublickey           bool           `json:"ssh_client_auth_publickey"`
	SSHClientAuthPassword            bool           `json:"ssh_client_auth_password"`
	SSHClientAuthKeyboardInteractive bool           `json:"ssh_client_auth_keyboard_interactive"`
	PasswordLoginMode                string         `json:"password_login_mode"`
	TicketSelfServiceEnabled         bool           `json:"ticket_self_service_enabled"`
	TicketAutoApproveExistingAccess  bool           `json:"ticket_auto_approve_existing_access"`
	TicketMaxDurationSeconds         int64          `json:"ticket_max_duration_seconds,omitempty"`
	TicketMaxUses                    int            `json:"ticket_max_uses,omitempty"`
	TicketRequireDescription         bool           `json:"ticket_require_description"`
	TicketRequestShowAllTargets      bool           `json:"ticket_request_show_all_targets"`
	TargetClickAction                string         `json:"target_click_action,omitempty"`
	ShowSessionMenu                  bool           `json:"show_session_menu"`
	PasswordPolicy                   PasswordPolicy `json:"password_policy"`
	MaxAPITokenDurationSeconds       int64          `json:"max_api_token_duration_seconds,omitempty"`
	RecordSCP                        bool           `json:"record_scp"`
	LoginProtectionEnabled           bool           `json:"login_protection_enabled"`
	LoginProtectionRetentionSeconds  int            `json:"login_protection_retention_seconds"`
	LPIPMaxAttempts                  int            `json:"lp_ip_max_attempts"`
	LPIPTimeWindowSeconds            int            `json:"lp_ip_time_window_seconds"`
	LPIPBaseBlockDurationSeconds     int            `json:"lp_ip_base_block_duration_seconds"`
	LPIPBlockDurationMultiplier      float64        `json:"lp_ip_block_duration_multiplier"`
	LPIPMaxBlockDurationSeconds      int            `json:"lp_ip_max_block_duration_seconds"`
	LPIPCooldownResetSeconds         int            `json:"lp_ip_cooldown_reset_seconds"`
	LPUserMaxAttempts                int            `json:"lp_user_max_attempts"`
	LPUserTimeWindowSeconds          int            `json:"lp_user_time_window_seconds"`
	LPUserAutoUnlock                 bool           `json:"lp_user_auto_unlock"`
	LPUserLockoutDurationSeconds     int            `json:"lp_user_lockout_duration_seconds"`
	LPUserExemptAdmins               bool           `json:"lp_user_exempt_admins"`
	SSHBanner                        string         `json:"ssh_banner"`
	WebSSHEnabled                    bool           `json:"web_ssh_enabled"`
	AnalyticsConsent                 string         `json:"analytics_consent"`
	AnalyticsNormal                  bool           `json:"analytics_normal"`
}

// PasswordPolicy represents the Warpgate password policy settings.
type PasswordPolicy struct {
	MinLength        int  `json:"min_length"`
	RequireUppercase bool `json:"require_uppercase"`
	RequireLowercase bool `json:"require_lowercase"`
	RequireDigits    bool `json:"require_digits"`
	RequireSpecial   bool `json:"require_special"`
}

// ParametersUpdateRequest is the request payload for updating parameters
type ParametersUpdateRequest struct {
	AllowOwnCredentialManagement     bool            `json:"allow_own_credential_management"`
	RateLimitBytesPerSecond          *int            `json:"rate_limit_bytes_per_second,omitempty"`
	SSHClientAuthPublickey           *bool           `json:"ssh_client_auth_publickey,omitempty"`
	SSHClientAuthPassword            *bool           `json:"ssh_client_auth_password,omitempty"`
	SSHClientAuthKeyboardInteractive *bool           `json:"ssh_client_auth_keyboard_interactive,omitempty"`
	PasswordLoginMode                *string         `json:"password_login_mode,omitempty"`
	TicketSelfServiceEnabled         *bool           `json:"ticket_self_service_enabled,omitempty"`
	TicketAutoApproveExistingAccess  *bool           `json:"ticket_auto_approve_existing_access,omitempty"`
	TicketMaxDurationSeconds         *int64          `json:"ticket_max_duration_seconds,omitempty"`
	TicketMaxUses                    *int            `json:"ticket_max_uses,omitempty"`
	TicketRequireDescription         *bool           `json:"ticket_require_description,omitempty"`
	TicketRequestShowAllTargets      *bool           `json:"ticket_request_show_all_targets,omitempty"`
	TargetClickAction                *string         `json:"target_click_action,omitempty"`
	ShowSessionMenu                  *bool           `json:"show_session_menu,omitempty"`
	PasswordPolicy                   *PasswordPolicy `json:"password_policy,omitempty"`
	MaxAPITokenDurationSeconds       *int64          `json:"max_api_token_duration_seconds,omitempty"`
	RecordSCP                        *bool           `json:"record_scp,omitempty"`
	LoginProtectionEnabled           *bool           `json:"login_protection_enabled,omitempty"`
	LoginProtectionRetentionSeconds  *int            `json:"login_protection_retention_seconds,omitempty"`
	LPIPMaxAttempts                  *int            `json:"lp_ip_max_attempts,omitempty"`
	LPIPTimeWindowSeconds            *int            `json:"lp_ip_time_window_seconds,omitempty"`
	LPIPBaseBlockDurationSeconds     *int            `json:"lp_ip_base_block_duration_seconds,omitempty"`
	LPIPBlockDurationMultiplier      *float64        `json:"lp_ip_block_duration_multiplier,omitempty"`
	LPIPMaxBlockDurationSeconds      *int            `json:"lp_ip_max_block_duration_seconds,omitempty"`
	LPIPCooldownResetSeconds         *int            `json:"lp_ip_cooldown_reset_seconds,omitempty"`
	LPUserMaxAttempts                *int            `json:"lp_user_max_attempts,omitempty"`
	LPUserTimeWindowSeconds          *int            `json:"lp_user_time_window_seconds,omitempty"`
	LPUserAutoUnlock                 *bool           `json:"lp_user_auto_unlock,omitempty"`
	LPUserLockoutDurationSeconds     *int            `json:"lp_user_lockout_duration_seconds,omitempty"`
	LPUserExemptAdmins               *bool           `json:"lp_user_exempt_admins,omitempty"`
	SSHBanner                        *string         `json:"ssh_banner,omitempty"`
	WebSSHEnabled                    *bool           `json:"web_ssh_enabled,omitempty"`
	AnalyticsConsent                 *string         `json:"analytics_consent,omitempty"`
	AnalyticsNormal                  *bool           `json:"analytics_normal,omitempty"`
}

// GetParameters retrieves the global parameters from Warpgate
func (c *Client) GetParameters(ctx context.Context) (*ParameterValues, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/parameters", nil)
	if err != nil {
		return nil, err
	}

	var parameters ParameterValues
	if err := handleResponse(resp, &parameters); err != nil {
		return nil, err
	}

	return &parameters, nil
}

// UpdateParameters updates the global parameters in Warpgate
// Note: The API returns HTTP 201 with no response body, so we fetch the current state after update
func (c *Client) UpdateParameters(ctx context.Context, req *ParametersUpdateRequest) (*ParameterValues, error) {
	resp, err := c.doRequest(ctx, http.MethodPut, "/parameters", req)
	if err != nil {
		return nil, err
	}

	// PUT /parameters returns 201 with no body, so we need to discard the response
	// and fetch the current state instead
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to update parameters: HTTP %d", resp.StatusCode)
	}

	// Fetch the updated parameters
	return c.GetParameters(ctx)
}
