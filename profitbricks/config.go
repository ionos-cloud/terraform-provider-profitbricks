package profitbricks

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"

	"github.com/profitbricks/profitbricks-sdk-go/v5"
	"github.com/smartystreets/go-aws-auth"

)

// Config represents
type Config struct {
	Username string
	Password string
	Endpoint string
	S3Endpoint string
	S3Key string
	S3Secret string
	Retries  int
	Token    string
}

type Clients struct {
	ApiClient *profitbricks.Client
	S3Client *S3Client
}

func GetUserAgent(providerVersion string, terraformVersion string) string {
	return fmt.Sprintf("Ionos Cloud Terraform Provider v%s for Terraform v%s", providerVersion, terraformVersion)
}

// Client returns a new client for accessing profitbricks.
func (c *Config) Client(providerVersion, terraformVersion string) (*profitbricks.Client, error) {
	var client *profitbricks.Client
	if c.Token != "" {
		client = profitbricks.NewClientbyToken(c.Token)
	} else {
		client = profitbricks.NewClient(c.Username, c.Password)
	}
	client.SetUserAgent(GetUserAgent(providerVersion, terraformVersion))

	log.Printf("[DEBUG] Terraform client UA set to %s", client.GetUserAgent())

	client.SetDepth(5)

	if len(c.Endpoint) > 0 {
		client.SetHostURL(c.Endpoint)
	}
	return client, nil
}

func (c *Config) S3Client(providerVersion, terraformVersion string) (*S3Client, error) {
	if len(c.S3Key) == 0 || len(c.S3Secret) == 0 {
		return nil, nil
	}

	client, err := c.Client(providerVersion, terraformVersion)
	if err != nil {
		return nil, err
	}

	if len(c.S3Endpoint) > 0 {
		client.SetHostURL(c.S3Endpoint)
	}

	/* add auth signature headers to all requests */
	var authMiddleware resty.RequestMiddleware = func (restyClient *resty.Client, request *resty.Request) error {
		awsauth.Sign(request.RawRequest, awsauth.Credentials{AccessKeyID: c.S3Key, SecretAccessKey: c.S3Secret})
		return nil
	}
	client.Client.OnBeforeRequest(authMiddleware)

	return NewS3Client(client), nil
}
