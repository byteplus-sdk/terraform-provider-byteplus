package cdn_edge_function

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusCdnEdgeFunctions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCdnEdgeFunctionsRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the function. \n100: running. \n400: unassociated. \n500: configuring.",
			},
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

			"edge_functions": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the edge function.",
						},
						"function_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the edge function.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the edge function.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The remark of the edge function.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of the edge function.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the project to which the edge function belongs.",
						},
						"account_identity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The account id of the edge function.",
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creator of the edge function.",
						},
						"user_identity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The user id of the edge function.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The create time of the edge function. Displayed in UNIX timestamp format.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The update time of the edge function. Displayed in UNIX timestamp format.",
						},
						"source_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest code content of the edge function. The code is transformed into a Base64-encoded format.",
						},
						"domain": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The domain name bound to the edge function.",
						},
						"envs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The environment variables of the edge function.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The key of the environment variable.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the environment variable.",
									},
								},
							},
						},
						"continent_cluster": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The canary cluster info of the edge function.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"country": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The country where the cluster is located.",
									},
									"continent": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The continent where the cluster is located.",
									},
									"cluster_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The type of the cluster.",
									},
									"traffics": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The versions of the function deployed on this cluster.",
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

func dataSourceByteplusCdnEdgeFunctionsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCdnEdgeFunctionService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCdnEdgeFunctions())
}
