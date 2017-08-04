package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
	"log"
	"time"
)

func resourceProfitBricksGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksGroupCreate,
		Read:   resourceProfitBricksGroupRead,
		Update: resourceProfitBricksGroupUpdate,
		Delete: resourceProfitBricksGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_datacenter": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reserve_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"access_activity_log": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"administrator": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"force_sec_auth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceProfitBricksGroupCreate(d *schema.ResourceData, meta interface{}) error {
	request := profitbricks.Group{
		Properties: profitbricks.GroupProperties{},
	}

	log.Printf("[DEBUG] NAME %s", d.Get("name"))
	if d.Get("name") != nil {
		request.Properties.Name = d.Get("name").(string)
	}

	tempCreateDataCenter := d.Get("create_datacenter").(bool)
	request.Properties.CreateDataCenter = &tempCreateDataCenter
	tempCreateSnapshot := d.Get("create_snapshot").(bool)
	request.Properties.CreateSnapshot = &tempCreateSnapshot
	tempReserveIp := d.Get("reserve_ip").(bool)
	request.Properties.ReserveIp = &tempReserveIp
	tempAccessActivityLog := d.Get("access_activity_log").(bool)
	request.Properties.AccessActivityLog = &tempAccessActivityLog

	usertoAdd := d.Get("user_id").(string)

	group := profitbricks.CreateGroup(request)

	log.Printf("[DEBUG] GROUP ID: %s", group.Id)

	if group.StatusCode > 299 {
		return fmt.Errorf("An error occured while creating a group: %s", group.Response)
	}

	err := waitTillProvisioned(meta, group.Headers.Get("Location"))
	if err != nil {
		return err
	}
	d.SetId(group.Id)

	//add users to group if any is provided
	if usertoAdd != "" {
		addedUser := profitbricks.AddUserToGroup(d.Id(), usertoAdd)
		if addedUser.StatusCode > 299 {
			return fmt.Errorf("An error occured while adding %s user to group ID %s %s", usertoAdd, d.Id(), group.Response)
		}
		err := waitTillProvisioned(meta, addedUser.Headers.Get("Location"))
		if err != nil {
			return err
		}
	}
	return resourceProfitBricksGroupRead(d, meta)
}

func resourceProfitBricksGroupRead(d *schema.ResourceData, meta interface{}) error {
	group := profitbricks.GetGroup(d.Id())

	if group.StatusCode > 299 {
		if group.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("An error occured while fetching a Group ID %s %s", d.Id(), group.Response)
	}

	d.Set("name", group.Properties.Name)
	d.Set("create_datacenter", group.Properties.CreateDataCenter)
	d.Set("create_snapshot", group.Properties.CreateSnapshot)
	d.Set("reserve_ip", group.Properties.ReserveIp)
	d.Set("access_activity_log", group.Properties.AccessActivityLog)

	users := profitbricks.ListGroupUsers(d.Id())
	var usersArray = []profitbricks.UserProperties{}
	if len(users.Items) > 0 {
		for _, usr := range users.Items {
			usersArray = append(usersArray, *usr.Properties)
		}
		d.Set("users", usersArray)
	}

	return nil
}

func resourceProfitBricksGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	tempCreateDataCenter := d.Get("create_datacenter").(bool)
	tempCreateSnapshot := d.Get("create_snapshot").(bool)
	tempReserveIp := d.Get("reserve_ip").(bool)
	tempAccessActivityLog := d.Get("access_activity_log").(bool)
	usertoAdd := d.Get("user_id").(string)
	groupReq := profitbricks.Group{
		Properties: profitbricks.GroupProperties{
			CreateDataCenter:  &tempCreateDataCenter,
			CreateSnapshot:    &tempCreateSnapshot,
			ReserveIp:         &tempReserveIp,
			AccessActivityLog: &tempAccessActivityLog,
		},
	}

	_, newValue := d.GetChange("name")
	groupReq.Properties.Name = newValue.(string)

	group := profitbricks.UpdateGroup(d.Id(), groupReq)
	if group.StatusCode > 299 {
		return fmt.Errorf("An error occured while patching a group ID %s %s", d.Id(), group.Response)
	}
	err := waitTillProvisioned(meta, group.Headers.Get("Location"))
	if err != nil {
		return err
	}
	//add users to group if any is provided
	if usertoAdd != "" {
		addedUser := profitbricks.AddUserToGroup(d.Id(), usertoAdd)
		if addedUser.StatusCode > 299 {
			return fmt.Errorf("An error occured while adding %s user to group ID %s %s", usertoAdd, d.Id(), group.Response)
		}
		err := waitTillProvisioned(meta, addedUser.Headers.Get("Location"))
		if err != nil {
			return err
		}
	}
	return resourceProfitBricksGroupRead(d, meta)
}

func resourceProfitBricksGroupDelete(d *schema.ResourceData, meta interface{}) error {
	resp := profitbricks.DeleteGroup(d.Id())
	if resp.StatusCode > 299 {
		//try again in 20 seconds
		time.Sleep(20 * time.Second)
		resp = profitbricks.DeleteGroup(d.Id())
		if resp.StatusCode > 299 && resp.StatusCode != 404 {
			return fmt.Errorf("An error occured while deleting a group %s %s", d.Id(), string(resp.Body))
		}
	}

	if resp.Headers.Get("Location") != "" {
		err := waitTillProvisioned(meta, resp.Headers.Get("Location"))
		if err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}
