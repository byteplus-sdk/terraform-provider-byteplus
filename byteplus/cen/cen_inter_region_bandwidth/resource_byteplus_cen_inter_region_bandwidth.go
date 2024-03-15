package cen_inter_region_bandwidth

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
CenInterRegionBandwidth can be imported using the id, e.g.
```
$ terraform import byteplus_cen_inter_region_bandwidth.default cirb-3tex2x1cwd4c6c0v****
```

*/

func ResourceByteplusCenInterRegionBandwidth() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCenInterRegionBandwidthCreate,
		Read:   resourceByteplusCenInterRegionBandwidthRead,
		Update: resourceByteplusCenInterRegionBandwidthUpdate,
		Delete: resourceByteplusCenInterRegionBandwidthDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Description: "The cen ID of the cen inter region bandwidth.",
			},
			"local_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The local region id of the cen inter region bandwidth.",
			},
			"peer_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The peer region id of the cen inter region bandwidth.",
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The bandwidth of the cen inter region bandwidth.",
			},
		},
	}
	s := DataSourceByteplusCenInterRegionBandwidths().Schema["inter_region_bandwidths"].Elem.(*schema.Resource).Schema
	delete(s, "id")
	bp.MergeDateSourceToResource(s, &resource.Schema)
	return resource
}

func resourceByteplusCenInterRegionBandwidthCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenInterRegionBandwidthService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusCenInterRegionBandwidth())
	if err != nil {
		return fmt.Errorf("error on creating cen inter region bandwidth %q, %s", d.Id(), err)
	}
	return resourceByteplusCenInterRegionBandwidthRead(d, meta)
}

func resourceByteplusCenInterRegionBandwidthRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenInterRegionBandwidthService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCenInterRegionBandwidth())
	if err != nil {
		return fmt.Errorf("error on reading cen inter region bandwidth %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCenInterRegionBandwidthUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenInterRegionBandwidthService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusCenInterRegionBandwidth())
	if err != nil {
		return fmt.Errorf("error on updating cen inter region bandwidth %q, %s", d.Id(), err)
	}
	return resourceByteplusCenInterRegionBandwidthRead(d, meta)
}

func resourceByteplusCenInterRegionBandwidthDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenInterRegionBandwidthService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCenInterRegionBandwidth())
	if err != nil {
		return fmt.Errorf("error on deleting cen inter region bandwidth %q, %s", d.Id(), err)
	}
	return err
}
