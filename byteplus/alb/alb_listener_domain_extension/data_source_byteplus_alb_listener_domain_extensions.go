package alb_listener_domain_extension

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusListenerDomainExtensions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusAlbListenerDomainExtensionsRead,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A Listener ID.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of Listener query.",
			},
			"domain_extensions": {
				Description: "The collection of domain extensions query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Listener.",
						},
						"domain_extension_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The extension domain ID.",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The server certificate ID that domain used.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The listener ID that domain belongs to.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusAlbListenerDomainExtensionsRead(d *schema.ResourceData, meta interface{}) error {
	listenerService := NewAlbListenerDomainExtensionService(meta.(*bp.SdkClient))
	return listenerService.Dispatcher.Data(listenerService, d, DataSourceByteplusListenerDomainExtensions())
}
