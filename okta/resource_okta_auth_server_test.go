package okta

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func authServerExists(id string) (bool, error) {
	server, resp, err := getOktaClientFromMetadata(testAccProvider.Meta()).AuthorizationServer.GetAuthorizationServer(context.Background(), id)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return server != nil && server.Id != "" && err == nil, err
}

func TestAccOktaAuthServer_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.sun_also_rises", authServer)
	name := buildResourceName(ri)
	mgr := newFixtureManager(authServer)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "The best way to find out if you can trust somebody is to trust them."),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "AUTO"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "The past is not dead. In fact, it's not even past."),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "AUTO"),
				),
			},
		},
	})
}

func TestAccOktaAuthServer_fullStack(t *testing.T) {
	ri := acctest.RandInt()
	name := buildResourceName(ri)
	resourceName := fmt.Sprintf("%s.test", authServer)
	claimName := fmt.Sprintf("%s.test", authServerClaim)
	ruleName := fmt.Sprintf("%s.test", authServerPolicyRule)
	policyName := fmt.Sprintf("%s.test", authServerPolicy)
	scopeName := fmt.Sprintf("%s.test", authServerScope)
	mgr := newFixtureManager(authServer)
	config := mgr.GetFixtures("full_stack.tf", ri, t)
	updatedConfig := mgr.GetFixtures("full_stack_with_client.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "AUTO"),

					resource.TestCheckResourceAttr(scopeName, "name", "test:something"),
					resource.TestCheckResourceAttr(claimName, "name", "test"),
					resource.TestCheckResourceAttr(policyName, "name", "test"),
					resource.TestCheckResourceAttr(ruleName, "name", "test"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "test_updated"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "AUTO"),

					resource.TestCheckResourceAttr(scopeName, "name", "test:something"),
					resource.TestCheckResourceAttr(claimName, "name", "test"),
					resource.TestCheckResourceAttr(policyName, "name", "test"),
					resource.TestCheckResourceAttr(policyName, "client_whitelist.#", "1"),
					resource.TestCheckResourceAttr(ruleName, "name", "test"),
				),
			},
		},
	})
}

func TestAccOktaAuthServer_gh299(t *testing.T) {
	ri := acctest.RandInt()
	name := buildResourceName(ri)
	resourceName := fmt.Sprintf("%s.test", authServer)
	resource2Name := fmt.Sprintf("%s.test1", authServer)
	mgr := newFixtureManager(authServer)
	config := mgr.GetFixtures("dependency.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "AUTO"),

					resource.TestCheckResourceAttr(resource2Name, "name", name+"1"),
					resource.TestCheckResourceAttr(resource2Name, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resource2Name, "credentials_rotation_mode", "MANUAL"),
				),
			},
		},
	})
}
