package profitbricks

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
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
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},

		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceProfitBricksSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	dcId := d.Get("datacenter_id").(string)
	volumeId := d.Get("volume_id").(string)
	name := d.Get("name").(string)

	snapshot := profitbricks.CreateSnapshot(dcId, volumeId, name, "")

	if snapshot.StatusCode > 299 {
		return fmt.Errorf("An error occured while creating a snapshot: %s", snapshot.Response)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, snapshot.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId(snapshot.Id)

	return resourceProfitBricksSnapshotRead(d, meta)
}

func resourceProfitBricksSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	snapshot := profitbricks.GetSnapshot(d.Id())

	if snapshot.StatusCode > 299 {
		if snapshot.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error occured while fetching a snapshot ID %s %s", d.Id(), snapshot.Response)
	}

	if snapshot.StatusCode > 299 {
		return fmt.Errorf("An error occured while fetching a snapshot ID %s %s", d.Id(), snapshot.Response)

	}

	d.Set("name", snapshot.Properties.Name)
	return nil
}

func resourceProfitBricksSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	dcId := d.Get("datacenter_id").(string)
	volumeId := d.Get("volume_id").(string)

	snapshot := profitbricks.RestoreSnapshot(dcId, volumeId, d.Id())
	if snapshot.StatusCode > 299 {
		return fmt.Errorf("An error occured while restoring a snapshot ID %s %d", d.Id(), snapshot.StatusCode)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, snapshot.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceProfitBricksSnapshotRead(d, meta)
}

func resourceProfitBricksSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	status := profitbricks.GetSnapshot(d.Id())
	for status.Metadata.State != "AVAILABLE" {
		time.Sleep(30 * time.Second)
		status = profitbricks.GetSnapshot(d.Id())
	}

	dcId := d.Get("datacenter_id").(string)
	dc := profitbricks.GetDatacenter(dcId)
	for dc.Metadata.State != "AVAILABLE" {
		time.Sleep(30 * time.Second)
		dc = profitbricks.GetDatacenter(dcId)
	}

	resp := profitbricks.DeleteSnapshot(d.Id())
	if resp.StatusCode > 299 {
		return fmt.Errorf("An error occured while deleting a snapshot ID %s %s", d.Id(), string(resp.Body))
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Headers.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
