package iam_role

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusIamRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusIamRolesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of Role.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of Role query.",
			},
			"role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Role, comma separated.",
			},
			"query": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The query field of Role.",
			},
			"roles": {
				Description: "The collection of Role query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Role.",
						},
						"trn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource name of the Role.",
						},
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Role.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time of the Role.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the Role.",
						},
						"trust_policy_document": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The trust policy document of the Role.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusIamRolesRead(d *schema.ResourceData, meta interface{}) error {
	iamRoleService := NewIamRoleService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(iamRoleService, d, DataSourceByteplusIamRoles())
}
