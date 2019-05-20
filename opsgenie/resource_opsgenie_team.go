package opsgenie

import (
	"context"
	"fmt"
	"log"
	"strings"

	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
	"github.com/opsgenie/opsgenie-go-sdk-v2/user"
)

func resourceOpsGenieTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpsGenieTeamCreate,
		Read:   resourceOpsGenieTeamRead,
		Update: resourceOpsGenieTeamUpdate,
		Delete: resourceOpsGenieTeamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOpsGenieTeamName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},

						"role": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "user",
							ValidateFunc: validateOpsGenieTeamRole,
						},
					},
				},
			},
		},
	}
}

func resourceOpsGenieTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).team
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	members, err := expandOpsGenieTeamMembers(d, meta)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating OpsGenie team '%s'", name)
	result, err := client.Create(context.Background(), &team.CreateTeamRequest{
		Name:        name,
		Description: description,
		Members:     members,
	})
	if err != nil {
		return err
	}

	d.SetId(result.Id)
	d.Set("name", result.Name)
	return resourceOpsGenieTeamRead(d, meta)
}

func resourceOpsGenieTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).team

	result, err := client.Get(context.Background(), &team.GetTeamRequest{
		IdentifierValue: d.Get("name").(string),
	})
	if err != nil {
		return err
	}

	d.Set("name", result.Name)
	d.Set("description", result.Description)
	d.Set("member", flattenOpsGenieTeamMembers(result.Members))

	return nil
}

func resourceOpsGenieTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).team
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	log.Printf("[INFO] Updating OpsGenie team '%s'", name)

	members, err := expandOpsGenieTeamMembers(d, meta)
	if err != nil {
		return err
	}
	result, err := client.Update(context.Background(), &team.UpdateTeamRequest{
		Id:          d.Id(),
		Name:        name,
		Description: description,
		Members:     members,
	})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Team %s updated", result.Name)
	return nil
}

func resourceOpsGenieTeamDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting OpsGenie team '%s'", d.Get("name").(string))
	client := meta.(*OpsGenieClient).team

	_, err := client.Delete(context.Background(), &team.DeleteTeamRequest{
		IdentifierValue: d.Get("name").(string),
	})
	if err != nil {
		return err
	}
	// log.Printf("[INFO] Team %s deleted", result.Name)
	return nil
}

func expandOpsGenieTeamMembers(d *schema.ResourceData, meta interface{}) ([]team.Member, error) {
	client := meta.(*OpsGenieClient).user

	input := d.Get("member").([]interface{})
	members := make([]team.Member, 0, len(input))
	if input == nil {
		return members, nil
	}

	for _, v := range input {
		config := v.(map[string]interface{})

		username := config["username"].(string)
		role := config["role"].(string)
		result, err := client.List(context.Background(), &user.ListRequest{})
		if err != nil {
			return nil, err
		}
		var user *user.User
		if len(result.Users) > 0 {
			for _, u := range result.Users {
				if u.Username == username {
					user = &u
					break
				}
			}
		}
		member := team.Member{
			User: team.User{
				ID:       user.Id,
				Username: user.Username,
			},
			Role: role,
		}

		members = append(members, member)
	}

	return members, nil
}

func validateOpsGenieTeamName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only alpha numeric characters and underscores are allowed in %q: %q", k, value))
	}

	if len(value) >= 100 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 100 characters: %q %d", k, value, len(value)))
	}

	return
}

func validateOpsGenieTeamRole(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"admin": true,
		"user":  true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("OpsGenie Team Role can only be 'Admin' or 'User'"))
	}

	return
}
