package opsgenie

import (
	"log"

	"golang.org/x/net/context"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
	"github.com/opsgenie/opsgenie-go-sdk-v2/user"
)

type OpsGenieClient struct {
	apiKey string

	StopContext context.Context

	team team.Client
	user user.Client
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
	// opsGenie, err := client.NewOpsGenieClient(&client.Config{
	// 	ApiKey:         c.ApiKey,
	// 	OpsGenieAPIURL: client.API_URL,
	// })
	// opsGenie := new(client.OpsGenieClient)
	// opsGenie.SetAPIKey(c.ApiKey)
	client := OpsGenieClient{}

	log.Printf("[INFO] OpsGenie client configured")

	// teamsClient, err := opsGenie.Team()
	teamClient, err := team.NewClient(opsGenieConfig)
	if err != nil {
		return nil, err
	}
	client.team = *teamClient

	// usersClient, err := opsGenie.UserV2()
	userClient, err := user.NewClient(opsGenieConfig)
	if err != nil {
		return nil, err
	}
	client.user = *userClient

	return &client, nil
}
