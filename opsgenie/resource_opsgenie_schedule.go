package opsgenie

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
)

func resourceOpsGenieSchedule() *schema.Resource {
	return &schema.Resource{
		Read:   resourceOpsGenieScheduleRead,
		Create: resourceOpsGenieScheduleCreate,
		Update: resourceOpsGenieScheduleUpdate,
		Delete: resourceOpsGenieScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "America/New_York",
			},
		},
	}
}

func resourceOpsGenieScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).schedule

	result, err := client.Get(context.Background(), &schedule.GetRequest{
		IdentifierValue: d.Get("name").(string),
	})
	if err != nil {
		return err
	}

	d.Set("name", result.Schedule.Name)
	d.Set("timezone", result.Schedule.Timezone)
	d.Set("description", result.Schedule.Description)
	d.Set("owner", result.Schedule.OwnerTeam)

	return nil
}

func resourceOpsGenieScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).schedule

	name := d.Get("name").(string)
	owner := d.Get("owner").(string)
	description := d.Get("description").(string)
	timeZone := d.Get("timezone").(string)

	result, err := client.Create(context.Background(), &schedule.CreateRequest{
		Name: name,
		OwnerTeam: &og.OwnerTeam{
			Name: owner,
		},
		Description: description,
		Timezone:    timeZone,
	})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created schedule: %v", result)

	d.SetId(result.Id)
	return resourceOpsGenieScheduleRead(d, meta)
}

func resourceOpsGenieScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).schedule

	name := d.Get("name").(string)
	owner := d.Get("owner").(string)
	description := d.Get("description").(string)
	timeZone := d.Get("timezone").(string)

	result, err := client.Update(context.Background(), &schedule.UpdateRequest{
		IdentifierValue: d.Id(),
		IdentifierType:  schedule.Id,
		Name:            name,
		OwnerTeam: &og.OwnerTeam{
			Name: owner,
		},
		Description: description,
		Timezone:    timeZone,
	})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Schedule %s updated", result.Name)

	return nil
}

func resourceOpsGenieScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).schedule

	result, err := client.Delete(context.Background(), &schedule.DeleteRequest{
		IdentifierValue: d.Get("name").(string),
	})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Deleted schedule: %v", result)
	return nil
}
