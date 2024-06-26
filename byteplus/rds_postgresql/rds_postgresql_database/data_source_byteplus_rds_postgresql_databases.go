package rds_postgresql_database

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusRdsPostgresqlDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusRdsPostgresqlDatabasesRead,
		Schema: map[string]*schema.Schema{
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
							Description: "The name of the RDS database.",
						},
						"db_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the RDS database.",
						},
						"character_set_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The character set of the RDS database.",
						},
						"collate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The collate of database.",
						},
						"c_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Character classification.",
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The owner of database.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusRdsPostgresqlDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	databaseService := NewRdsPostgresqlDatabaseService(meta.(*bp.SdkClient))
	return databaseService.Dispatcher.Data(databaseService, d, DataSourceByteplusRdsPostgresqlDatabases())
}
