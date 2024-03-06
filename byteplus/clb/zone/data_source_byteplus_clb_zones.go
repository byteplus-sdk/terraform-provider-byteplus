package zone

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusClbZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusClbZonesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of zone query.",
			},
			"master_zones": {
				Description: "The master zones list.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The master zone id.",
						},
						"slave_zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The slave zones list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The slave zone id.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusClbZonesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewClbZoneService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(service, d, DataSourceByteplusClbZones())
}
