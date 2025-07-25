package kafka_sasl_user

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusKafkaSaslUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusKafkaSaslUsersRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of instance.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user name, support fuzzy matching.",
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
			"users": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of user.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of user.",
						},
						"password_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of password.",
						},
						"all_authority": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this user has read and write permissions for all topics.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusKafkaSaslUsersRead(d *schema.ResourceData, meta interface{}) error {
	service := NewKafkaSaslUserService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusKafkaSaslUsers())
}
