package opsgenie

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	user "github.com/opsgenie/opsgenie-go-sdk-v2/user"
)

func init() {
	resource.AddTestSweepers("opsgenie_user", &resource.Sweeper{
		Name: "opsgenie_user",
		F:    testSweepUser,
	})

}

func testSweepUser(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*OpsGenieClient).user

	resp, err := client.List(context.Background(), &user.ListRequest{})
	if err != nil {
		return err
	}

	for _, u := range resp.Users {
		if strings.HasPrefix(u.Username, "acctest-") {
			log.Printf("Destroying user %s", u.Username)
			if _, err := client.Delete(context.Background(), &user.DeleteRequest{
				Identifier: u.Id,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccOpsGenieUserUsername_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "hello",
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(99),
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(100),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateOpsGenieUserUsername(tc.Value, "opsgenie_team")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the OpsGenie User Username Validation to trigger a validation error: %v", errors)
		}
	}
}

func TestAccOpsGenieUserFullName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "hello",
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(100),
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(511),
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(512),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateOpsGenieUserFullName(tc.Value, "opsgenie_team")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the OpsGenie User Full Name Validation to trigger a validation error: %v", errors)
		}
	}
}

func TestAccOpsGenieUserRole_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "hello",
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(100),
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(511),
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(512),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateOpsGenieUserRole(tc.Value, "opsgenie_team")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the OpsGenie User Role Validation to trigger a validation error: %v", errors)
		}
	}
}

func TestAccOpsGenieUser_basic(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieUser_basic(rs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieUserExists("opsgenie_user.test"),
				),
			},
		},
	})
}

func TestAccOpsGenieUser_complete(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccOpsGenieUser_complete(rs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckOpsGenieUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckOpsGenieUserExists("opsgenie_user.test"),
				),
			},
		},
	})
}

func testCheckOpsGenieUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*OpsGenieClient).user

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opsgenie_user" {
			continue
		}

		result, _ := client.Get(context.Background(), &user.GetRequest{
			Identifier: rs.Primary.Attributes["id"],
		})
		if result != nil {
			return fmt.Errorf("User still exists:\n%#v", result)
		}
	}

	return nil
}

func testCheckOpsGenieUserExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.Attributes["id"]
		username := rs.Primary.Attributes["username"]

		client := testAccProvider.Meta().(*OpsGenieClient).user

		result, _ := client.Get(context.Background(), &user.GetRequest{
			Identifier: rs.Primary.Attributes["id"],
		})
		if result == nil {
			return fmt.Errorf("Bad: User %q (username: %q) does not exist", id, username)
		}

		return nil
	}
}

func testAccOpsGenieUser_basic(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "terraform-acctest+%s@hashicorp.com"
  full_name = "Acceptance Test User"
  role      = "User"
}
`, rString)
}

func testAccOpsGenieUser_complete(rString string) string {
	return fmt.Sprintf(`
resource "opsgenie_user" "test" {
  username  = "terraform-acctest+%s@hashicorp.com"
  full_name = "Acceptance Test User"
  role      = "User"
  locale    = "en_GB"
  timezone  = "Etc/GMT"
}
`, rString)
}
