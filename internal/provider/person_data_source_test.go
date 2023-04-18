package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/tcarreira/api-server/pkg/client"
	apiTypes "github.com/tcarreira/api-server/pkg/types"
)

func TestAccPersonDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			c, _ := client.NewAPIClient(client.Config{Endpoint: "http://localhost:18080"})
			c.People().Create(&apiTypes.Person{Name: "personXXX", Age: 49})
		},
		CheckDestroy: func(s *terraform.State) error {
			c, _ := client.NewAPIClient(client.Config{Endpoint: "http://localhost:18080"})
			c.People().Delete(0)
			return nil
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccPersonDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.apiserver_person.test", "name", "personXXX"),
					resource.TestCheckResourceAttr("data.apiserver_person.test", "age", "49"),
				),
			},
		},
	})
}

const testAccPersonDataSourceConfig = `
data "apiserver_person" "test" {
  id = 0
}
`
