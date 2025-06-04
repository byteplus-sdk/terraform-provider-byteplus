package account

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusRedisAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusRedisAccountsRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of redis accounts query.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the Redis instance.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the redis account.",
			},
			"accounts": {
				Description: "The collection of redis instance account query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the redis account.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of instance.",
						},
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The role info.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the redis account.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusRedisAccountsRead(d *schema.ResourceData, meta interface{}) error {
	redisAccountService := NewAccountService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(redisAccountService, d, DataSourceByteplusRedisAccounts())
}
