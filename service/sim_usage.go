package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// SimsUsageItem represents usage information for a single SIM.
type SimsUsageItem struct {
	ICCID                string `json:"iccid"`
	MSISDN               string `json:"msisdn"`
	MonthToDateUp        string `json:"month_to_date_bytes_up"`
	MonthToDateDown      string `json:"month_to_date_bytes_down"`
	MonthToDateVoiceUp   string `json:"month_to_date_voice_up"`
	MonthToDateVoiceDown string `json:"month_to_date_voice_down"`
	MonthToDateSmsUp     string `json:"month_to_date_sms_up"`
	MonthToDateSmsDown   string `json:"month_to_date_sms_down"`
	LastSeen             string `json:"last_seen"`
	InCurrentSession     bool   `json:"in_current_session"`
}

// SimsUsageResponse represents the top-level response from /api/v3/sims/usage
type SimsUsageResponse struct {
	Sims []SimsUsageItem `json:"sims"`
}

// GetSimUsage fetches current month usage stats for SIMs, optionally filtered by ICCID(s).
func (s *SimService) GetSimUsage(ctx context.Context, iccids []string, page, limit int) (*SimsUsageResponse, error) {
	params := url.Values{}
	if len(iccids) > 0 {
		params.Set("iccid", joinComma(iccids))
	}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	endpoint := simUsageEndpoint
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	respBody, err := s.Client.DoRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var result SimsUsageResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode SIM usage: %w", err)
	}
	return &result, nil
}

// joinComma joins a slice of strings with commas for query parameters.
func joinComma(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return url.QueryEscape(fmt.Sprintf("%s", values[0])) + func() string {
		if len(values) == 1 {
			return ""
		}
		result := ""
		for _, v := range values[1:] {
			result += "," + url.QueryEscape(v)
		}
		return result
	}()
}
