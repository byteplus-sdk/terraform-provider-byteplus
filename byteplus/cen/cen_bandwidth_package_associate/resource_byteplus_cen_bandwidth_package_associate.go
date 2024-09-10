package cen_bandwidth_package_associate

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Cen bandwidth package associate can be imported using the CenBandwidthPackageId:CenId, e.g.
```
$ terraform import byteplus_cen_bandwidth_package_associate.default cbp-4c2zaavbvh5fx****:cen-7qthudw0ll6jmc****
```

*/

func ResourceByteplusCenBandwidthPackageAssociate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCenBandwidthPackageAssociateCreate,
		Read:   resourceByteplusCenBandwidthPackageAssociateRead,
		Delete: resourceByteplusCenBandwidthPackageAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: cenGrantInstanceImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_bandwidth_package_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the cen bandwidth package.",
			},
			"cen_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the cen.",
			},
		},
	}
	return resource
}

func resourceByteplusCenBandwidthPackageAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenBandwidthPackageAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusCenBandwidthPackageAssociate())
	if err != nil {
		return fmt.Errorf("error on creating cen bandwidth package associate %q, %s", d.Id(), err)
	}
	return resourceByteplusCenBandwidthPackageAssociateRead(d, meta)
}

func resourceByteplusCenBandwidthPackageAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenBandwidthPackageAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCenBandwidthPackageAssociate())
	if err != nil {
		return fmt.Errorf("error on reading cen bandwidth package associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCenBandwidthPackageAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenBandwidthPackageAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCenBandwidthPackageAssociate())
	if err != nil {
		return fmt.Errorf("error on deleting cen bandwidth package associate %q, %s", d.Id(), err)
	}
	return err
}
