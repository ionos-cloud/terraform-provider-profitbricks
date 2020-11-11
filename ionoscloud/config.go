package ionoscloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	ionoscloud "github.com/profitbricks/profitbricks-sdk-go/v5"
)

// Config represents
type Config struct {
	Username string
	Password string
	Endpoint string
	Retries  int
	Token    string
}

// Client returns a new client for accessing IonosCloud.
func (c *Config) Client(terraformVersion string) (*ionoscloud.Client, error) {
	var client *ionoscloud.Client
	if c.Token != "" {
		client = ionoscloud.NewClientbyToken(c.Token)
	} else {
		client = ionoscloud.NewClient(c.Username, c.Password)
	}
	client.SetUserAgent(httpclient.TerraformUserAgent(terraformVersion))

	log.Printf("[DEBUG] Terraform client UA set to %s", client.GetUserAgent())

	client.SetDepth(5)

	if len(c.Endpoint) > 0 {
		client.SetHostURL(c.Endpoint)
	}
	return client, nil
}
