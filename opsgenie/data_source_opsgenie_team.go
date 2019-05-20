package opsgenie

import (
	"context"
	// "fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
)

func dataSourceOpsGenieTeam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpsGenieTeamRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceOpsGenieTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).team
	name := d.Get("name").(string)
	log.Printf("[INFO] Reading OpsGenie team '%s'", name)

	result, err := client.Get(context.Background(), &team.GetTeamRequest{
		IdentifierValue: name,
	})
	if err != nil {
		return err
	}
	d.SetId(result.Id)
	d.Set("name", result.Name)
	d.Set("description", result.Description)
	d.Set("member", flattenOpsGenieTeamMembers(result.Members))

	return nil
}
