package cloud_monitor_rule

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CloudMonitorRule can be imported using the id, e.g.
```
$ terraform import byteplus_cloud_monitor_rule.default 174284623567451****
```

*/

func ResourceByteplusCloudMonitorRule() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCloudMonitorRuleCreate,
		Read:   resourceByteplusCloudMonitorRuleRead,
		Update: resourceByteplusCloudMonitorRuleUpdate,
		Delete: resourceByteplusCloudMonitorRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the cloud monitor rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the cloud monitor rule.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The namespace of the cloud monitor rule.",
			},
			"sub_namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The sub namespace of the cloud monitor rule.",
			},
			"level": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"level", "level_conditions"},
				ValidateFunc: validation.StringInSlice([]string{"critical", "warning", "notice"}, false),
				Description:  "The severity level of the cloud monitor rule. Valid values: `critical`, `warning`, `notice`. One of `level` and `level_conditions` must be specified.",
			},
			"enable_state": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, false),
				Description:  "Whether to enable the cloud monitor rule. Valid values: `enable`, `disable`.",
			},
			"evaluation_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The number of successive times for which the threshold is reached before the alarm is triggered. Unit in minutes. Supports configurations of 1, 3, 5, 10, 15, 30, 60, and 120.",
			},
			"effect_start_at": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The effect start time of the cloud monitor rule. The expression is `HH:MM`.",
			},
			"effect_end_at": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The effect end time of the cloud monitor rule. The expression is `HH:MM`.",
			},
			"silence_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The interval at which alarms are sent. Unit in minutes. Valid values: 5, 30, 60, 180, 360, 720, 1440.",
			},
			"multiple_conditions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to use multiple metrics in the cloud monitor rule.",
			},
			"condition_operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether the alarm is triggered only when the conditions on multiple metrics are met. Valid values: `&&`, `||`.",
			},
			"notify_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Alarm sending aggregation strategy.\n\nrule（default）: aggregation by rule.\nresource: aggregation by rule and resource.",
			},
			"notification_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The notification id of the cloud monitor rule.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the cloud monitor rule.",
			},
			"alert_methods": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The notification method of the cloud monitor rule. Valid values: `Email`, `Webhook`.",
			},
			"web_hook": {
				Type:          schema.TypeString,
				Optional:      true,
				AtLeastOneOf:  []string{"web_hook", "contact_group_ids", "webhook_ids"},
				ConflictsWith: []string{"webhook_ids"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if methods, ok := d.GetOk("alert_methods"); ok {
						methodArr := methods.(*schema.Set).List()
						if contains("Webhook", methodArr) {
							return false
						}
					}
					return true
				},
				Description: "The webhook URL that is used when an alarm is triggered. When the alert method is `Webhook`, one of `web_hook` and `webhook_ids` must be specified.",
			},
			"webhook_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				AtLeastOneOf:  []string{"web_hook", "contact_group_ids", "webhook_ids"},
				ConflictsWith: []string{"web_hook"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if methods, ok := d.GetOk("alert_methods"); ok {
						methodArr := methods.(*schema.Set).List()
						if contains("Webhook", methodArr) {
							return false
						}
					}
					return true
				},
				Description: "The web hook id list of the cloud monitor rule. When the alert method is `Webhook`, one of `web_hook` and `webhook_ids` must be specified.",
			},
			"contact_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:          schema.HashString,
				AtLeastOneOf: []string{"web_hook", "contact_group_ids", "webhook_ids"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if methods, ok := d.GetOk("alert_methods"); ok {
						methodsArr := methods.(*schema.Set).List()
						if contains("Phone", methodsArr) || contains("Email", methodsArr) ||
							contains("SMS", methodsArr) {
							return false
						}
					}
					return true
				},
				Description: "The contact group ids of the cloud monitor rule. When the alert method is `Email`,, This field must be specified.",
			},
			"recovery_notify": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "The recovery notify of the cloud monitor rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Specifies whether alarm recovery notifications are sent.",
						},
					},
				},
			},
			"regions": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The region ids of the cloud monitor rule.",
			},
			"conditions": {
				Type:         schema.TypeSet,
				Optional:     true,
				ExactlyOneOf: []string{"conditions", "level_conditions"},
				Description:  "The conditions that trigger the alarm.\nSpecify an array that contains a maximum of 10 metric math expressions. One of `conditions` and `level_conditions` must be specified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The metric name of the cloud monitor rule.",
						},
						"metric_unit": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The metric unit of the cloud monitor rule.",
						},
						"statistics": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The statistics of the cloud monitor rule. Valid values: `avg`, `max`, `min`.",
						},
						"comparison_operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The comparison operation of the cloud monitor rule. Valid values: `>`, `>=`, `<`, `<=`, `!=`, `=`.",
						},
						"threshold": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The threshold of the cloud monitor rule.",
						},
						"period": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The period of the cloud monitor rule.",
						},
					},
				},
			},
			"level_conditions": {
				Type:         schema.TypeSet,
				Optional:     true,
				ExactlyOneOf: []string{"conditions", "level_conditions"},
				Description:  "The level conditions that trigger the alarm. One of `conditions` and `level_conditions` must be specified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"critical", "warning", "notice"}, false),
							Description:  "The severity level of the cloud monitor rule. Valid values: `critical`, `warning`, `notice`.",
						},
						"conditions": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The conditions that trigger the alarm.\nSpecify an array that contains a maximum of 10 metric math expressions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The metric name of the cloud monitor rule.",
									},
									"metric_unit": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The metric unit of the cloud monitor rule.",
									},
									"statistics": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The statistics of the cloud monitor rule. Valid values: `avg`, `max`, `min`.",
									},
									"comparison_operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The comparison operation of the cloud monitor rule. Valid values: `>`, `>=`, `<`, `<=`, `!=`, `=`.",
									},
									"threshold": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The threshold of the cloud monitor rule.",
									},
									"period": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The period of the cloud monitor rule.",
									},
								},
							},
						},
					},
				},
			},
			"original_dimensions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The original dimensions of the cloud monitor rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The key of the dimension.",
						},
						"value": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set:         schema.HashString,
							Description: "The value of the dimension. If you want to specify all possible values of the dimension, set the value to an asterisk ( * ).",
						},
					},
				},
			},
			"no_data": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "No-data alarm.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Specifies whether to enable no-data alarm. The default value is false.",
						},
						"evaluation_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "No data alarm triggering threshold. When `enable` is set to true, `evaluation_count` is mandatory. The range of values is integers between 3 and 20.",
						},
					},
				},
			},

			// computed fields
			"alert_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alert state of the cloud monitor rule.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The created time of the cloud monitor rule.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The updated time of the cloud monitor rule.",
			},
		},
	}
	return resource
}

func resourceByteplusCloudMonitorRuleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorRuleService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCloudMonitorRule())
	if err != nil {
		return fmt.Errorf("error on creating cloud_monitor_rule %q, %s", d.Id(), err)
	}
	return resourceByteplusCloudMonitorRuleRead(d, meta)
}

func resourceByteplusCloudMonitorRuleRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorRuleService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCloudMonitorRule())
	if err != nil {
		return fmt.Errorf("error on reading cloud_monitor_rule %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCloudMonitorRuleUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorRuleService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCloudMonitorRule())
	if err != nil {
		return fmt.Errorf("error on updating cloud_monitor_rule %q, %s", d.Id(), err)
	}
	return resourceByteplusCloudMonitorRuleRead(d, meta)
}

func resourceByteplusCloudMonitorRuleDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorRuleService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCloudMonitorRule())
	if err != nil {
		return fmt.Errorf("error on deleting cloud_monitor_rule %q, %s", d.Id(), err)
	}
	return err
}

func contains(target string, arr []interface{}) bool {
	for _, v := range arr {
		if target == v.(string) {
			return true
		}
	}
	return false
}
