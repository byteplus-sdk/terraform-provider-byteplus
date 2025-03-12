package cloud_monitor_contact_group

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusCloudMonitorContactGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCloudMonitorContactGroupsRead,
		Schema: map[string]*schema.Schema{
			//"ids": {
			//	Type:     schema.TypeSet,
			//	Optional: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//	Set:           schema.HashString,
			//	ConflictsWith: []string{"name"},
			//	Description:   "A list of cloud monitor contact group ids.",
			//},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				//ConflictsWith: []string{"ids"},
				Description: "The keyword of the contact group names. Fuzzy match is supported.",
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
			"groups": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the contact group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the contact group.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the account.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the contact group.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The update time.",
						},
						"contacts": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Contact information in the contact group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of the contact.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of contact.",
									},
									"phone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The phone of contact.",
									},
									"email": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The email of contact.",
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

func dataSourceByteplusCloudMonitorContactGroupsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCloudMonitorContactGroupService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCloudMonitorContactGroups())
}
