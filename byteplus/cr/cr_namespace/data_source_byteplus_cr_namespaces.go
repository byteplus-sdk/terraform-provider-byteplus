package cr_namespace

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusCrNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCrNamespacesRead,
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The target cr instance name.",
			},
			"names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The list of instance IDs.",
			},
			"projects": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The list of project names to query.",
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
			"namespaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The collection of namespaces query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of OCI repository.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when namespace created.",
						},
						"project": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ProjectName of the CrNamespace.",
						},
						"repository_default_access_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default access level of repository. Valid values: `Private`, `Public`.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusCrNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCrNamespaceService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(service, d, DataSourceByteplusCrNamespaces())
}
