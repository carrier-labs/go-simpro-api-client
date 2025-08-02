package service

import "github.com/carrier-labs/go-simpro-api-client/client"

// SimService provides access to all SIM-related endpoints (list, usage, details, etc).
type SimService struct {
	Client *client.Client
}

// NewSimService creates a new SimService instance.
func NewSimService(c *client.Client) *SimService {
	return &SimService{Client: c}
}

const (
	simListEndpoint  = apiPrefix + "/sims"
	simUsageEndpoint = apiPrefix + "/sims/usage"
)
