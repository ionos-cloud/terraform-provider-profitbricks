package profitbricks

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	profitbricks "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourcek8sNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourcek8sNodePoolCreate,
		Read:   resourcek8sNodePoolRead,
		Update: resourcek8sNodePoolUpdate,
		Delete: resourcek8sNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProfitBricksK8sNodepoolImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The desired name for the node pool",
				Required:    true,
			},
			"k8s_version": {
				Type:        schema.TypeString,
				Description: "The desired kubernetes version",
				Required:    true,
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "A maintenance window comprise of a day of the week and a time for maintenance to be allowed",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Description: "A clock time in the day when maintenance is allowed",
							Required:    true,
						},
						"day_of_the_week": {
							Type:        schema.TypeString,
							Description: "Day of the week when maintenance is allowed",
							Required:    true,
						},
					},
				},
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The UUID of the VDC",
				Required:    true,
			},
			"k8s_cluster_id": {
				Type:        schema.TypeString,
				Description: "The UUID of an existing kubernetes cluster",
				Required:    true,
			},
			"cpu_family": {
				Type:        schema.TypeString,
				Description: "CPU Family",
				Required:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "The compute availability zone in which the nodes should exist",
				Required:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "Storage type to use",
				Required:    true,
			},
			"node_count": {
				Type:        schema.TypeInt,
				Description: "The number of nodes in this node pool",
				Required:    true,
			},
			"cores_count": {
				Type:        schema.TypeInt,
				Description: "CPU cores count",
				Required:    true,
			},
			"ram_size": {
				Type:        schema.TypeInt,
				Description: "The amount of RAM in MB",
				Required:    true,
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Description: "The total allocated storage capacity of a node in GB",
				Required:    true,
			},
		},
	}
}

func resourcek8sNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	k8sNodepool := profitbricks.KubernetesNodePool{
		Properties: &profitbricks.KubernetesNodePoolProperties{
			Name:             d.Get("name").(string),
			DatacenterID:     d.Get("datacenter_id").(string),
			K8sVersion:       d.Get("k8s_version").(string),
			AvailabilityZone: d.Get("availability_zone").(string),
			CPUFamily:        d.Get("cpu_family").(string),
			StorageType:      d.Get("storage_type").(string),
			NodeCount:        uint32(d.Get("node_count").(int)),
			CoresCount:       uint32(d.Get("cores_count").(int)),
			StorageSize:      uint32(d.Get("storage_size").(int)),
			RAMSize:          uint32(d.Get("ram_size").(int)),
		},
	}

	if _, mwOk := d.GetOk("maintenance_window.0"); mwOk {
		k8sNodepool.Properties.MaintenanceWindow = &profitbricks.MaintenanceWindow{}
	}

	if mtVal, mtOk := d.GetOk("maintenance_window.0.time"); mtOk {
		log.Printf("[INFO] Setting Maintenance Window Time to : %s", mtVal.(string))
		k8sNodepool.Properties.MaintenanceWindow.Time = mtVal.(string)
	}

	if mdVal, mdOk := d.GetOk("maintenance_window.0.day_of_the_week"); mdOk {
		k8sNodepool.Properties.MaintenanceWindow.DayOfTheWeek = mdVal.(string)
	}

	createdNodepool, err := client.CreateKubernetesNodePool(d.Get("k8s_cluster_id").(string), k8sNodepool)

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error creating k8s node pool: %s", err)
	}

	d.SetId(createdNodepool.ID)

	log.Printf("[INFO] Successfully created k8s node pool: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be ready...", d.Id())
		time.Sleep(10 * time.Second)

		nodepoolReady, rsErr := k8sNodepoolReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of k8s node pool %s: %s", d.Id(), rsErr)
		}

		if nodepoolReady && rsErr == nil {
			log.Printf("[INFO] k8s node pool ready: %s", d.Id())
			break
		}
	}

	return resourcek8sNodePoolRead(d, meta)
}

func resourcek8sNodePoolRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*profitbricks.Client)
	k8sNodepool, err := client.GetKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id())

	if err != nil {
		log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
	}

	log.Printf("[INFO] Successfully retreived k8s node pool %s: %+v", d.Id(), k8sNodepool)

	return nil
}

func resourcek8sNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*profitbricks.Client)
	request := profitbricks.KubernetesNodePool{}

	request.Properties = &profitbricks.KubernetesNodePoolProperties{
		NodeCount: uint32(d.Get("node_count").(int)),
	}

	if d.HasChange("node_count") {
		oldNc, newNc := d.GetChange("name")
		log.Printf("[INFO] k8s node pool node count changed from %+v to %+v", oldNc, newNc)
		if oldNc.(int) != newNc.(int) {
			request.Properties.NodeCount = uint32(newNc.(int))
		}
	}

	if d.HasChange("k8s_version") {
		oldk8sVersion, newk8sVersion := d.GetChange("k8s_version")
		log.Printf("[INFO] k8s pool k8s version changed from %+v to %+v", oldk8sVersion, newk8sVersion)
		if newk8sVersion != nil {
			request.Properties.K8sVersion = newk8sVersion.(string)
		}
	}

	if d.HasChange("maintenance_window.0") {

		_, newMw := d.GetChange("maintenance_window.0")

		if newMw.(map[string]interface{}) != nil {

			updateMaintenanceWindow := false
			maintenanceWindow := &profitbricks.MaintenanceWindow{
				DayOfTheWeek: d.Get("maintenance_window.0.day_of_the_week").(string),
				Time:         d.Get("maintenance_window.0.time").(string),
			}

			if d.HasChange("maintenance_window.0.day_of_the_week") {

				oldMd, newMd := d.GetChange("maintenance_window.0.day_of_the_week")
				if newMd.(string) != "" {
					log.Printf("[INFO] k8s node pool maintenance window DOW changed from %+v to %+v", oldMd, newMd)
					updateMaintenanceWindow = true
					maintenanceWindow.DayOfTheWeek = newMd.(string)
				}
			}

			if d.HasChange("maintenance_window.0.time") {
				oldMt, newMt := d.GetChange("maintenance_window.0.time")
				if newMt.(string) != "" {
					log.Printf("[INFO] k8s node pool maintenance window time changed from %+v to %+v", oldMt, newMt)
					updateMaintenanceWindow = true
					maintenanceWindow.Time = newMt.(string)
				}
			}

			if updateMaintenanceWindow == true {
				request.Properties.MaintenanceWindow = maintenanceWindow
			}
		}
	}

	b, jErr := json.Marshal(request)

	if jErr == nil {
		log.Printf("[INFO] Update req: %s", string(b))
	}

	_, err := client.UpdateKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id(), request)

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while updating k8s node pool: %s", err)
		}
		return fmt.Errorf("Error while updating k8s node pool %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be ready...", d.Id())
		time.Sleep(10 * time.Second)

		nodepoolReady, rsErr := k8sNodepoolReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of k8s node pool %s: %s", d.Id(), rsErr)
		}

		if nodepoolReady && rsErr == nil {
			log.Printf("[INFO] k8s node pool ready: %s", d.Id())
			break
		}
	}

	return resourcek8sNodePoolRead(d, meta)
}

func resourcek8sNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	_, err := client.DeleteKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while deleting k8s node pool: %s", err)
		}

		return fmt.Errorf("Error while deleting k8s node pool %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for k8s node pool %s to be deleted...", d.Id())
		time.Sleep(10 * time.Second)

		nodepoolDeleted, dsErr := k8sNodepoolDeleted(client, d)

		if dsErr != nil {
			return fmt.Errorf("Error while checking deletion status of k8s node pool %s: %s", d.Id(), dsErr)
		}

		if nodepoolDeleted && dsErr == nil {
			log.Printf("[INFO] Successfully deleted k8s node pool: %s", d.Id())
			break
		}
	}

	return nil
}

func k8sNodepoolReady(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	subjectNodepool, err := client.GetKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id())

	if err != nil {
		return true, fmt.Errorf("Error checking k8s node pool status: %s", err)
	}
	return subjectNodepool.Metadata.State == "ACTIVE", nil
}

func k8sNodepoolDeleted(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	_, err := client.GetKubernetesNodePool(d.Get("k8s_cluster_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				return true, nil
			}
			return true, fmt.Errorf("Error checking k8s node pool deletion status: %s", err)
		}
	}
	return false, nil
}
