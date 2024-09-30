package iam_user_group_policy_attachment

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusIamUserGroupPolicyAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusIamUserGroupPolicyAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"user_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A name of user group.",
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
			"policies": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_trn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource name of the strategy.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the policy.",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the policy.",
						},
						"attach_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Attached time.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusIamUserGroupPolicyAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewIamUserGroupPolicyAttachmentService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusIamUserGroupPolicyAttachments())
}
