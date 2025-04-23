package cloud_monitor_contact_group

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CloudMonitorContactGroup can be imported using the id, e.g.
```
$ terraform import byteplus_cloud_monitor_contact_group.default resource_id
```

*/

func ResourceByteplusCloudMonitorContactGroup() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCloudMonitorContactGroupCreate,
		Read:   resourceByteplusCloudMonitorContactGroupRead,
		Update: resourceByteplusCloudMonitorContactGroupUpdate,
		Delete: resourceByteplusCloudMonitorContactGroupDelete,
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
				Description: "The name of the contact group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the contact group.",
			},
			"contacts_id_list": {
				Type:     schema.TypeSet,
				Set:      schema.HashString,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "When creating a contact group, contacts should be added with their contact ID. " +
					"The maximum number of IDs allowed is 100, meaning that the maximum number of members in a single contact group is 100.",
			},
		},
	}
	return resource
}

func resourceByteplusCloudMonitorContactGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorContactGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCloudMonitorContactGroup())
	if err != nil {
		return fmt.Errorf("error on creating cloud_monitor_contact_group %q, %s", d.Id(), err)
	}
	return resourceByteplusCloudMonitorContactGroupRead(d, meta)
}

func resourceByteplusCloudMonitorContactGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorContactGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCloudMonitorContactGroup())
	if err != nil {
		return fmt.Errorf("error on reading cloud_monitor_contact_group %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCloudMonitorContactGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorContactGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCloudMonitorContactGroup())
	if err != nil {
		return fmt.Errorf("error on updating cloud_monitor_contact_group %q, %s", d.Id(), err)
	}
	return resourceByteplusCloudMonitorContactGroupRead(d, meta)
}

func resourceByteplusCloudMonitorContactGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorContactGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCloudMonitorContactGroup())
	if err != nil {
		return fmt.Errorf("error on deleting cloud_monitor_contact_group %q, %s", d.Id(), err)
	}
	return err
}
