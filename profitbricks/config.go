package profitbricks

import (
	"github.com/hashicorp/terraform/terraform"
	"github.com/profitbricks/profitbricks-sdk-go"
)

type Config struct {
	Username string
	Password string
	Endpoint string
	Retries  int
}

// Client() returns a new client for accessing ProfitBricks.
func (c *Config) Client() (*Config, error) {
	profitbricks.SetAuth(c.Username, c.Password)
	profitbricks.SetDepth("5")
	profitbricks.SetUserAgent(terraform.UserAgentString())

	if len(c.Endpoint) > 0 {
		profitbricks.SetEndpoint(c.Endpoint)
	}
	return c, nil
}
