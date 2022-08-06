package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOktaBrand_read(t *testing.T) {
	mgr := newFixtureManager(brand, t.Name())
	config := mgr.GetFixtures("datasource.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		Steps: []resource.TestStep{
			{
				Config:  config,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_brand.example", "id"),
					resource.TestCheckResourceAttrSet("data.okta_brand.example", "links"),
				),
			},
		},
	})
}
