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
	connection := meta.(*profitbricks.Client)

	var results []profitbricks.Resource

	resource_type := d.Get("resource_type").(string)
	resource_id := d.Get("resource_id").(string)

	if resource_type != "" && resource_id != "" {
		results, err := connection.GetResourceByType(resource_type, resource_id)
		if err != nil {
			return fmt.Errorf("An error occured while fetching resource by type %s", err)
		}

		d.Set("resource_type", results.PBType)
		d.Set("resource_id", results.ID)
	} else if resource_type != "" {
		items, err := connection.ListResourcesByType(resource_type)
		if err != nil {
			return fmt.Errorf("An error occured while fetching resources by type %s", err)
		}

		results = items.Items
		d.Set("resource_type", results[0].PBType)
	} else {
		items, err := connection.ListResources()
		if err != nil {
			return fmt.Errorf("An error occured while fetching resources %s", err)
		}
		results = items.Items
	}

	if len(results) > 1 {
		return fmt.Errorf("There is more than one resource that match the search criteria")
	}

	if len(results) == 0 {
		return fmt.Errorf("There are no resources that match the search criteria")
	}

	d.SetId(results[0].ID)

	return nil
}
