package cloud_monitor_contact

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CloudMonitor Contact can be imported using the id, e.g.
```
$ terraform import byteplus_cloud_monitor_contact.default 145258255725730****
```

*/

func ResourceByteplusCloudMonitorContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusCloudMonitorContactCreate,
		Read:   resourceByteplusCloudMonitorContactRead,
		Update: resourceByteplusCloudMonitorContactUpdate,
		Delete: resourceByteplusCloudMonitorContactDelete,
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
				Description: "The name of contact.",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The email of contact.",
			},
		},
	}
}

func resourceByteplusCloudMonitorContactCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusCloudMonitorContact())
	if err != nil {
		return fmt.Errorf("error on creating Contact %q, %w", d.Id(), err)
	}
	return resourceByteplusCloudMonitorContactRead(d, meta)
}

func resourceByteplusCloudMonitorContactRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCloudMonitorContact())
	if err != nil {
		return fmt.Errorf("error on reading Contact %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusCloudMonitorContactUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusCloudMonitorContact())
	if err != nil {
		return fmt.Errorf("error on updating Contact %q, %w", d.Id(), err)
	}
	return resourceByteplusCloudMonitorContactRead(d, meta)
}

func resourceByteplusCloudMonitorContactDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCloudMonitorContact())
	if err != nil {
		return fmt.Errorf("error on deleting Contact %q, %w", d.Id(), err)
	}
	return err
}
