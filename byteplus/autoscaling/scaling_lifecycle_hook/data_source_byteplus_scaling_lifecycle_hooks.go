package scaling_lifecycle_hook

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusScalingLifecycleHooks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusScalingLifecycleHooksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of lifecycle hook ids.",
			},
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Set:         schema.HashString,
				Description: "An id of scaling group id.",
			},
			"lifecycle_hook_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of lifecycle hook names.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of lifecycle hook.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of lifecycle hook query.",
			},
			"lifecycle_hooks": {
				Description: "The collection of lifecycle hook query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the lifecycle hook.",
						},
						"lifecycle_hook_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the lifecycle hook.",
						},
						"scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the scaling group.",
						},
						"lifecycle_hook_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the lifecycle hook.",
						},
						"lifecycle_hook_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timeout of the lifecycle hook.",
						},
						"lifecycle_hook_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the lifecycle hook.",
						},
						"lifecycle_hook_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy of the lifecycle hook.",
						},
						"lifecycle_command": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Batch job command.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"command_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Batch job command ID, which indicates the batch job command to be executed after triggering the lifecycle hook and installed in the instance.",
									},
									"parameters": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameters and parameter values in batch job commands.\nThe number of parameters ranges from 0 to 60.",
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

func dataSourceByteplusScalingLifecycleHooksRead(d *schema.ResourceData, meta interface{}) error {
	lifecycleHookService := NewScalingLifecycleHookService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(lifecycleHookService, d, DataSourceByteplusScalingLifecycleHooks())
}
