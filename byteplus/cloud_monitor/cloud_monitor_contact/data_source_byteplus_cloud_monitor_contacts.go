package cloud_monitor_contact

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusCloudMonitorContacts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCloudMonitorContactsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"email"},
				Description:   "A list of Contact IDs.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "The keyword of contact names. Fuzzy match is supported.",
			},
			"email": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "The email of the cloud monitor contact. This field support fuzzy query.",
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
			"contacts": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of contact.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of contact.",
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
	}
}

func dataSourceByteplusCloudMonitorContactsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(service, d, DataSourceByteplusCloudMonitorContacts())
}
