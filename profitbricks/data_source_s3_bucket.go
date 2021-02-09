package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func dataSourceS3Bucket() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceS3ReadBucket,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:	  schema.TypeString,
				Required: true,
			},
			"contents": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"etag": {
							Type:	  schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:	  schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"display_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceS3ReadBucket(d *schema.ResourceData, meta interface{}) error {
	s3Client := meta.(Clients).S3Client
	name := d.Get("name").(string)
	s3Bucket, err := s3Client.GetBucket(name)
	if err != nil {
		switch err.(type) {
		case profitbricks.ApiError:
			return fmt.Errorf("error: %s - %s",
			 err.Error(), string((err.(profitbricks.ApiError)).Body()))

		case profitbricks.ClientError:
			return fmt.Errorf("error: %s", err.Error())

		}
	}

	if s3Bucket == nil {
		return fmt.Errorf("bucket not found")
	}

	contents:= make([]interface{}, len(s3Bucket.Contents), len(s3Bucket.Contents))
	for i, content := range s3Bucket.Contents {
		entry := make(map[string]interface{})
		entry["key"] = content.Key
		entry["etag"] = content.ETag
		entry["last_modified"] = content.LastModified
		entry["storage_class"] = content.StorageClass
		entry["size"] = content.Size

		ownerList := make([]interface{}, 1, 1)
		owner := make(map[string]string)
		owner["id"] = content.Owner.ID
		owner["display_name"] = content.Owner.DisplayName
		ownerList[0] = owner

		entry["owner"] = ownerList

		contents[i] = entry
	}

	d.SetId(s3Bucket.Name)

	if err := d.Set("contents", contents); err != nil {
		return err
	}

	return nil
}
