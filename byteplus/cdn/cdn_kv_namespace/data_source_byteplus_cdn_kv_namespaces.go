package cdn_kv_namespace

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusCdnKvNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCdnKvNamespacesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of Resource.",
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

			"kv_namespaces": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the kv namespace.",
						},
						"namespace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the kv namespace.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the kv namespace.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the kv namespace.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the project to which the namespace belongs, defaulting to `default`.",
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creator of the kv namespace.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The creation time of the kv namespace. Displayed in UNIX timestamp format.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The update time of the kv namespace. Displayed in UNIX timestamp format.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusCdnKvNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCdnKvNamespaceService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCdnKvNamespaces())
}
