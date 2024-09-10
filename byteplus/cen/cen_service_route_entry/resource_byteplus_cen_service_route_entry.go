package cen_service_route_entry

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
CenServiceRouteEntry can be imported using the CenId:DestinationCidrBlock:ServiceRegionId:ServiceVpcId, e.g.
```
$ terraform import byteplus_cen_service_route_entry.default cen-2nim00ybaylts7trquyzt****:100.XX.XX.0/24:cn-beijing:vpc-3rlkeggyn6tc010exd32q****
```

*/

func ResourceByteplusCenServiceRouteEntry() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCenServiceRouteEntryCreate,
		Update: resourceByteplusCenServiceRouteEntryUpdate,
		Read:   resourceByteplusCenServiceRouteEntryRead,
		Delete: resourceByteplusCenServiceRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: cenServiceRouteEntryImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The cen ID of the cen service route entry.",
			},
			"destination_cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "The destination cidr block of the cen service route entry.",
			},
			"service_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service region id of the cen service route entry.",
			},
			"service_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service VPC id of the cen service route entry.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the cen service route entry.",
			},
			"publish_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LocalDCGW",
				ValidateFunc: validation.StringInSlice([]string{
					"LocalDCGW",
					"Custom",
				}, false),
				Description: "Publishing scope of cloud service access routes. Valid values are `LocalDCGW`(default), `Custom`.",
			},
			"publish_to_instances": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    100,
				Description: "The publish instances. A maximum of 100 can be uploaded in one request. This field needs to be filled in when the `publish_mode` is `Custom`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_region_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The region where the cloud service access route needs to be published.",
						},
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"VPC",
								"DCGW",
							}, false),
							Description: "The network instance type that needs to be published for cloud service access routes. The values are as follows: `VPC`, `DCGW`.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud service access routes need to publish the network instance ID.",
						},
					},
				},
			},
		},
	}
	s := DataSourceByteplusCenServiceRouteEntries().Schema["service_route_entries"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(s, &resource.Schema)
	return resource
}

func resourceByteplusCenServiceRouteEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenServiceRouteEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusCenServiceRouteEntry())
	if err != nil {
		return fmt.Errorf("error on creating cen service route entry %q, %s", d.Id(), err)
	}
	return resourceByteplusCenServiceRouteEntryRead(d, meta)
}

func resourceByteplusCenServiceRouteEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenServiceRouteEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCenServiceRouteEntry())
	if err != nil {
		return fmt.Errorf("error on reading cen service route entry %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCenServiceRouteEntryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenServiceRouteEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusCenServiceRouteEntry())
	if err != nil {
		return fmt.Errorf("error on updating cen service route entry %q, %s", d.Id(), err)
	}
	return resourceByteplusCenServiceRouteEntryRead(d, meta)
}

func resourceByteplusCenServiceRouteEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenServiceRouteEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCenServiceRouteEntry())
	if err != nil {
		return fmt.Errorf("error on deleting cen service route entry %q, %s", d.Id(), err)
	}
	return err
}
