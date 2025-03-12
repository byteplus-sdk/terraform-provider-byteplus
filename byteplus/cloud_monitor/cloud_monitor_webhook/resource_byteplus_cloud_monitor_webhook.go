package cloud_monitor_webhook

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CloudMonitorWebhook can be imported using the id, e.g.
```
$ terraform import byteplus_cloud_monitor_webhook.default resource_id
```

*/

func ResourceByteplusCloudMonitorWebhook() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCloudMonitorWebhookCreate,
		Read:   resourceByteplusCloudMonitorWebhookRead,
		Update: resourceByteplusCloudMonitorWebhookUpdate,
		Delete: resourceByteplusCloudMonitorWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the webhook.\n\nLength limit must not exceed 512 bytes.\nThe name can be repeated.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the webhook.\n\ncustom：custom webhook\nwecom：WeChat webhook\nlark：Lark webhook\ndingtalk：DingTalk webhook.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The address of the webhook.",
			},

			// computed fields
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
	}
	return resource
}

func resourceByteplusCloudMonitorWebhookCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorWebhookService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCloudMonitorWebhook())
	if err != nil {
		return fmt.Errorf("error on creating cloud_monitor_webhook %q, %s", d.Id(), err)
	}
	return resourceByteplusCloudMonitorWebhookRead(d, meta)
}

func resourceByteplusCloudMonitorWebhookRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorWebhookService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCloudMonitorWebhook())
	if err != nil {
		return fmt.Errorf("error on reading cloud_monitor_webhook %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCloudMonitorWebhookUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorWebhookService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCloudMonitorWebhook())
	if err != nil {
		return fmt.Errorf("error on updating cloud_monitor_webhook %q, %s", d.Id(), err)
	}
	return resourceByteplusCloudMonitorWebhookRead(d, meta)
}

func resourceByteplusCloudMonitorWebhookDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorWebhookService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCloudMonitorWebhook())
	if err != nil {
		return fmt.Errorf("error on deleting cloud_monitor_webhook %q, %s", d.Id(), err)
	}
	return err
}
