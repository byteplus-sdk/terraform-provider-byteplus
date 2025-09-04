package kafka_zone

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusZonesRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Id of Region.",
			},
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
			"zones": {
				Description: "The collection of zone query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the zone.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the zone.",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the zone.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the zone.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the zone.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusZonesRead(d *schema.ResourceData, meta interface{}) error {
	zoneService := NewZoneService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(zoneService, d, DataSourceByteplusZones())
}
