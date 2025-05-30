package alb_certificate

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusAlbCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusAlbCertificatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The list of Certificate IDs.",
			},
			"certificate_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Name of Certificate.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of Certificate query.",
			},
			"certificates": {
				Description: "The collection of Certificate query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Certificate.",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Certificate.",
						},
						"certificate_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Certificate.",
						},
						"certificate_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the Certificate.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the Certificate.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time of the Certificate.",
						},
						"expired_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expire time of the Certificate.",
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain name of the Certificate.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the Certificate.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ProjectName of the Certificate.",
						},
						"listeners": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The ID list of the Listener.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusAlbCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	certificateService := NewCertificateService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(certificateService, d, DataSourceByteplusAlbCertificates())
}
