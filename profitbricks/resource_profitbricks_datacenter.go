package profitbricks

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
)

func resourceProfitBricksDatacenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksDatacenterCreate,
		Read:   resourceProfitBricksDatacenterRead,
		Update: resourceProfitBricksDatacenterUpdate,
		Delete: resourceProfitBricksDatacenterDelete,
		Schema: map[string]*schema.Schema{

			//Datacenter parameters
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},

		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceProfitBricksDatacenterCreate(d *schema.ResourceData, meta interface{}) error {
	datacenter := profitbricks.Datacenter{
		Properties: profitbricks.DatacenterProperties{
			Name:     d.Get("name").(string),
			Location: d.Get("location").(string),
		},
	}

	if attr, ok := d.GetOk("description"); ok {
		datacenter.Properties.Description = attr.(string)
	}
	dc := profitbricks.CreateDatacenter(datacenter)

	if dc.StatusCode > 299 {
		return fmt.Errorf(
			"Error creating data center (%s) (%s)", d.Id(), dc.Response)
	}
	d.SetId(dc.Id)

	log.Printf("[INFO] DataCenter Id: %s", d.Id())

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, dc.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceProfitBricksDatacenterRead(d, meta)
}

func resourceProfitBricksDatacenterRead(d *schema.ResourceData, meta interface{}) error {
	datacenter := profitbricks.GetDatacenter(d.Id())
	if datacenter.StatusCode > 299 {
		if datacenter.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error while fetching a data center ID %s %s", d.Id(), datacenter.Response)
	}

	d.Set("name", datacenter.Properties.Name)
	d.Set("location", datacenter.Properties.Location)
	d.Set("description", datacenter.Properties.Description)
	return nil
}

func resourceProfitBricksDatacenterUpdate(d *schema.ResourceData, meta interface{}) error {
	obj := profitbricks.DatacenterProperties{}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")

		obj.Name = newName.(string)
	}

	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")
		obj.Description = newDescription.(string)
	}

	if d.HasChange("location") {
		oldLocation, newLocation := d.GetChange("location")
		return fmt.Errorf("Data center is created in %s location. You can not change location of the data center to %s. It requires recreation of the data center.", oldLocation, newLocation)
	}

	resp := profitbricks.PatchDatacenter(d.Id(), obj)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceProfitBricksDatacenterRead(d, meta)
}

func resourceProfitBricksDatacenterDelete(d *schema.ResourceData, meta interface{}) error {
	dcid := d.Id()
	resp := profitbricks.DeleteDatacenter(dcid)

	if resp.StatusCode > 299 {
		return fmt.Errorf("An error occured while deleting the data center ID %s %s", d.Id(), string(resp.Body))
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Headers.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}

func getImageId(dcId string, imageName string, imageType string) string {
	if imageName == "" {
		return ""
	}
	dc := profitbricks.GetDatacenter(dcId)
	if dc.StatusCode > 299 {
		log.Print(fmt.Errorf("Error while fetching a data center ID %s %s", dcId, dc.Response))
	}

	images := profitbricks.ListImages()
	if images.StatusCode > 299 {
		log.Print(fmt.Errorf("Error while fetching the list of images %s", images.Response))
	}

	if len(images.Items) > 0 {
		for _, i := range images.Items {
			imgName := ""
			if i.Properties.Name != "" {
				imgName = i.Properties.Name
			}

			if imageType == "SSD" {
				imageType = "HDD"
			}
			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(imageName)) && i.Properties.ImageType == imageType && i.Properties.Location == dc.Properties.Location && i.Properties.Public == true {
				return i.Id
			}
		}
	}
	return ""
}

func getSnapshotId(snapshotName string) string {
	if snapshotName == "" {
		return ""
	}
	snapshots := profitbricks.ListSnapshots()
	if snapshots.StatusCode > 299 {
		log.Print(fmt.Errorf("Error while fetching the list of snapshots %s", snapshots.Response))
	}

	if len(snapshots.Items) > 0 {
		for _, i := range snapshots.Items {
			imgName := ""
			if i.Properties.Name != "" {
				imgName = i.Properties.Name
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(snapshotName)) {
				return i.Id
			}
		}
	}
	return ""
}

func getImageAlias(imageAlias string, location string) string {
	if imageAlias == "" {
		return ""
	}
	locations := profitbricks.GetLocation(location)
	if locations.StatusCode > 299 {
		log.Print(fmt.Errorf("Error while fetching the list of snapshots %s", locations.Response))
	}

	if len(locations.Properties.ImageAliases) > 0 {
		for _, i := range locations.Properties.ImageAliases {
			alias := ""
			if i != "" {
				alias = i
			}

			if alias != "" && strings.ToLower(alias) == strings.ToLower(imageAlias) {
				return i
			}
		}
	}
	return ""
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
