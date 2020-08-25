package profitbricks

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	profitbricks "github.com/profitbricks/profitbricks-sdk-go/v5"
)

// Provider returns a schema.Provider for ProfitBricks.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("PROFITBRICKS_USERNAME", nil),
				Description:   "ProfitBricks username for API operations. If token is provided, token is preferred",
				ConflictsWith: []string{"token"},
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("PROFITBRICKS_PASSWORD", nil),
				Description:   "ProfitBricks password for API operations. If token is provided, token is preferred",
				ConflictsWith: []string{"token"},
			},
			"token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("PROFITBRICKS_TOKEN", nil),
				Description:   "Profitbricks bearer token for API operations.",
				ConflictsWith: []string{"username", "password"},
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROFITBRICKS_API_URL", ""),
				Description: "ProfitBricks REST API URL.",
			},
			"retries": {
				Type:       schema.TypeInt,
				Optional:   true,
				Default:    50,
				Deprecated: "Timeout is used instead of this functionality",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"profitbricks_datacenter":           resourceProfitBricksDatacenter(),
			"profitbricks_ipblock":              resourceProfitBricksIPBlock(),
			"profitbricks_firewall":             resourceProfitBricksFirewall(),
			"profitbricks_lan":                  resourceProfitBricksLan(),
			"profitbricks_loadbalancer":         resourceProfitBricksLoadbalancer(),
			"profitbricks_nic":                  resourceProfitBricksNic(),
			"profitbricks_server":               resourceProfitBricksServer(),
			"profitbricks_volume":               resourceProfitBricksVolume(),
			"profitbricks_group":                resourceProfitBricksGroup(),
			"profitbricks_share":                resourceProfitBricksShare(),
			"profitbricks_user":                 resourceProfitBricksUser(),
			"profitbricks_snapshot":             resourceProfitBricksSnapshot(),
			"profitbricks_ipfailover":           resourceProfitBricksLanIPFailover(),
			"profitbricks_k8s_cluster":          resourcek8sCluster(),
			"profitbricks_k8s_node_pool":        resourcek8sNodePool(),
			"profitbricks_private_crossconnect": resourcePrivateCrossConnect(),
			"profitbricks_backup_unit":          resourceBackupUnit(),
			"profitbricks_s3_key":               resourceS3Key(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"profitbricks_datacenter": dataSourceDataCenter(),
			"profitbricks_location":   dataSourceLocation(),
			"profitbricks_image":      dataSourceImage(),
			"profitbricks_resource":   dataSourceResource(),
			"profitbricks_snapshot":   dataSourceSnapshot(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {

		terraformVersion := provider.TerraformVersion

		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		log.Printf("[DEBUG] Setting terraformVersion to %s", terraformVersion)

		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {

	username, usernameOk := d.GetOk("username")
	password, passwordOk := d.GetOk("password")
	token, tokenOk := d.GetOk("token")

	if !tokenOk {
		if !usernameOk {
			return nil, fmt.Errorf("Neither ProfitBricks token, nor ProfitBricks username has been provided")
		}

		if !passwordOk {
			return nil, fmt.Errorf("Neither ProfitBricks token, nor ProfitBricks password has been provided")
		}
	} else {
		if usernameOk || passwordOk {
			return nil, fmt.Errorf("Only provide ProfitBricks token OR ProfitBricks username/password.")
		}
	}

	config := Config{
		Username: username.(string),
		Password: password.(string),
		Endpoint: cleanURL(d.Get("endpoint").(string)),
		Retries:  d.Get("retries").(int),
		Token:    token.(string),
	}

	return config.Client(terraformVersion)
}

// cleanURL makes sure trailing slash does not corrupte the state
func cleanURL(url string) string {
	length := len(url)
	if length > 1 && url[length-1] == '/' {
		url = url[:length-1]
	}

	return url
}

// getStateChangeConf gets the default configuration for tracking a request progress
func getStateChangeConf(meta interface{}, d *schema.ResourceData, location string, timeoutType string) *resource.StateChangeConf {
	stateConf := &resource.StateChangeConf{
		Pending:        resourcePendingStates,
		Target:         resourceTargetStates,
		Refresh:        resourceStateRefreshFunc(meta, location),
		Timeout:        d.Timeout(timeoutType),
		MinTimeout:     10 * time.Second,
		Delay:          10 * time.Second, // Wait 10 secs before starting
		NotFoundChecks: 600,              //Setting high number, to support long timeouts
	}

	return stateConf
}

type RequestFailedError struct {
	msg string
}

func (e RequestFailedError) Error() string {
	return e.msg
}

func IsRequestFailed(err error) bool {
	_, ok := err.(RequestFailedError)
	return ok
}

// resourceStateRefreshFunc tracks progress of a request
func resourceStateRefreshFunc(meta interface{}, path string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*profitbricks.Client)

		fmt.Printf("[INFO] Checking PATH %s", path)
		if path == "" {
			return nil, "", fmt.Errorf("Can not check a state when path is empty")
		}

		request, err := client.GetRequestStatus(path)

		if err != nil {
			return nil, "", fmt.Errorf("Request failed with following error: %s", err)
		}

		if request.Metadata.Status == "FAILED" {
			return nil, "", RequestFailedError{fmt.Sprintf("Request failed with following error: %s", request.Metadata.Message)}
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
