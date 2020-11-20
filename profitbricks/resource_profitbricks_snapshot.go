package profitbricks

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	profitbricks "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourceProfitBricksSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksSnapshotCreate,
		Read:   resourceProfitBricksSnapshotRead,
		Update: resourceProfitBricksSnapshotUpdate,
		Delete: resourceProfitBricksSnapshotDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceProfitBricksSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)
	volumeId := d.Get("volume_id").(string)
	name := d.Get("name").(string)

	snapshot, err := client.CreateSnapshot(dcId, volumeId, name, "")

	if err != nil {
		return fmt.Errorf("An error occured while creating a snapshot: %s", err)
	}

	d.SetId(snapshot.ID)
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, snapshot.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceProfitBricksSnapshotRead(d, meta)
}

func resourceProfitBricksSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	snapshot, err := client.GetSnapshot(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error occured while fetching a snapshot ID %s %s", d.Id(), err)
	}

	d.Set("name", snapshot.Properties.Name)
	return nil
}

func resourceProfitBricksSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	dcId := d.Get("datacenter_id").(string)
	volumeId := d.Get("volume_id").(string)

	snapshot, err := client.RestoreSnapshot(dcId, volumeId, d.Id())
	if err != nil {
		return fmt.Errorf("An error occured while restoring a snapshot ID %s %d", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, snapshot.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceProfitBricksSnapshotRead(d, meta)
}

func resourceProfitBricksSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	status, err := client.GetSnapshot(d.Id())
	if err != nil {
		return fmt.Errorf("An error occured while fetching a snapshot ID %s %s", d.Id(), err)
	}
	for status.Metadata.State != "AVAILABLE" {
		time.Sleep(30 * time.Second)
		status, err = client.GetSnapshot(d.Id())

		if err != nil {
			return fmt.Errorf("An error occured while fetching a snapshot ID %s %s", d.Id(), err)
		}
	}

	dcId := d.Get("datacenter_id").(string)
	dc, err := client.GetDatacenter(dcId)

	if err != nil {
		return fmt.Errorf("An error occured while fetching a Datacenter ID %s %s", dcId, err)
	}

	for dc.Metadata.State != "AVAILABLE" {
		time.Sleep(30 * time.Second)
		dc, err = client.GetDatacenter(dcId)

		if err != nil {
			return fmt.Errorf("An error occured while fetching a Datacenter ID %s %s", dcId, err)
		}
	}

	resp, err := client.DeleteSnapshot(d.Id())
	if err != nil {
		return fmt.Errorf("An error occured while deleting a snapshot ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
