package profitbricks

import (
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

const DefaultS3ApiUrl = "https://s3-de-central.profitbricks.com"

type S3Client struct {
	*profitbricks.Client
}

func NewS3Client(client *profitbricks.Client) *S3Client {
	return &S3Client{client}
}

