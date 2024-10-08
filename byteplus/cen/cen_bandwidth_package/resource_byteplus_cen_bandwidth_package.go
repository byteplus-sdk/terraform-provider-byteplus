package cen_bandwidth_package

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
CenBandwidthPackage can be imported using the id, e.g.
```
$ terraform import byteplus_cen_bandwidth_package.default cbp-4c2zaavbvh5f42****
```

*/

func ResourceByteplusCenBandwidthPackage() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCenBandwidthPackageCreate,
		Read:   resourceByteplusCenBandwidthPackageRead,
		Update: resourceByteplusCenBandwidthPackageUpdate,
		Delete: resourceByteplusCenBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"local_geographic_region_set_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "China",
				ValidateFunc: validation.StringInSlice([]string{"China", "Asia"}, false),
				Description:  "The local geographic region set id of the cen bandwidth package. Valid value: `China`, `Asia`.",
			},
			"peer_geographic_region_set_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "China",
				ValidateFunc: validation.StringInSlice([]string{"China", "Asia"}, false),
				Description:  "The peer geographic region set id of the cen bandwidth package. Valid value: `China`, `Asia`.",
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(2, 100000),
				Description:  "The bandwidth of the cen bandwidth package. Value: 2~10000.",
			},
			"cen_bandwidth_package_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the cen bandwidth package.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the cen bandwidth package.",
			},
			"billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PrePaid",
				ValidateFunc: validation.StringInSlice([]string{"PrePaid"}, false),
				Description:  "The billing type of the cen bandwidth package. Only support `PrePaid` and default value is `PrePaid`.",
			},
			"period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "Month",
				ValidateFunc:     validation.StringInSlice([]string{"Month", "Year"}, false),
				DiffSuppressFunc: periodDiffSuppress,
				Description:      "The period unit of the cen bandwidth package. Value: `Month`, `Year`. Default value is `Month`.",
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: periodDiffSuppress,
				Description:      "The period of the cen bandwidth package. Default value is 1.",
			},
			"tags": bp.TagsSchema(),
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ProjectName of the cen bandwidth package.",
			},
		},
	}
	s := DataSourceByteplusCenBandwidthPackages().Schema["bandwidth_packages"].Elem.(*schema.Resource).Schema
	delete(s, "id")
	bp.MergeDateSourceToResource(s, &resource.Schema)
	return resource
}

func resourceByteplusCenBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenBandwidthPackageService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusCenBandwidthPackage())
	if err != nil {
		return fmt.Errorf("error on creating cen bandwidth package %q, %s", d.Id(), err)
	}
	return resourceByteplusCenBandwidthPackageRead(d, meta)
}

func resourceByteplusCenBandwidthPackageRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenBandwidthPackageService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCenBandwidthPackage())
	if err != nil {
		return fmt.Errorf("error on reading cen bandwidth package %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCenBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenBandwidthPackageService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusCenBandwidthPackage())
	if err != nil {
		return fmt.Errorf("error on updating cen bandwidth package %q, %s", d.Id(), err)
	}
	return resourceByteplusCenBandwidthPackageRead(d, meta)
}

func resourceByteplusCenBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCenBandwidthPackageService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCenBandwidthPackage())
	if err != nil {
		return fmt.Errorf("error on deleting cen bandwidth package %q, %s", d.Id(), err)
	}
	return err
}
