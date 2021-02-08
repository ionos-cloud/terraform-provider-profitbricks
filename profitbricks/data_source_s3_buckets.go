package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
	"log"
)

func dataSourceS3Buckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceS3ReadBuckets,
		Schema: map[string]*schema.Schema{
			"owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owner_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"buckets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceS3ReadBuckets(d *schema.ResourceData, meta interface{}) error {
	s3Client := meta.(Clients).S3Client
	s3Buckets, err := s3Client.ListS3Buckets()
	if err != nil {
		log.Printf("--------> %v\n", err)
		return fmt.Errorf("error: %s - %s",
			err.Error(), string((err.(profitbricks.ApiError)).Body()))
	}

	log.Printf("[DEBUG] S3 Owner ID = %s", s3Buckets.Owner.ID)
	if err := d.Set("owner_id", s3Buckets.Owner.ID); err != nil {
		return err
	}

	log.Printf("[DEBUG] S3 Owner Name = %s", s3Buckets.Owner.DisplayName)
	if err := d.Set("owner_name", s3Buckets.Owner.DisplayName); err != nil {
		return err
	}

	return nil
}
