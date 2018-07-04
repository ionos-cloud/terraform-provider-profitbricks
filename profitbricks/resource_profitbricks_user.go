package profitbricks

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go"
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

		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceProfitBricksUserCreate(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
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
	user, err := connection.CreateUser(request)

	log.Printf("[DEBUG] USER ID: %s", user.ID)

	if err != nil {
		return fmt.Errorf("An error occured while creating a user: %s", err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, user.Headers.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId(user.ID)
	return resourceProfitBricksUserRead(d, meta)
}

func resourceProfitBricksUserRead(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	user, err := connection.GetUser(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("An error occured while fetching a User ID %s %s", d.Id(), err)
	}

	d.Set("first_name", user.Properties.Firstname)
	d.Set("last_name", user.Properties.Lastname)
	d.Set("email", user.Properties.Email)
	d.Set("administrator", user.Properties.Administrator)
	d.Set("force_sec_auth", user.Properties.ForceSecAuth)
	return nil
}

func resourceProfitBricksUserUpdate(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	originalUser, _ := connection.GetUser(d.Id())
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

	user, err := connection.UpdateUser(d.Id(), userReq)
	if err != nil {
		return fmt.Errorf("An error occured while patching a user ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, user.Headers.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceProfitBricksUserRead(d, meta)
}

func resourceProfitBricksUserDelete(d *schema.ResourceData, meta interface{}) error {
	connection := meta.(*profitbricks.Client)
	resp, err := connection.DeleteUser(d.Id())
	if err != nil {
		//try again in 20 seconds
		time.Sleep(20 * time.Second)
		resp, err = connection.DeleteUser(d.Id())
		if err != nil {
			if apiError, ok := err.(profitbricks.ApiError); ok {
				if apiError.HttpStatusCode() != 404 {
					return fmt.Errorf("An error occured while deleting a user %s %s", d.Id(), err)
				}
			}
		}
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, resp.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
