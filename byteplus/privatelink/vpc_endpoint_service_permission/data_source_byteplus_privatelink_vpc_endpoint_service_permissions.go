package vpc_endpoint_service_permission

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusPrivatelinkVpcEndpointServicePermissions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusPrivatelinkVpcEndpointServicePermissionRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Id of service.",
			},
			"permit_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of permit account.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Returns the total amount of the data list.",
			},
			"permissions": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permit_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The permit account id.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusPrivatelinkVpcEndpointServicePermissionRead(d *schema.ResourceData, meta interface{}) error {
	service := NewService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(service, d, DataSourceByteplusPrivatelinkVpcEndpointServicePermissions())
}
