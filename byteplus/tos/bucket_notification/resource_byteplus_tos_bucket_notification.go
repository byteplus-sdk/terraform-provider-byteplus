package tos_bucket_notification

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
TosBucketNotification can be imported using the bucketName, e.g.
```
$ terraform import byteplus_tos_bucket_notification.default resource_id
```

*/

func ResourceByteplusTosBucketNotification() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusTosBucketNotificationCreate,
		Read:   resourceByteplusTosBucketNotificationRead,
		Update: resourceByteplusTosBucketNotificationUpdate,
		Delete: resourceByteplusTosBucketNotificationDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
				}
				if err := data.Set("bucket_name", items[0]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				if err := data.Set("rules.0.rule_id", items[1]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				return []*schema.ResourceData{data}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the bucket.",
			},
			"rules": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				MinItems:    1,
				Description: "The notification rule of the bucket.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The rule name of the notification.",
						},
						"events": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The event type of the notification.",
						},
						"destination": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "The destination info of the notification.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ve_faas": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The VeFaas info of the destination.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"function_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The function id of the destination.",
												},
											},
										},
									},
								},
							},
						},
						"filter": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The filter of the notification.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tos_key": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The tos filter of the notification.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filter_rules": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The filter rules of the notification.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The name of the filter rule. Valid values: `prefix`, `suffix`.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The value of the filter rule.",
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
				},
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the notification.",
			},
		},
	}
	return resource
}

func resourceByteplusTosBucketNotificationCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketNotificationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusTosBucketNotification())
	if err != nil {
		return fmt.Errorf("error on creating tos_bucket_notification %q, %s", d.Id(), err)
	}
	return resourceByteplusTosBucketNotificationRead(d, meta)
}

func resourceByteplusTosBucketNotificationRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketNotificationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusTosBucketNotification())
	if err != nil {
		return fmt.Errorf("error on reading tos_bucket_notification %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusTosBucketNotificationUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketNotificationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusTosBucketNotification())
	if err != nil {
		return fmt.Errorf("error on updating tos_bucket_notification %q, %s", d.Id(), err)
	}
	return resourceByteplusTosBucketNotificationRead(d, meta)
}

func resourceByteplusTosBucketNotificationDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketNotificationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusTosBucketNotification())
	if err != nil {
		return fmt.Errorf("error on deleting tos_bucket_notification %q, %s", d.Id(), err)
	}
	return err
}
