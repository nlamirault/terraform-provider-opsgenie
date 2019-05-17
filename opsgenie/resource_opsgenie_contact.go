package opsgenie

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/opsgenie/opsgenie-go-sdk-v2/contact"
)

func resourceOpsGenieContact() *schema.Resource {
	return &schema.Resource{
		Read:   resourceOpsGenieContactRead,
		Create: resourceOpsGenieContactCreate,
		Delete: resourceOpsGenieContactDelete,
		Update: resourceOpsGenieContactUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"to": {
				Type:     schema.TypeString,
				Required: true,
			},
			"method": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceOpsGenieContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).contact

	result, err := client.Get(context.Background(), &contact.GetRequest{
		ContactIdentifier: d.Id(),
		UserIdentifier:    d.Get("username").(string),
	})
	if err != nil {
		return err
	}
	d.Set("to", result.To)
	d.Set("method", result.MethodOfContact)
	return nil
}

func resourceOpsGenieContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).contact

	username := d.Get("username").(string)
	to := d.Get("to").(string)
	method := d.Get("method").(string)

	log.Printf("[INFO] Creating OpsGenie contact for '%s' '%s' '%s'", username, to, method)
	var methodType contact.MethodType
	switch method {
	case "sms":
		methodType = contact.Sms
	case "email":
		methodType = contact.Email
	case "voice":
		methodType = contact.Voice
	default:
		return fmt.Errorf("Invalid method type: %s", method)
	}

	result, err := client.Create(context.Background(), &contact.CreateRequest{
		UserIdentifier:  username,
		To:              to,
		MethodOfContact: methodType,
	})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created contact %v", result)
	d.SetId(result.Id)

	return resourceOpsGenieContactRead(d, meta)
}

func resourceOpsGenieContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OpsGenieClient).contact
	log.Printf("[INFO] Deleting OpsGenie contact")

	result, err := client.Delete(context.Background(), &contact.DeleteRequest{
		ContactIdentifier: d.Id(),
		UserIdentifier:    d.Get("username").(string),
	})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted contact: %v", result)
	return nil
}

func resourceOpsGenieContactUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating OpsGenie contact")

	client := meta.(*OpsGenieClient).contact
	username := d.Get("username").(string)
	to := d.Get("to").(string)

	result, err := client.Update(context.Background(), &contact.UpdateRequest{
		ContactIdentifier: d.Id(),
		UserIdentifier:    username,
		To:                to,
	})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Updated contact: %v", result)
	return nil
}
