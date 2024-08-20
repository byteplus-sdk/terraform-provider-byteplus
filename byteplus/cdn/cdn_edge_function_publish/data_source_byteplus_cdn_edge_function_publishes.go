package cdn_edge_function_publish

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusCdnEdgeFunctionPublishs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCdnEdgeFunctionPublishsRead,
		Schema: map[string]*schema.Schema{
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the function.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of query.",
			},

			"tickets": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ticket_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The release record id.",
						},
						"function_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The function id.",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The content of the release record.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the release record.",
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creator of the release record.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The create time of the release record. Displayed in UNIX timestamp format.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The update time of the release record. Displayed in UNIX timestamp format.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusCdnEdgeFunctionPublishsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCdnEdgeFunctionPublishService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCdnEdgeFunctionPublishs())
}
