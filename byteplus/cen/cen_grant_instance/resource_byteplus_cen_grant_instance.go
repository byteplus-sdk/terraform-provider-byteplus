package cen_grant_instance

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
Cen grant instance can be imported using the CenId:CenOwnerId:InstanceId:InstanceType:RegionId, e.g.
```
$ terraform import byteplus_cen_grant_instance.default cen-7qthudw0ll6jmc***:210000****:vpc-2fexiqjlgjif45oxruvso****:VPC:cn-beijing
```

*/

func ResourceByteplusCenGrantInstance() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCenGrantInstanceCreate,
		Read:   resourceByteplusCenGrantInstanceRead,
		Delete: resourceByteplusCenGrantInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: cenGrantInstanceImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the cen.",
			},
			"cen_owner_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The owner ID of the cen.",
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
				Description:  "The type of the instance.",
				ValidateFunc: validation.StringInSlice([]string{"VPC", "DCGW"}, false),
			},
			"instance_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region ID of the instance.",
			},
		},
	}
	return resource
}

func resourceByteplusCenGrantInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	grantInstanceService := NewCenGrantInstanceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(grantInstanceService, d, ResourceByteplusCenGrantInstance())
	if err != nil {
		return fmt.Errorf("error on creating cen grant instance  %q, %s", d.Id(), err)
	}
	return resourceByteplusCenGrantInstanceRead(d, meta)
}

func resourceByteplusCenGrantInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	grantInstanceService := NewCenGrantInstanceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(grantInstanceService, d, ResourceByteplusCenGrantInstance())
	if err != nil {
		return fmt.Errorf("error on reading cen grant instance %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCenGrantInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	grantInstanceService := NewCenGrantInstanceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(grantInstanceService, d, ResourceByteplusCenGrantInstance())
	if err != nil {
		return fmt.Errorf("error on deleting cen grant instance %q, %s", d.Id(), err)
	}
	return err
}
