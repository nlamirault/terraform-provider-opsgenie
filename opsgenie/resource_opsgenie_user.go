package opsgenie

import (
	"log"

	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk/userv2"
)

func resourceOpsGenieUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieUserCreate,
		Read:   resourceOpsGenieUserRead,
		Update: resourceOpsGenieUserUpdate,
		Delete: resourceOpsGenieUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateOpsGenieUserUsername,
			},
			"full_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieUserFullName,
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieUserRole,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en_US",
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "America/New_York",
			},
		},
	}
}

func resourceOpsGenieUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).users

	username := d.Get("username").(string)
	fullName := d.Get("full_name").(string)
	role := d.Get("role").(string)
	locale := d.Get("locale").(string)
	timeZone := d.Get("timezone").(string)

	createRequest := userv2.CreateUserRequest{
		UserName: username,
		FullName: fullName,
		Role:     &userv2.UserRole{Name: role},
		Locale:   locale,
		Timezone: timeZone,
	}

	log.Printf("[INFO] Creating OpsGenie user '%s' %v", username, createRequest)
	createResponse, err := client.Create(createRequest)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created user: %v", createResponse)
	// err = checkOpsGenieResponse(createResponse.Code, createResponse.Status)
	// if err != nil {
	// 	return err
	// }

	getRequest := userv2.GetUserRequest{
		Identifier: &userv2.Identifier{
			Username: username,
		},
	}

	getResponse, err := client.Get(getRequest)
	if err != nil {
		return err
	}

	d.SetId(getResponse.User.ID)

	return resourceOpsGenieUserRead(d, meta)
}

func resourceOpsGenieUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).users

	listRequest := userv2.ListUsersRequest{}
	listResponse, err := client.List(listRequest)
	if err != nil {
		return err
	}

	var found *userv2.User
	for _, user := range listResponse.Users {
		if user.ID == d.Id() {
			found = &user
			break
		}
	}

	if found == nil {
		d.SetId("")
		log.Printf("[INFO] User %q not found. Removing from state", d.Get("username").(string))
		return nil
	}

	getRequest := userv2.GetUserRequest{
		Identifier: &userv2.Identifier{
			Username: found.Username,
		},
	}

	getResponse, err := client.Get(getRequest)
	if err != nil {
		return err
	}

	d.Set("username", getResponse.User.Username)
	d.Set("full_name", getResponse.User.FullName)
	d.Set("role", getResponse.User.Role)
	d.Set("locale", getResponse.User.Locale)
	d.Set("timezone", getResponse.User.TimeZone)

	return nil
}

func resourceOpsGenieUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).users

	username := d.Get("username").(string)
	fullName := d.Get("full_name").(string)
	role := d.Get("role").(string)
	locale := d.Get("locale").(string)
	timeZone := d.Get("timezone").(string)

	log.Printf("[INFO] Updating OpsGenie user '%s'", username)

	updateRequest := userv2.UpdateUserRequest{
		Identifier: &userv2.Identifier{
			Username: username,
		},
		FullName: fullName,
		Role:     userv2.UserRole{Name: role},
		Locale:   locale,
		Timezone: timeZone,
	}

	updateResponse, err := client.Update(updateRequest)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Updated user: %v", updateResponse)

	// err = checkOpsGenieResponse(updateResponse.Code, updateResponse.Status)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func resourceOpsGenieUserDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie user '%s'", d.Get("username").(string))
	client := meta.(*OpsGenieClient).users

	deleteRequest := userv2.DeleteUserRequest{
		Identifier: &userv2.Identifier{
			ID: d.Id(),
		},
	}

	deleteResponse, err := client.Delete(deleteRequest)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Deleted user: %v", deleteResponse)
	return nil
}

func validateOpsGenieUserUsername(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) >= 100 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 100 characters: %q %d", k, value, len(value)))
	}

	return
}

func validateOpsGenieUserFullName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) >= 512 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 512 characters: %q %d", k, value, len(value)))
	}

	return
}

func validateOpsGenieUserRole(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) >= 512 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 512 characters: %q %d", k, value, len(value)))
	}

	return
}
