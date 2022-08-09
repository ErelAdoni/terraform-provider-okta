package okta

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOktaAppSignOnPolicyRule(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", appSignOnPolicyRule)
	mgr := newFixtureManager(appSignOnPolicyRule, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      appSignOnPolicyRuleExists,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttrSet(resourceName, "priority"),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "factor_mode", "2FA"),
					resource.TestCheckResourceAttr(resourceName, "groups_excluded.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "groups_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "user_types_excluded.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "user_types_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "users_included.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_includes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_excludes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_connection", "ANYWHERE"),
					resource.TestCheckResourceAttr(resourceName, "constraints.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "re_authentication_frequency", "PT2H"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_period", "PT1H"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttrSet(resourceName, "priority"),
					resource.TestCheckResourceAttr(resourceName, "access", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "factor_mode", "2FA"),
					resource.TestCheckResourceAttr(resourceName, "groups_excluded.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "groups_included.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "user_types_excluded.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_types_included.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users_excluded.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "users_included.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_includes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_excludes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_connection", "ZONE"),
					resource.TestCheckResourceAttr(resourceName, "platform_include.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "re_authentication_frequency", "PT43800H"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_period", "PT2H"),
					resource.TestCheckResourceAttr(resourceName, "type", "ASSURANCE"),
					resource.TestCheckResourceAttr(resourceName, "constraints.#", "2"),
				),
			},
		},
	})
}

func appSignOnPolicyRuleExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != appSignOnPolicyRule {
			continue
		}
		client := apiSupplementForTest()
		rule, resp, err := client.GetAppSignOnPolicyRule(context.Background(), rs.Primary.Attributes["policy_id"], rs.Primary.ID)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		} else if err != nil {
			return err
		}
		if rule != nil {
			return fmt.Errorf("app sign-on policy rule still exists, ID: %s, PolicyID: %s", rs.Primary.ID, rs.Primary.Attributes["policy_id"])
		}
		return nil
	}
	return nil
}
