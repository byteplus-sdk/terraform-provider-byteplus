package iam_access_key

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusIamAccessKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusIamAccessKeysRead,
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user names.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of IAM.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of user query.",
			},
			"access_key_metadata": {
				Description: "The collection of access keys.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user name.",
						},
						"access_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user access key id.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user access key create date.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user access key status.",
						},
						"update_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user access key update date.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusIamAccessKeysRead(d *schema.ResourceData, meta interface{}) error {
	iamAccessKeyService := NewIamAccessKeyService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(iamAccessKeyService, d, DataSourceByteplusIamAccessKeys())
}
