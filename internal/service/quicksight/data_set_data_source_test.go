package quicksight_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/quicksight"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccQuickSightDataSetDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rId := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_quicksight_data_set.test"
	dataSourceName := "data.aws_quicksight_data_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, quicksight.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSetDataSourceConfig_basic(rId, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
				),
			},
		},
	})
}

func testAccDataSetDataSourceConfig_basic(rId, rName string) string {
	return acctest.ConfigCompose(
		testAccDataSetConfigBase(rId, rName),
		fmt.Sprintf(`
resource "aws_quicksight_data_set" "test" {
  data_set_id = %[1]q
  name        = %[2]q
  import_mode = "SPICE"

  physical_table_map {
    physical_table_map_id = %[1]q
    s3_source {
      data_source_arn = aws_quicksight_data_source.test.arn
      input_columns {
        name = "Column1"
        type = "STRING"
      }
      upload_settings {
        format = "JSON"
      }
    }
  }
}

data "aws_quicksight_data_set" "test" {
  data_set_id = aws_quicksight_data_set.test.data_set_id
}
`, rId, rName))
}
