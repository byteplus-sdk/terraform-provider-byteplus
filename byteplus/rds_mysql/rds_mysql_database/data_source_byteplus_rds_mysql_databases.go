package rds_mysql_database

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusRdsMysqlDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusRdsMysqlDatabasesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of RDS database.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of RDS database query.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the RDS instance.",
			},
			"db_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the RDS database.",
			},
			"databases": {
				Description: "The collection of RDS instance account query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the RDS database. This field supports fuzzy queries.",
						},
						"character_set_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The character set of the RDS database.",
						},
						"database_privileges": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The privilege detail list of RDS mysql instance database.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of account.",
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

func dataSourceByteplusRdsMysqlDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	databaseService := NewRdsMysqlDatabaseService(meta.(*bp.SdkClient))
	return databaseService.Dispatcher.Data(databaseService, d, DataSourceByteplusRdsMysqlDatabases())
}
