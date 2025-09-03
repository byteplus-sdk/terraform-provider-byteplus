package tos_bucket_realtime_log

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
TosBucketRealtimeLog can be imported using the bucket_name, e.g.
```
$ terraform import byteplus_tos_bucket_realtime_log.default resource_id
```

*/

func ResourceByteplusTosBucketRealtimeLog() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusTosBucketRealtimeLogCreate,
		Read:   resourceByteplusTosBucketRealtimeLogRead,
		Update: resourceByteplusTosBucketRealtimeLogUpdate,
		Delete: resourceByteplusTosBucketRealtimeLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The role name used to grant TOS access to create resources such as projects and topics, and write logs to the TLS logging service. You can use the default TOS role `TOSLogArchiveTLSRole`.",
			},
			"access_log_configuration": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The export schedule of the bucket inventory.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ttl": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     7,
							Description: "The TLS log retention duration. Unit in days. Valid values range is 1~3650. default is 7.",
						},
						// computed fields
						"tls_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the tls project.",
						},
						"tls_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the tls topic.",
						},
						"tls_dashboard_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the tls dashboard.",
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceByteplusTosBucketRealtimeLogCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRealtimeLogService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusTosBucketRealtimeLog())
	if err != nil {
		return fmt.Errorf("error on creating tos_bucket_realtime_log %q, %s", d.Id(), err)
	}
	return resourceByteplusTosBucketRealtimeLogRead(d, meta)
}

func resourceByteplusTosBucketRealtimeLogRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRealtimeLogService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusTosBucketRealtimeLog())
	if err != nil {
		return fmt.Errorf("error on reading tos_bucket_realtime_log %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusTosBucketRealtimeLogUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRealtimeLogService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusTosBucketRealtimeLog())
	if err != nil {
		return fmt.Errorf("error on updating tos_bucket_realtime_log %q, %s", d.Id(), err)
	}
	return resourceByteplusTosBucketRealtimeLogRead(d, meta)
}

func resourceByteplusTosBucketRealtimeLogDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketRealtimeLogService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusTosBucketRealtimeLog())
	if err != nil {
		return fmt.Errorf("error on deleting tos_bucket_realtime_log %q, %s", d.Id(), err)
	}
	return err
}
