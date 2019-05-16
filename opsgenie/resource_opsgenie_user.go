package opsgenie

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/user"
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
	client := meta.(*OpsGenieClient).user

	username := d.Get("username").(string)
	fullName := d.Get("full_name").(string)
	role := d.Get("role").(string)
	locale := d.Get("locale").(string)
	timeZone := d.Get("timezone").(string)

	log.Printf("[INFO] Creating OpsGenie user '%s'", username)

	result, err := client.Create(context.Background(), &user.CreateRequest{
		Username: username,
		FullName: fullName,
		Role: &user.UserRoleRequest{
			RoleName: role,
		},
		Locale:   locale,
		TimeZone: timeZone,
	})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created user: %v", result)

	d.SetId(result.Id)
	return resourceOpsGenieUserRead(d, meta)
}

func resourceOpsGenieUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).user

	result, err := client.Get(context.Background(), &user.GetRequest{
		Identifier: d.Id(),
	})
	if err != nil {
		return err
	}

	d.Set("username", result.Username)
	d.Set("full_name", result.FullName)
	d.Set("role", result.Role)
	d.Set("locale", result.Locale)
	d.Set("timezone", result.TimeZone)

	return nil
}

func resourceOpsGenieUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).user

	username := d.Get("username").(string)
	fullName := d.Get("full_name").(string)
	role := d.Get("role").(string)
	locale := d.Get("locale").(string)
	timeZone := d.Get("timezone").(string)

	log.Printf("[INFO] Updating OpsGenie user '%s'", username)

	result, err := client.Update(context.Background(), &user.UpdateRequest{
		Identifier: d.Id(),
		FullName:   fullName,
		Role:       &user.UserRoleRequest{RoleName: role},
		Locale:     locale,
		TimeZone:   timeZone,
	})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Updated user: %v", result)

	// err = checkOpsGenieResponse(updateResponse.Code, updateResponse.Status)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func resourceOpsGenieUserDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie user '%s'", d.Get("username").(string))
	client := meta.(*OpsGenieClient).user

	result, err := client.Delete(context.Background(), &user.DeleteRequest{
		Identifier: d.Id(),
	})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted user: %v", result)
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
