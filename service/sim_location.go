package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// SimLocation represents the location data returned for a SIM.
type SimLocation struct {
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
	PostalCode string `json:"postal_code"`
	Timestamp  string `json:"timestamp"`
}

// GetSimLocation fetches the last known location of a SIM by ICCID.
func (s *SimService) GetSimLocation(ctx context.Context, iccid string) ([]SimLocation, error) {
	if iccid == "" {
		return nil, fmt.Errorf("ICCID must not be empty")
	}

	// Although the path includes {iccid}, it's also required as a query param per OpenAPI spec.
	endpoint := fmt.Sprintf("%s/sims/%s/location?iccid=%s", apiPrefix, url.PathEscape(iccid), url.QueryEscape(iccid))

	respBody, err := s.Client.DoRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var result []SimLocation
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode SIM location: %w", err)
	}
	return result, nil
}
