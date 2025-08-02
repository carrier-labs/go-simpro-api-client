package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/carrier-labs/go-simpro-api-client/models"
)

// TariffListFilter optionally filters the tariffs by account number(s).
type TariffListFilter struct {
	AccountNumbers []string // One or more billing account numbers
	Page           int
	Limit          int
}

// TariffListItem represents a single tariff as returned by /api/v3/tariffs
type TariffListItem struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	MNO            models.MNO `json:"mno"`
	ContractLength int        `json:"contract_length"`
	CustomerName   string     `json:"customer_name"`
	AccountNumber  string     `json:"account_number"`
	Bearers        []struct {
		Name string `json:"name"`
	} `json:"bearers"`
}

// TariffListResponse represents the response from /api/v3/tariffs
type TariffListResponse []TariffListItem

// GetTariffs retrieves available tariffs with optional filtering.
func (s *TariffService) GetTariffs(ctx context.Context, filter *TariffListFilter) (*TariffListResponse, error) {
	params := url.Values{}
	if filter != nil {
		if len(filter.AccountNumbers) > 0 {
			params.Set("account_numbers", joinComma(filter.AccountNumbers))
		}
		if filter.Page > 0 {
			params.Set("page", fmt.Sprintf("%d", filter.Page))
		}
		if filter.Limit > 0 {
			params.Set("limit", fmt.Sprintf("%d", filter.Limit))
		}
	}

	endpoint := tariffListEndpoint
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	respBody, err := s.Client.DoRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var result TariffListResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode tariffs: %w", err)
	}
	return &result, nil
}
