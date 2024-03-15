package cen_attach_instance

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
Cen attach instance can be imported using the CenId:InstanceId:InstanceType:RegionId, e.g.
```
$ terraform import byteplus_cen_attach_instance.default cen-7qthudw0ll6jmc***:vpc-2fexiqjlgjif45oxruvso****:VPC:cn-beijing
```

*/

func ResourceByteplusCenAttachInstance() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCenAttachInstanceCreate,
		Read:   resourceByteplusCenAttachInstanceRead,
		Delete: resourceByteplusCenAttachInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: cenAttachInstanceImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the cen.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the instance.",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The type of the instance. Valid values: `VPC`, `DCGW`.",
				ValidateFunc: validation.StringInSlice([]string{"VPC", "DCGW"}, false),
			},
			"instance_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region ID of the instance.",
			},
			"instance_owner_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The owner ID of the instance.",
			},
		},
	}
	s := DataSourceByteplusCenAttachInstances().Schema["attach_instances"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(s, &resource.Schema)
	return resource
}

func resourceByteplusCenAttachInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	cenAttachInstanceService := NewCenAttachInstanceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(cenAttachInstanceService, d, ResourceByteplusCenAttachInstance())
	if err != nil {
		return fmt.Errorf("error on creating cen attach instance  %q, %s", d.Id(), err)
	}
	return resourceByteplusCenAttachInstanceRead(d, meta)
}

func resourceByteplusCenAttachInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	cenAttachInstanceService := NewCenAttachInstanceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(cenAttachInstanceService, d, ResourceByteplusCenAttachInstance())
	if err != nil {
		return fmt.Errorf("error on reading cen attach instance %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCenAttachInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	cenAttachInstanceService := NewCenAttachInstanceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(cenAttachInstanceService, d, ResourceByteplusCenAttachInstance())
	if err != nil {
		return fmt.Errorf("error on deleting cen attach instance %q, %s", d.Id(), err)
	}
	return err
}
