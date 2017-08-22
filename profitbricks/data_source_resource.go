package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
)

func dataSourceResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceResourceRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceResourceRead(d *schema.ResourceData, meta interface{}) error {
	profitbricks.SetDepth("5")

	results := []profitbricks.Resource{}

	resource_type := d.Get("resource_type").(string)
	resource_id := d.Get("resource_id").(string)

	if resource_type != "" && resource_id != "" {
		results = append(results, profitbricks.GetResourceByType(resource_type, resource_id))
		d.Set("resource_type", results[0].Type_)
		d.Set("resource_id", results[0].Id)
	} else if resource_type != "" {
		results = profitbricks.ListResourcesByType(resource_type).Items
		d.Set("resource_type", results[0].Type_)
	} else {
		results = profitbricks.ListResources().Items
	}

	if len(results) > 1 {
		return fmt.Errorf("There is more than one resource that match the search criteria")
	}

	if len(results) == 0 {
		return fmt.Errorf("There are no resources that match the search criteria")
	}

	d.SetId(results[0].Id)

	return nil
}
