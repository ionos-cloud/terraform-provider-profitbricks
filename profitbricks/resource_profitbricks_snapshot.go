package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
	"time"
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
	}
}

func resourceProfitBricksSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	var err error
	dcId := d.Get("datacenter_id").(string)
	volumeId := d.Get("volume_id").(string)
	name := d.Get("name").(string)

	snapshot := profitbricks.CreateSnapshot(dcId, volumeId, name, "")

	if snapshot.StatusCode > 299 {
		return fmt.Errorf("An error occured while creating a snapshot: %s", snapshot.Response)
	}

	err = waitTillProvisioned(meta, snapshot.Headers.Get("Location"))
	if err != nil {
		return err
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

	err := waitTillProvisioned(meta, snapshot.Headers.Get("Location"))
	if err != nil {
		return err
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
	err := waitTillProvisioned(meta, resp.Headers.Get("Location"))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
