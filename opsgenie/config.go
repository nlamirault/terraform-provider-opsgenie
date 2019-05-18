package opsgenie

import (
	"log"

	"golang.org/x/net/context"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/contact"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
	"github.com/opsgenie/opsgenie-go-sdk-v2/user"
)

type OpsGenieClient struct {
	apiKey string

	StopContext context.Context

	team     team.Client
	user     user.Client
	contact  contact.Client
	schedule schedule.Client
}

// Config defines the configuration options for the OpsGenie client
type Config struct {
	APIKey string
}

// Client returns a new OpsGenie client
func (c *Config) Client() (*OpsGenieClient, error) {
	opsGenieConfig := &client.Config{
		ApiKey:         c.APIKey,
		OpsGenieAPIURL: client.API_URL,
	}
	client := OpsGenieClient{}

	log.Printf("[INFO] OpsGenie client configured")

	teamClient, err := team.NewClient(opsGenieConfig)
	if err != nil {
		return nil, err
	}
	client.team = *teamClient

	userClient, err := user.NewClient(opsGenieConfig)
	if err != nil {
		return nil, err
	}
	client.user = *userClient

	contactClient, err := contact.NewClient(opsGenieConfig)
	if err != nil {
		return nil, err
	}
	client.contact = *contactClient

	scheduleClient, err := schedule.NewClient(opsGenieConfig)
	if err != nil {
		return nil, err
	}
	client.schedule = *scheduleClient

	return &client, nil
}
