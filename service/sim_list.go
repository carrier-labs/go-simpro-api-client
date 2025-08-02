package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// SimsListFilter defines the supported query parameters for GET /api/v3/sims
type SimsListFilter struct {
	Status        string
	AccountNumber string
	TariffName    string
	MNO           string
	CustomField1  string
	// Add more fields if needed
}

// SimsListItem represents a single SIM record in the response from GET /api/v3/sims
type SimsListItem struct {
	ID             int    `json:"id"`
	ICCID          string `json:"iccid"`
	EID            string `json:"eid"`
	MSISDN         string `json:"msisdn"`
	IMSI           string `json:"imsi"`
	Status         string `json:"status"`
	WorkflowStatus string `json:"workflow_status"`
}

// SimsListResponse represents the response structure from GET /api/v3/sims
type SimsListResponse struct {
	Sims     []SimsListItem `json:"sims"`
	SimCount int            `json:"sim_count"`
}

// GetSims fetches a list of SIMs with optional filters applied.
func (s *SimService) GetSims(ctx context.Context, filter *SimsListFilter) (*SimsListResponse, error) {
	params := url.Values{}
	if filter != nil {
		if filter.Status != "" {
			params.Set("status", filter.Status)
		}
		if filter.AccountNumber != "" {
			params.Set("account_number", filter.AccountNumber)
		}
		if filter.TariffName != "" {
			params.Set("tariff_name", filter.TariffName)
		}
		if filter.MNO != "" {
			params.Set("mno", filter.MNO)
		}
		if filter.CustomField1 != "" {
			params.Set("custom_field1", filter.CustomField1)
		}
	}

	endpoint := simListEndpoint
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	respBody, err := s.Client.DoRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var result SimsListResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode SIMs: %w", err)
	}
	return &result, nil
}
