package rds_mysql_account

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusRdsMysqlAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusRdsMysqlAccountsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of database account.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of database account query.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the RDS instance.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the database account. This field supports fuzzy query.",
			},
			"accounts": {
				Description: "The collection of RDS instance account query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the database account.",
						},
						"account_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the database account.",
						},
						"account_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the database account.",
						},
						"account_privileges": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The privilege detail list of RDS mysql instance account.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of database.",
									},
									"account_privilege": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The privilege type of the account.",
									},
									"account_privilege_detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The privilege detail of the account.",
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

func dataSourceByteplusRdsMysqlAccountsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewRdsMysqlAccountService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusRdsMysqlAccounts())
}
