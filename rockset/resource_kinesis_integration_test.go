package rockset

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/openapi"
)

const testKinesisIntegrationName = "terraform-provider-acceptance-test-kinesis-integration"
const testKinesisIntegrationDescription = "Terraform provider acceptance tests."
const testKinesisIntegrationRoleArn = "arn:aws:iam::469279130686:role/terraform-provider-rockset-tests-kinesis"

func TestAccKinesisIntegration_Basic(t *testing.T) {
	var kinesisIntegration openapi.KinesisIntegration

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckRocksetKinesisIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: getHCL("kinesis_integration.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRocksetKinesisIntegrationExists("rockset_kinesis_integration.test",
						&kinesisIntegration),
					resource.TestCheckResourceAttr("rockset_kinesis_integration.test", "name",
						testKinesisIntegrationName),
					resource.TestCheckResourceAttr("rockset_kinesis_integration.test", "description",
						testKinesisIntegrationDescription),
					resource.TestCheckResourceAttr("rockset_kinesis_integration.test", "aws_role_arn",
						testKinesisIntegrationRoleArn),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccCheckRocksetKinesisIntegrationDestroy(s *terraform.State) error {
	rc := testAccProvider.Meta().(*rockset.RockClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rockset_kinesis_integration" {
			continue
		}

		name := rs.Primary.ID
		_, err := rc.GetIntegration(testCtx, name)
		// An error would mean we didn't find the it, we expect an error
		if err == nil {
			return err
		}
	}

	return nil
}

func testAccCheckRocksetKinesisIntegrationExists(resource string, kinesisIntegration *openapi.KinesisIntegration) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rc := testAccProvider.Meta().(*rockset.RockClient)

		rs, err := getResourceFromState(state, resource)
		if err != nil {
			return err
		}

		name := rs.Primary.ID
		resp, err := rc.GetIntegration(testCtx, name)
		if err != nil {
			return err
		}

		*kinesisIntegration = *resp.Kinesis

		return nil
	}
}
