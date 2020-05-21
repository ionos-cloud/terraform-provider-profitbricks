package profitbricks

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	profitbricks "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourcek8sCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourcek8sClusterCreate,
		Read:   resourcek8sClusterRead,
		Update: resourcek8sClusterUpdate,
		Delete: resourcek8sClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProfitBricksK8sClusterImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"k8s_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"maintenance_window": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"day_of_the_week": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourcek8sClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	cluster := profitbricks.KubernetesCluster{
		Properties: &profitbricks.KubernetesClusterProperties{
			Name: d.Get("name").(string),
		},
	}

	if k8svVal, k8svOk := d.GetOk("k8s_version"); k8svOk {
		log.Printf("[INFO] Setting K8s version to : %s", k8svVal.(string))
		cluster.Properties.K8sVersion = k8svVal.(string)
	}

	if _, mwOk := d.GetOk("maintenance_window.0"); mwOk {
		cluster.Properties.MaintenanceWindow = &profitbricks.MaintenanceWindow{}
	}

	if mtVal, mtOk := d.GetOk("maintenance_window.0.time"); mtOk {
		log.Printf("[INFO] Setting Maintenance Window Time to : %s", mtVal.(string))
		cluster.Properties.MaintenanceWindow.Time = mtVal.(string)
	}

	if mdVal, mdOk := d.GetOk("maintenance_window.0.day_of_the_week"); mdOk {
		cluster.Properties.MaintenanceWindow.DayOfTheWeek = mdVal.(string)
	}

	createdCluster, err := client.CreateKubernetesCluster(cluster)

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error creating k8s cluster: %s", err)
	}

	d.SetId(createdCluster.ID)
	log.Printf("[INFO] Created k8s cluster: %s", d.Id())

	for !k8sClusterReady(client, d) {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)
	}

	return resourcek8sClusterRead(d, meta)
}

func resourcek8sClusterRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*profitbricks.Client)
	cluster, err := client.GetKubernetesCluster(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error while fetching k8s cluster %s: %s", d.Id(), err)
	}

	log.Printf("[INFO] Successfully retreived cluster %s", d.Id())

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

	return nil
}

func resourcek8sClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	request := profitbricks.UpdatedKubernetesCluster{}

	request.Properties = &profitbricks.KubernetesClusterProperties{
		Name: d.Get("name").(string),
	}

	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		log.Printf("[INFO] k8s cluster name changed from %+v to %+v", oldName, newName)
		request.Properties.Name = newName.(string)
	}

	log.Printf("[INFO] Attempting update cluster Id %s", d.Id())

	if d.HasChange("k8s_version") {
		oldk8sVersion, newk8sVersion := d.GetChange("k8s_version")
		log.Printf("[INFO] k8s version changed from %+v to %+v", oldk8sVersion, newk8sVersion)
		if newk8sVersion != nil {
			request.Properties.K8sVersion = newk8sVersion.(string)
		}
	}

	if d.HasChange("maintenance_window.0") {

		_, newMw := d.GetChange("maintenance_window.0")

		if newMw.(map[string]interface{}) != nil {

			updateMaintenanceWindow := false
			maintenanceWindow := &profitbricks.MaintenanceWindow{}

			if d.HasChange("maintenance_window.0.day_of_the_week") {
				oldMd, newMd := d.GetChange("maintenance_window.0.day_of_the_week")
				if newMd.(string) != "" {
					log.Printf("[INFO] k8s maintenance window DOW changed from %+v to %+v", oldMd, newMd)
					updateMaintenanceWindow = true
					maintenanceWindow.DayOfTheWeek = newMd.(string)
				}
			}

			if d.HasChange("maintenance_window.0.time") {

				oldMt, newMt := d.GetChange("maintenance_window.0.time")
				if newMt.(string) != "" {
					log.Printf("[INFO] k8s maintenance window time changed from %+v to %+v", oldMt, newMt)
					updateMaintenanceWindow = true
					maintenanceWindow.Time = newMt.(string)
				}
			}

			if updateMaintenanceWindow == true {
				request.Properties.MaintenanceWindow = maintenanceWindow
			}
		}
	}

	_, err := client.UpdateKubernetesCluster(d.Id(), request)

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while updating k8s cluster: %s", err)
		}
		return fmt.Errorf("Error while updating k8s cluster %s: %s", d.Id(), err)
	}

	for !k8sClusterReady(client, d) {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)
	}

	return resourcek8sClusterRead(d, meta)
}

func resourcek8sClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	_, err := client.DeleteKubernetesCluster(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while deleting k8s cluster: %s", err)
		}

		return fmt.Errorf("Error while deleting k8s cluster %s: %s", d.Id(), err)
	}

	for !k8sClusterDeleted(client, d) {
		log.Printf("[INFO] Waiting for cluster %s to be deleted...", d.Id())
		time.Sleep(5 * time.Second)
	}

	log.Printf("[INFO] Successfully deleted k8s cluster: %s", d.Id())

	return nil
}

func k8sClusterReady(client *profitbricks.Client, d *schema.ResourceData) bool {
	subjectCluster, err := client.GetKubernetesCluster(d.Id())

	if err == nil {
		return subjectCluster.Metadata.State == "ACTIVE"
	}
	return false
}

func k8sClusterDeleted(client *profitbricks.Client, d *schema.ResourceData) bool {
	_, err := client.GetKubernetesCluster(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				return true
			}
		}
	}
	return false
}
