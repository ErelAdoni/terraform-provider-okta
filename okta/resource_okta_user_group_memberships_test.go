package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaUserGroupMemberships_crud(t *testing.T) {
	mgr := newFixtureManager(userGroupMemberships, t.Name())
	start := mgr.GetFixtures("basic.tf", t)
	update := mgr.GetFixtures("basic_update.tf", t)
	remove := mgr.GetFixtures("basic_removal.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: start,
			},
			{
				Config: update,
			},
			{
				Config: remove,
			},
		},
	})
}
