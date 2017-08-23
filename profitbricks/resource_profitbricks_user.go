package profitbricks

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
	"log"
	"time"
)

func resourceProfitBricksUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfitBricksUserCreate,
		Read:   resourceProfitBricksUserRead,
		Update: resourceProfitBricksUserUpdate,
		Delete: resourceProfitBricksUserDelete,
		Schema: map[string]*schema.Schema{
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"administrator": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"force_sec_auth": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceProfitBricksUserCreate(d *schema.ResourceData, meta interface{}) error {
	request := profitbricks.User{
		Properties: &profitbricks.UserProperties{},
	}

	log.Printf("[DEBUG] NAME %s", d.Get("first_name"))

	if d.Get("first_name") != nil {
		request.Properties.Firstname = d.Get("first_name").(string)
	}
	if d.Get("last_name") != nil {
		request.Properties.Lastname = d.Get("last_name").(string)
	}
	if d.Get("email") != nil {
		request.Properties.Email = d.Get("email").(string)
	}
	if d.Get("password") != nil {
		request.Properties.Password = d.Get("password").(string)
	}

	request.Properties.Administrator = d.Get("administrator").(bool)
	request.Properties.ForceSecAuth = d.Get("force_sec_auth").(bool)
	user := profitbricks.CreateUser(request)

	log.Printf("[DEBUG] USER ID: %s", user.Id)

	if user.StatusCode > 299 {
		return fmt.Errorf("An error occured while creating a user: %s", user.Response)
	}

	err := waitTillProvisioned(meta, user.Headers.Get("Location"))
	if err != nil {
		return err
	}
	d.SetId(user.Id)
	return resourceProfitBricksUserRead(d, meta)
}

func resourceProfitBricksUserRead(d *schema.ResourceData, meta interface{}) error {
	user := profitbricks.GetUser(d.Id())

	if user.StatusCode > 299 {
		if user.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("An error occured while fetching a User ID %s %s", d.Id(), user.Response)
	}

	d.Set("first_name", user.Properties.Firstname)
	d.Set("last_name", user.Properties.Lastname)
	d.Set("email", user.Properties.Email)
	d.Set("administrator", user.Properties.Administrator)
	d.Set("force_sec_auth", user.Properties.ForceSecAuth)
	return nil
}

func resourceProfitBricksUserUpdate(d *schema.ResourceData, meta interface{}) error {
	originalUser := profitbricks.GetUser(d.Id())
	userReq := profitbricks.User{
		Properties: &profitbricks.UserProperties{
			Administrator: d.Get("administrator").(bool),
			ForceSecAuth:  d.Get("force_sec_auth").(bool),
		},
	}

	if d.HasChange("first_name") {
		_, newValue := d.GetChange("first_name")
		userReq.Properties.Firstname = newValue.(string)

	} else {
		userReq.Properties.Firstname = originalUser.Properties.Firstname
	}

	if d.HasChange("last_name") {
		_, newValue := d.GetChange("last_name")
		userReq.Properties.Lastname = newValue.(string)
	} else {
		userReq.Properties.Lastname = originalUser.Properties.Lastname
	}

	if d.HasChange("email") {
		_, newValue := d.GetChange("email")
		userReq.Properties.Email = newValue.(string)
	} else {
		userReq.Properties.Email = originalUser.Properties.Email
	}

	user := profitbricks.UpdateUser(d.Id(), userReq)
	if user.StatusCode > 299 {
		return fmt.Errorf("An error occured while patching a user ID %s %s", d.Id(), user.Response)
	}
	err := waitTillProvisioned(meta, user.Headers.Get("Location"))
	if err != nil {
		return err
	}
	return resourceProfitBricksUserRead(d, meta)
}

func resourceProfitBricksUserDelete(d *schema.ResourceData, meta interface{}) error {
	resp := profitbricks.DeleteUser(d.Id())
	if resp.StatusCode > 299 {
		//try again in 20 seconds
		time.Sleep(20 * time.Second)
		resp = profitbricks.DeleteUser(d.Id())
		if resp.StatusCode > 299 && resp.StatusCode != 404 {
			return fmt.Errorf("An error occured while deleting a user %s %s", d.Id(), string(resp.Body))
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
