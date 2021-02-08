package profitbricks

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	awsauth "github.com/smartystreets/go-aws-auth"
	"log"
	"net/http"

	"github.com/profitbricks/profitbricks-sdk-go/v5"
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

	client.Client.Debug = true
	/* add auth signature headers to all requests */
	var authMiddleware resty.PreRequestHook = func (restyClient *resty.Client, request *http.Request) error {
		log.Printf(">>> signing\n")
		awsauth.Sign(request, awsauth.Credentials{AccessKeyID: c.S3Key, SecretAccessKey: c.S3Secret})
		log.Printf("%v\n", request.Header)
		return nil
	}
	client.Client.SetPreRequestHook(authMiddleware)
	client.Client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		log.Printf(">>> RESPONSE")
		log.Printf("%s\n", string(response.Body()))
		return nil
	})

	return NewS3Client(client), nil
}
