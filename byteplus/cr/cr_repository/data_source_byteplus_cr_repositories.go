package cr_repository

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusCrRepositories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCrRepositoriesRead,
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The CR instance name.",
			},
			"namespaces": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The list of instance namespace.",
			},
			"names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The list of instance names.",
			},
			"access_levels": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The list of instance access level.",
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
			"repositories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The collection of repository query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The namespace of repository.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of repository.",
						},
						"access_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The access level of repository.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of repository.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of repository.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last update time of repository.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusCrRepositoriesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCrRepositoryService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(service, d, DataSourceByteplusCrRepositories())
}
