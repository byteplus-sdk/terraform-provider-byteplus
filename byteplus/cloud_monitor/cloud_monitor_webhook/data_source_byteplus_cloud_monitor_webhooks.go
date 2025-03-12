package cloud_monitor_webhook

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusCloudMonitorWebhooks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCloudMonitorWebhooksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"name", "rule_id", "event_rule_id", "type", "url"},
				Description:   "A list of webhook IDs.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "Webhook name, fuzzy search by name.",
			},
			"rule_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "Alarm strategy ID.",
			},
			"event_rule_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "Event Rule ID.",
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "Type of the webhook.\n\ncustom：Custom webhook\nwecom：WeChat webhook\nlark：Lark webhook\ndingtalk：DingTalk webhook.",
			},
			"url": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "The address of the webhook.",
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
			"webhooks": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the webhook.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the webhook.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the webhook.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The address of the webhook.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the webhook.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The update time of the webhook.",
						},
						"rule_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Alarm strategy IDs.",
						},
						"event_rule_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Event rule IDs.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusCloudMonitorWebhooksRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCloudMonitorWebhookService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCloudMonitorWebhooks())
}
