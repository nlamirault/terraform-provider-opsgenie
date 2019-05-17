package opsgenie

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfigForRegion(region string) (interface{}, error) {
	if os.Getenv("OPSGENIE_API_KEY") == "" {
		return nil, fmt.Errorf("OPSGENIE_API_KEY must be set")
	}

	config := Config{
		APIKey: os.Getenv("OPSGENIE_API_KEY"),
	}

	client, err := config.Client()
	if err != nil {
		return nil, fmt.Errorf("error getting OpsGenie client")
	}

	return client, nil
}
