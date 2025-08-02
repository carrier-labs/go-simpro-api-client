package service

import "github.com/carrier-labs/go-simpro-api-client/client"

// Internal endpoint path for Tariff operations
const tariffListEndpoint = apiPrefix + "/tariffs"

// TariffService provides access to /api/v3/tariffs and related endpoints.
type TariffService struct {
	Client *client.Client
}

// NewTariffService creates a new instance of TariffService.
func NewTariffService(c *client.Client) *TariffService {
	return &TariffService{Client: c}
}
