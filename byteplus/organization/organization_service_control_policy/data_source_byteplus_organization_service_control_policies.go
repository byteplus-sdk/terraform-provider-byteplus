package organization_service_control_policy

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusServiceControlPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusServiceControlPoliciesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of Policy query.",
			},
			"query": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query policies, support policy name or description.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of policy. The value can be System or Custom.",
			},
			"policies": {
				Description: "The collection of Policy query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Policy.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Policy.",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the Policy.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time of the Policy.",
						},
						"update_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The update time of the Policy.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the Policy.",
						},
						"statement": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The statement of the Policy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusServiceControlPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(iamPolicyService, d, DataSourceByteplusServiceControlPolicies())
}
