package cloud_monitor_object_group

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusCloudMonitorObjectGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCloudMonitorObjectGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of cloud monitor object group ids.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The keyword of the object group names. Fuzzy match is supported.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of Resource.",
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
			"object_groups": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource group ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource group name.",
						},
						"alert_template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alarm template ID associated with the resource group.",
						},
						"alert_template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alarm template name associated with the resource group.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the resource group.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The update time of the resource group.",
						},
						"objects": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of cloud product resources under the resource group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource grouping ID.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The product space to which the cloud product belongs in cloud monitoring.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of resource collection.",
									},
									"region": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Availability zone associated with the cloud product under the current resource.",
									},
									"dimensions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of cloud product resource IDs.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Key for retrieving metrics.",
												},
												"value": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "Value corresponding to the Key.",
												},
											},
										},
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

func dataSourceByteplusCloudMonitorObjectGroupsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCloudMonitorObjectGroupService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCloudMonitorObjectGroups())
}
