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
	connection := meta.(*profitbricks.Client)
	datacenter := profitbricks.Datacenter{
		Properties: profitbricks.DatacenterProperties{
			Name:     d.Get("name").(string),
			Location: d.Get("location").(string),
		},
	}

	if attr, ok := d.GetOk("description"); ok {
		datacenter.Properties.Description = attr.(string)
	}
	dc, err := connection.CreateDatacenter(datacenter)

	if err != nil {
		return fmt.Errorf(
			"Error creating data center (%s) (%s)", d.Id(), err)
	}
	d.SetId(dc.ID)

	log.Printf("[INFO] DataCenter Id: %s", d.Id())

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, dc.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceProfitBricksDatacenterRead(d, meta)
}

func resourceProfitBricksDatacenterRead(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	datacenter, err := connection.GetDatacenter(d.Id())

	if err != nil {
		if err2, ok := err.(profitbricks.ApiError); ok {
			if err2.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error while fetching a data center ID %s %s", d.Id(), err)
	}

	d.Set("name", datacenter.Properties.Name)
	d.Set("location", datacenter.Properties.Location)
	d.Set("description", datacenter.Properties.Description)
	return nil
}

func resourceProfitBricksDatacenterUpdate(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
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

	resp, err := connection.UpdateDataCenter(d.Id(), obj)

	if err != nil {
		return fmt.Errorf("An error occured while update the data center ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceProfitBricksDatacenterRead(d, meta)
}

func resourceProfitBricksDatacenterDelete(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	dcid := d.Id()
	resp, err := connection.DeleteDatacenter(dcid)

	if err != nil {
		return fmt.Errorf("An error occured while deleting the data center ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}

func getImageId(connection *profitbricks.Client, dcId string, imageName string, imageType string) string {
	if imageName == "" {
		return ""
	}
	dc, err := connection.GetDatacenter(dcId)
	if err != nil {
		log.Print(fmt.Errorf("Error while fetching a data center ID %s %s", dcId, err))
	}

	images, err := connection.ListImages()
	if err != nil {
		log.Print(fmt.Errorf("Error while fetching the list of images %s", err))
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
				return i.ID
			}
		}
	}
	return ""
}

func getSnapshotId(connection *profitbricks.Client, snapshotName string) string {
	if snapshotName == "" {
		return ""
	}
	snapshots, err := connection.ListSnapshots()
	if err != nil {
		log.Print(fmt.Errorf("Error while fetching the list of snapshots %s", err))
	}

	if len(snapshots.Items) > 0 {
		for _, i := range snapshots.Items {
			imgName := ""
			if i.Properties.Name != "" {
				imgName = i.Properties.Name
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(snapshotName)) {
				return i.ID
			}
		}
	}
	return ""
}

func getImageAlias(connection *profitbricks.Client, imageAlias string, location string) string {
	if imageAlias == "" {
		return ""
	}
	locations, err := connection.GetLocation(location)
	if err != nil {
		log.Print(fmt.Errorf("Error while fetching the list of snapshots %s", err))
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
