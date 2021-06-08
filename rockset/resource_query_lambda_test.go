package rockset

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/stretchr/testify/assert"
)

const testQueryLambdaName = "terraform-provider-acceptance-tests-query-lambda-basic"

func TestAccQueryLambda_Basic(t *testing.T) {
	var queryLambda openapi.QueryLambda

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRocksetQueryLambdaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckQueryLambdaBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRocksetQueryLambdaExists("rockset_query_lambda.test", &queryLambda),
					resource.TestCheckResourceAttr("rockset_query_lambda.test", "name", testQueryLambdaName),
					resource.TestCheckResourceAttr("rockset_query_lambda.test", "description", "basic lambda"),
					testAccCheckSql(t, &queryLambda, "SELECT * FROM commons._events WHERE _events._event_time > :timestamp LIMIT 1"),
					resource.TestCheckResourceAttrSet("rockset_query_lambda.test", "version"),
					resource.TestCheckResourceAttrSet("rockset_query_lambda.test", "state"),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccCheckQueryLambdaUpdated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRocksetQueryLambdaExists("rockset_query_lambda.test", &queryLambda),
					resource.TestCheckResourceAttr("rockset_query_lambda.test", "name", testQueryLambdaName),
					resource.TestCheckResourceAttr("rockset_query_lambda.test", "description", "updated description"),
					testAccCheckSql(t, &queryLambda, "SELECT * FROM commons._events WHERE _events._event_time >= :timestamp LIMIT 2"),
					resource.TestCheckResourceAttrSet("rockset_query_lambda.test", "version"),
					resource.TestCheckResourceAttrSet("rockset_query_lambda.test", "state"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccCheckQueryLambdaBasic() string {
	hclPath := filepath.Join("..", "testdata", "query_lambda_basic.tf")
	hclString, err := getFileContents(hclPath)
	if err != nil {
		log.Fatalf("Unexpected error loading test data %s", hclPath)
	}

	return hclString
}

func testAccCheckQueryLambdaUpdated() string {
	hclPath := filepath.Join("..", "testdata", "query_lambda_basic_updated.tf")
	hclString, err := getFileContents(hclPath)
	if err != nil {
		log.Fatalf("Unexpected error loading test data %s", hclPath)
	}

	return hclString
}

func testAccCheckRocksetQueryLambdaDestroy(s *terraform.State) error {
	rc := testAccProvider.Meta().(*rockset.RockClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rockset_api_key" {
			continue
		}

		workspace, name := workspaceAndNameFromID(rs.Primary.ID)
		_, err := getQueryLambda(testCtx, rc, workspace, name)

		// An error would mean we didn't find the key, we expect an error
		if err == nil {
			// We did not get an error, so we failed to delete the key.
			return fmt.Errorf("Query Lambda %s still exists.", name)
		}
	}

	return nil
}

func testAccCheckRocksetQueryLambdaExists(resource string, queryLambda *openapi.QueryLambda) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rc := testAccProvider.Meta().(*rockset.RockClient)

		rs, err := getResourceFromState(state, resource)
		if err != nil {
			return err
		}

		workspace, name := workspaceAndNameFromID(rs.Primary.ID)

		resp, err := getQueryLambda(testCtx, rc, workspace, name)
		if err != nil {
			return err
		}

		*queryLambda = *resp

		return nil
	}
}

func testAccCheckSql(t *testing.T, queryLambda *openapi.QueryLambda, expectedSql string) resource.TestCheckFunc {
	assert := assert.New(t)

	return func(state *terraform.State) error {
		sql := queryLambda.LatestVersion.Sql.Query

		assert.Equal(sql, expectedSql, "SQL string didn't match.")

		return nil
	}
}