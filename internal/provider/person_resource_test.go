package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPersonResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccPersonResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apiserver_person.test", "name", "MyName"),
					resource.TestCheckResourceAttr("apiserver_person.test", "age", "42"),
					resource.TestCheckResourceAttr("apiserver_person.test", "description", "one"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "apiserver_person.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"id", "last_update"},
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccPersonResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apiserver_person.test", "name", "MyName"),
					resource.TestCheckResourceAttr("apiserver_person.test", "description", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPersonResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "apiserver_person" "test" {
  name        = "MyName"
	age         = 42
	description = %[1]q
}
`, configurableAttribute)
}
