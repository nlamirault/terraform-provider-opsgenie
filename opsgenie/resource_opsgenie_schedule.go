package opsgenie

import (
	"context"
	"fmt"
	"log"
	"time"

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
			"rotation": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end_date": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"participant": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"username": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// "name": "First Rotation",
//             "startDate": "2017-02-06T05:00:00Z",
//             "endDate": "2017-02-23T06:00:00Z",
//             "type": "hourly",
//             "length": 6,
//     "participants": [
//         {
//             "type": "team",
//             "id": "b3578948-55b3-4acc-9bf1-2ce2db3alpa2"
//         },
//         {
//             "type": "user",
//             "username": "user@opsgenie.com"
//         },
//         {
//             "type": "none"
//         }
//     ]
// }

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

	rotations, err := expandOpsGenieScheduleRotations(d, meta)
	if err != nil {
		return err
	}
	result, err := client.Create(context.Background(), &schedule.CreateRequest{
		Name: name,
		OwnerTeam: &og.OwnerTeam{
			Name: owner,
		},
		Description: description,
		Timezone:    timeZone,
		Rotations:   rotations,
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
	rotations, err := expandOpsGenieScheduleRotations(d, meta)
	if err != nil {
		return err
	}
	result, err := client.Update(context.Background(), &schedule.UpdateRequest{
		IdentifierValue: d.Id(),
		IdentifierType:  schedule.Id,
		Name:            name,
		OwnerTeam: &og.OwnerTeam{
			Name: owner,
		},
		Description: description,
		Timezone:    timeZone,
		Rotations:   rotations,
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

func expandOpsGenieScheduleRotations(d *schema.ResourceData, meta interface{}) ([]og.Rotation, error) {
	input := d.Get("rotation").([]interface{})
	rotations := make([]og.Rotation, 0, len(input))
	if input == nil {
		return rotations, nil
	}

	for _, v := range input {
		config := v.(map[string]interface{})

		name := config["name"].(string)
		startDate, err := time.Parse(time.RFC3339, config["start_date"].(string))
		if err != nil {
			return nil, err
		}
		endDate, err := time.Parse(time.RFC3339, config["end_date"].(string))
		if err != nil {
			return nil, err
		}
		rotationType, err := extractRotationType(config["type"].(string))
		if err != nil {
			return nil, err
		}

		input := config["participant"].([]interface{})
		participants := make([]og.Participant, 0, len(input))
		if input != nil {
			for _, elem := range input {
				participantConfig := elem.(map[string]interface{})
				usernameValue := participantConfig["username"].(string)
				idValue := participantConfig["id"].(string)
				participantType, err := extractParticipantType(participantConfig["type"].(string))
				if err != nil {
					return nil, err
				}
				participants = append(participants, og.Participant{
					Type:     *participantType,
					Username: usernameValue,
					Id:       idValue,
				})
			}
		}

		rotations = append(rotations, og.Rotation{
			Name:         name,
			StartDate:    &startDate,
			EndDate:      &endDate,
			Type:         *rotationType,
			Participants: participants,
		})
	}

	return rotations, nil
}

func extractRotationType(value string) (*og.RotationType, error) {
	var rotationType og.RotationType
	switch value {
	case "daily":
		rotationType = og.Daily
	case "weekly":
		rotationType = og.Weekly
	case "hourly":
		rotationType = og.Hourly
	default:
		return nil, fmt.Errorf("Invalid rotation type: %s", value)
	}
	return &rotationType, nil
}

func extractParticipantType(typeValue string) (*og.ParticipantType, error) {
	var participantType og.ParticipantType
	switch typeValue {
	case "user":
		participantType = og.User
	case "team":
		participantType = og.Team
	case "escalation":
		participantType = og.Escalation
	case "schedule":
		participantType = og.Schedule
	default:
		return nil, fmt.Errorf("Invalid participation type: %s", typeValue)
	}
	return &participantType, nil
}
