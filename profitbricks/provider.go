package profitbricks

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/profitbricks/profitbricks-sdk-go"
)

// Provider returns a schema.Provider for ProfitBricks.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROFITBRICKS_USERNAME", nil),
				Description: "ProfitBricks username for API operations.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROFITBRICKS_PASSWORD", nil),
				Description: "ProfitBricks password for API operations.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROFITBRICKS_API_URL", profitbricks.Endpoint),
				Description: "ProfitBricks REST API URL.",
			},
			"retries": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  50,
				Removed:  "Timeout is used instead of this functionality",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"profitbricks_datacenter":   resourceProfitBricksDatacenter(),
			"profitbricks_ipblock":      resourceProfitBricksIPBlock(),
			"profitbricks_firewall":     resourceProfitBricksFirewall(),
			"profitbricks_lan":          resourceProfitBricksLan(),
			"profitbricks_loadbalancer": resourceProfitBricksLoadbalancer(),
			"profitbricks_nic":          resourceProfitBricksNic(),
			"profitbricks_server":       resourceProfitBricksServer(),
			"profitbricks_volume":       resourceProfitBricksVolume(),
			"profitbricks_group":        resourceProfitBricksGroup(),
			"profitbricks_share":        resourceProfitBricksShare(),
			"profitbricks_user":         resourceProfitBricksUser(),
			"profitbricks_snapshot":     resourceProfitBricksSnapshot(),
			"profitbricks_ipfailover":   resourceProfitBricksLanIPFailover(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"profitbricks_datacenter": dataSourceDataCenter(),
			"profitbricks_location":   dataSourceLocation(),
			"profitbricks_image":      dataSourceImage(),
			"profitbricks_resource":   dataSourceResource(),
			"profitbricks_snapshot":   dataSourceSnapshot(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	if _, ok := d.GetOk("username"); !ok {
		return nil, fmt.Errorf("ProfitBricks username has not been provided")
	}

	if _, ok := d.GetOk("password"); !ok {
		return nil, fmt.Errorf("ProfitBricks password has not been provided")
	}

	config := Config{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Endpoint: cleanURL(d.Get("endpoint").(string)),
		Retries:  d.Get("retries").(int),
	}

	return config.Client()
}

// cleanURL makes sure trailing slash does not corrupte the state
func cleanURL(url string) string {
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	return url
}

// getStateChangeConf gets the default configuration for tracking a request progress
func getStateChangeConf(meta interface{}, d *schema.ResourceData, location string, timeoutType string) *resource.StateChangeConf {
	stateConf := &resource.StateChangeConf{
		Pending:    resourcePendingStates,
		Target:     resourceTargetStates,
		Refresh:    resourceStateRefreshFunc(meta, location),
		Timeout:    d.Timeout(timeoutType),
		MinTimeout: 10 * time.Second,
		Delay:      10 * time.Second, // Wait 10 secs before starting
	}

	return stateConf
}

// resourceStateRefreshFunc tracks progress of a request
func resourceStateRefreshFunc(meta interface{}, path string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		if path == "" {
			return nil, "", fmt.Errorf("Can not check a state when path is empty")
		}

		request := profitbricks.GetRequestStatus(path)

		if request.Metadata.Status == "FAILED" {

			return nil, "", fmt.Errorf("Request failed with following error: %s", request.Metadata.Message)
		}

		if request.Metadata.Status == "DONE" {
			return request, "DONE", nil
		}

		return nil, request.Metadata.Status, nil
	}
}

// resourcePendingStates defines states of working in progress
var resourcePendingStates = []string{
	"RUNNING",
	"QUEUED",
}

// resourceTargetStates defines states of completion
var resourceTargetStates = []string{
	"DONE",
}

// resourceDefaultTimeouts sets default value for each Timeout type
var resourceDefaultTimeouts = schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(60 * time.Minute),
	Update:  schema.DefaultTimeout(60 * time.Minute),
	Delete:  schema.DefaultTimeout(60 * time.Minute),
	Default: schema.DefaultTimeout(60 * time.Minute),
}
