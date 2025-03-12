package cloud_monitor_object_group

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CloudMonitorObjectGroup can be imported using the id, e.g.
```
$ terraform import byteplus_cloud_monitor_object_group.default resource_id
```

*/

func ResourceByteplusCloudMonitorObjectGroup() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCloudMonitorObjectGroupCreate,
		Read:   resourceByteplusCloudMonitorObjectGroupRead,
		Update: resourceByteplusCloudMonitorObjectGroupUpdate,
		Delete: resourceByteplusCloudMonitorObjectGroupDelete,
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
				Description: "The name of resource group.\n\nCan only contain Chinese, English, or underscores\nThe length is limited to 1-64 characters.",
			},
			"objects": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Need to group the list of cloud product resources, the maximum length of the list is 100.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The product space to which the cloud product belongs in cloud monitoring.",
						},
						"region": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Availability zone associated with the cloud product under the current resource. Only one region id can be specified currently.",
						},
						"dimensions": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Collection of cloud product resource IDs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Key for retrieving metrics.",
									},
									"value": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Set:         schema.HashString,
										Description: "Value corresponding to the Key.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceByteplusCloudMonitorObjectGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorObjectGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCloudMonitorObjectGroup())
	if err != nil {
		return fmt.Errorf("error on creating cloud_monitor_object_group %q, %s", d.Id(), err)
	}
	return resourceByteplusCloudMonitorObjectGroupRead(d, meta)
}

func resourceByteplusCloudMonitorObjectGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorObjectGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCloudMonitorObjectGroup())
	if err != nil {
		return fmt.Errorf("error on reading cloud_monitor_object_group %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCloudMonitorObjectGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorObjectGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCloudMonitorObjectGroup())
	if err != nil {
		return fmt.Errorf("error on updating cloud_monitor_object_group %q, %s", d.Id(), err)
	}
	return resourceByteplusCloudMonitorObjectGroupRead(d, meta)
}

func resourceByteplusCloudMonitorObjectGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCloudMonitorObjectGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCloudMonitorObjectGroup())
	if err != nil {
		return fmt.Errorf("error on deleting cloud_monitor_object_group %q, %s", d.Id(), err)
	}
	return err
}
