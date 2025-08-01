package cr_authorization_token

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusCrAuthorizationTokens() *schema.Resource {
	return &schema.Resource{
		Read: func(data *schema.ResourceData, meta interface{}) error {
			service := NewCrAuthorizationTokenService(meta.(*bp.SdkClient))
			return bp.DefaultDispatcher().Data(service, data, DataSourceByteplusCrAuthorizationTokens())
		},
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cr instance name want to query.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of instance query.",
			},
			"tokens": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The collection of users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Temporary access token.",
						},
						"username": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username for login repository instance.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expiration time of the temporary access token.",
						},
					},
				},
			},
		},
	}
}
