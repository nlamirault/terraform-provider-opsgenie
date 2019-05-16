package opsgenie

import (
	"context"
	// "fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	user "github.com/opsgenie/opsgenie-go-sdk-v2/user"
)

func dataSourceOpsGenieUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpsGenieUserRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOpsGenieUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).user

	username := d.Get("username").(string)

	log.Printf("[INFO] Reading OpsGenie user '%s'", username)

	// result, err := client.List(context.Background(), &user.ListRequest{})
	// if err != nil {
	// 	return nil
	// }

	// var found *user.User
	// if len(result.Users) > 0 {
	// 	for _, user := range result.Users {
	// 		if user.Username == username {
	// 			found = &user
	// 			break
	// 		}
	// 	}
	// }
	// if found == nil {
	// 	return fmt.Errorf("Unable to locate any user with the username: %s", username)
	// }

	result, err := client.Get(context.Background(), &user.GetRequest{
		Identifier: d.Id(),
	})
	if err != nil {
		return err
	}

	d.SetId(result.Id)
	d.Set("username", result.Username)
	d.Set("full_name", result.FullName)
	d.Set("role", result.Role)

	return nil
}
