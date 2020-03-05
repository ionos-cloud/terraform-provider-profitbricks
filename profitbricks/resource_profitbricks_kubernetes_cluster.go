package profitbricks

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: nil,
		Read:   nil,
		Update: nil,
		Delete: nil,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
