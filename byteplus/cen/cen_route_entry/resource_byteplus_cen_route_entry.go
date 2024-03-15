package cen_route_entry

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
CenRouteEntry can be imported using the CenId:DestinationCidrBlock:InstanceId:InstanceType:InstanceRegionId, e.g.
```
$ terraform import byteplus_cen_route_entry.default cen-2nim00ybaylts7trquyzt****:100.XX.XX.0/24:vpc-vtbnbb04qw3k2hgi12cv****:VPC:cn-beijing
```

*/

func ResourceByteplusCenRouteEntry() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCenRouteEntryCreate,
		Read:   resourceByteplusCenRouteEntryRead,
		Delete: resourceByteplusCenRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: cenRouteEntryImporter,
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
				Description: "The cen ID of the cen route entry.",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "VPC",
				ValidateFunc: validation.StringInSlice([]string{"VPC"}, false),
				Description:  "The instance type of the next hop of the cen route entry.",
			},
			"instance_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The instance region id of the next hop of the cen route entry.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The instance id of the next hop of the cen route entry.",
			},
			"destination_cidr_block": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The destination cidr block of the cen route entry.",
			},
		},
	}
	s := DataSourceByteplusCenRouteEntries().Schema["cen_route_entries"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(s, &resource.Schema)
	return resource
}

func resourceByteplusCenRouteEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenRouteEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusCenRouteEntry())
	if err != nil {
		return fmt.Errorf("error on creating cen route entry %q, %s", d.Id(), err)
	}
	return resourceByteplusCenRouteEntryRead(d, meta)
}

func resourceByteplusCenRouteEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenRouteEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCenRouteEntry())
	if err != nil {
		return fmt.Errorf("error on reading cen route entry %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCenRouteEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenRouteEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCenRouteEntry())
	if err != nil {
		return fmt.Errorf("error on deleting cen route entry %q, %s", d.Id(), err)
	}
	return err
}
