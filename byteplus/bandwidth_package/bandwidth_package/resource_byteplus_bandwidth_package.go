package bandwidth_package

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
BandwidthPackage can be imported using the id, e.g.
```
$ terraform import byteplus_bandwidth_package.default bwp-2zeo05qre24nhrqpy****
```

*/

func ResourceByteplusBandwidthPackage() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusBandwidthPackageCreate,
		Read:   resourceByteplusBandwidthPackageRead,
		Update: resourceByteplusBandwidthPackageUpdate,
		Delete: resourceByteplusBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_package_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the bandwidth package.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the bandwidth package.",
			},
			"isp": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "BGP",
				ForceNew:    true,
				Description: "Route type, default to BGP.",
			},
			"billing_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "PostPaidByBandwidth",
				Description: "BillingType of the bandwidth package. Valid values: `PrePaid`, `PostPaidByBandwidth`(Default), `PostPaidByTraffic`, `PayBy95Peak`." +
					" The billing method of IPv6 does not include `PrePaid`, and the billing method is only based on the `PostPaidByBandwidth`.",
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "Bandwidth upper limit of shared bandwidth package, unit: Mbps. " +
					"When BillingType is set to PrePaid: the value range is 5 to 5000. " +
					"When BillingType is set to PostPaidByBandwidth: the value range is 2 to 5000. " +
					"When BillingType is set to PostPaidByTraffic: the value range is 2 to 2000. " +
					"When BillingType is set to PayBy95Peak: the value range is 2 to 5000.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The IP protocol values for shared bandwidth packages are as follows: `IPv4`: IPv4 protocol. `IPv6`: IPv6 protocol.",
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("billing_type").(string) != "PrePaid"
				},
				Description: "Duration of purchasing shared bandwidth package on an annual or monthly basis. " +
					"The valid value range in 1~9 or 12, 24 or 36. Default value is 1. The period unit defaults to `Month`.",
			},
			"security_protection_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security protection types for shared bandwidth packages. " +
					"Parameter - N: Indicates the number of security protection types, currently only supports taking 1. Value: `AntiDDoS_Enhanced` or left blank." +
					"If the value is `AntiDDoS_Enhanced`, then will create a shared bandwidth package with enhanced protection," +
					" which supports adding basic protection type public IP addresses." +
					"If left blank, it indicates a shared bandwidth package with basic protection, " +
					"which supports the addition of public IP addresses with enhanced protection.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the bandwidth package.",
			},
			"tags": bp.TagsSchema(),
		},
	}
	return resource
}

func resourceByteplusBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewBandwidthPackageService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusBandwidthPackage())
	if err != nil {
		return fmt.Errorf("error on creating bandwidth_package %q, %s", d.Id(), err)
	}
	return resourceByteplusBandwidthPackageRead(d, meta)
}

func resourceByteplusBandwidthPackageRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewBandwidthPackageService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusBandwidthPackage())
	if err != nil {
		return fmt.Errorf("error on reading bandwidth_package %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewBandwidthPackageService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusBandwidthPackage())
	if err != nil {
		return fmt.Errorf("error on updating bandwidth_package %q, %s", d.Id(), err)
	}
	return resourceByteplusBandwidthPackageRead(d, meta)
}

func resourceByteplusBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewBandwidthPackageService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusBandwidthPackage())
	if err != nil {
		return fmt.Errorf("error on deleting bandwidth_package %q, %s", d.Id(), err)
	}
	return err
}
