package profitbricks

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	profitbricks "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourceProfitBricksResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {datacenter}/{resource}", d.Id())
	}

	d.Set("datacenter_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceProfitBricksServerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) > 4 || len(parts) < 3 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {datacenter}/{server}/{primary_nic} or {datacenter}/{server}/{primary_nic}/{firewall}", d.Id())
	}

	d.Set("datacenter_id", parts[0])
	d.Set("primary_nic", parts[2])
	if len(parts) > 3 {
		d.Set("firewallrule_id", parts[3])
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func resourceProfitBricksK8sClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*profitbricks.Client)
	cluster, err := client.GetKubernetesCluster(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil, fmt.Errorf("Unable to find k8s cluster %q", d.Id())
			}
		}
		return nil, fmt.Errorf("Unable to retreive k8s cluster %q", d.Id())
	}

	log.Printf("[INFO] K8s cluster found: %+v", cluster)
	d.SetId(cluster.ID)
	d.Set("name", cluster.Properties.Name)
	d.Set("k8s_version", cluster.Properties.K8sVersion)

	if cluster.Properties.MaintenanceWindow != nil {
		d.Set("maintenance_window", []map[string]string{
			{
				"day_of_the_week": cluster.Properties.MaintenanceWindow.DayOfTheWeek,
				"time":            cluster.Properties.MaintenanceWindow.Time,
			},
		})
		log.Printf("[INFO] Setting maintenance window for k8s cluster %s to %+v...", d.Id(), cluster.Properties.MaintenanceWindow)
	}

	log.Printf("[INFO] Importing k8s cluster %q...", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceProfitBricksK8sNodepoolImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {k8sClusterId}/{k8sNodePoolId}", d.Id())
	}

	client := meta.(*profitbricks.Client)
	k8sNodepool, err := client.GetKubernetesNodePool(parts[0], parts[1])

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil, fmt.Errorf("Unable to find k8s node pool %q", d.Id())
			}
		}
		return nil, fmt.Errorf("Unable to retreive k8s node pool %q", d.Id())
	}

	log.Printf("[INFO] K8s node pool found: %+v", k8sNodepool)

	d.SetId(k8sNodepool.ID)
	d.Set("name", k8sNodepool.Properties.Name)
	d.Set("k8s_version", k8sNodepool.Properties.K8sVersion)
	d.Set("k8s_cluster_id", parts[0])
	d.Set("datacenter_id", k8sNodepool.Properties.DatacenterID)
	d.Set("cpu_family", k8sNodepool.Properties.CPUFamily)
	d.Set("availability_zone", k8sNodepool.Properties.AvailabilityZone)
	d.Set("storage_type", k8sNodepool.Properties.StorageType)
	d.Set("node_count", k8sNodepool.Properties.NodeCount)
	d.Set("cores_count", k8sNodepool.Properties.CoresCount)
	d.Set("ram_size", k8sNodepool.Properties.RAMSize)
	d.Set("storage_size", k8sNodepool.Properties.StorageSize)

	if k8sNodepool.Properties.MaintenanceWindow != nil {
		d.Set("maintenance_window", []map[string]string{
			{
				"day_of_the_week": k8sNodepool.Properties.MaintenanceWindow.DayOfTheWeek,
				"time":            k8sNodepool.Properties.MaintenanceWindow.Time,
			},
		})
		log.Printf("[INFO] Setting maintenance window for k8s node pool %s to %+v...", d.Id(), k8sNodepool.Properties.MaintenanceWindow)
	}

	log.Printf("[INFO] Importing k8s node pool %q...", d.Id())

	return []*schema.ResourceData{d}, nil
}

func resourceProfitBricksFirewallImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {datacenter}/{server}/{nic}/{firewall}", d.Id())
	}

	d.Set("datacenter_id", parts[0])
	d.Set("server_id", parts[1])
	d.Set("nic_id", parts[2])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}

func resourceProfitBricksNicImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {datacenter}/{server}/{nic}", d.Id())
	}

	d.Set("datacenter_id", parts[0])
	d.Set("server_id", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}

func convertSlice(slice []interface{}) []string {
	s := make([]string, len(slice))
	for i, v := range slice {
		s[i] = v.(string)
	}
	return s
}

func diffSlice(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}
