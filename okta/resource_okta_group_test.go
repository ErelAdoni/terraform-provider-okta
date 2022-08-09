package okta

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaGroup_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", group)
	mgr := newFixtureManager(group, t.Name())
	config := mgr.GetFixtures("okta_group.tf", t)
	updatedConfig := mgr.GetFixtures("okta_group_updated.tf", t)
	addUsersConfig := mgr.GetFixtures("okta_group_with_users.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(group, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc")),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testAccDifferent")),
			},
			{
				Config: addUsersConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "4"),
				),
			},
		},
	})
}

func TestAccOktaGroup_customschema(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", group)
	mgr := newFixtureManager(group, t.Name())
	base := mgr.GetFixtures("okta_group_custom_base.tf", t)
	updated := mgr.GetFixtures("okta_group_custom_updated.tf", t)
	removal := mgr.GetFixtures("okta_group_custom_removal.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(group, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(mgr.Seed))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"testing1234\",\"testSchema2_%s\":true,\"testSchema3_%s\":54321}", strconv.Itoa(mgr.Seed), strconv.Itoa(mgr.Seed), strconv.Itoa(mgr.Seed))),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(mgr.Seed))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"moretesting1234\",\"testSchema2_%s\":false,\"testSchema3_%s\":12345}", strconv.Itoa(mgr.Seed), strconv.Itoa(mgr.Seed), strconv.Itoa(mgr.Seed))),
				),
			},
			{
				Config: removal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(mgr.Seed))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"moretesting1234\"}", strconv.Itoa(mgr.Seed))),
				),
			},
		},
	})
}

func doesGroupExist(id string) (bool, error) {
	client := oktaClientForTest()
	_, response, err := client.Group.GetGroup(context.Background(), id)
	return doesResourceExist(response, err)
}
