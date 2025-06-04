package region

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusRedisRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusRegionsRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target region info.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of region query.",
			},
			"regions": {
				Description: "The collection of region query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the region.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of region.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusRegionsRead(d *schema.ResourceData, meta interface{}) error {
	redisRegionService := NewRegionService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(redisRegionService, d, DataSourceByteplusRedisRegions())
}
