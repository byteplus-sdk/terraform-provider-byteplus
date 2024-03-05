package ipv6_address_bandwidth

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Ipv6AddressBandwidth can be imported using the id, e.g.
```
$ terraform import byteplus_vpc_ipv6_address_bandwidth.default eip-2fede9fsgnr4059gp674m6ney
```

*/

func ResourceByteplusIpv6AddressBandwidth() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusIpv6AddressBandwidthCreate,
		Read:   resourceByteplusIpv6AddressBandwidthRead,
		Update: resourceByteplusIpv6AddressBandwidthUpdate,
		Delete: resourceByteplusIpv6AddressBandwidthDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ipv6_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Ipv6 address.",
			},
			"billing_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PostPaidByBandwidth",
					"PostPaidByTraffic",
				}, false),
				Description: "BillingType of the Ipv6 bandwidth. Valid values: `PostPaidByBandwidth`; `PostPaidByTraffic`.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Peek bandwidth of the Ipv6 address. Valid values: 1 to 200. Unit: Mbit/s.",
			},
		},
	}
	dataSource := DataSourceByteplusIpv6AddressBandwidths().Schema["ipv6_address_bandwidths"].Elem.(*schema.Resource).Schema
	delete(dataSource, "id")
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusIpv6AddressBandwidthCreate(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6AddressBandwidthService := NewIpv6AddressBandwidthService(meta.(*bp.SdkClient))
	err = ipv6AddressBandwidthService.Dispatcher.Create(ipv6AddressBandwidthService, d, ResourceByteplusIpv6AddressBandwidth())
	if err != nil {
		return fmt.Errorf("error on creating Ipv6AddressBandwidth %q, %w", d.Id(), err)
	}
	return resourceByteplusIpv6AddressBandwidthRead(d, meta)
}

func resourceByteplusIpv6AddressBandwidthRead(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6AddressBandwidthService := NewIpv6AddressBandwidthService(meta.(*bp.SdkClient))
	err = ipv6AddressBandwidthService.Dispatcher.Read(ipv6AddressBandwidthService, d, ResourceByteplusIpv6AddressBandwidth())
	if err != nil {
		return fmt.Errorf("error on reading Ipv6AddressBandwidth %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusIpv6AddressBandwidthUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6AddressBandwidthService := NewIpv6AddressBandwidthService(meta.(*bp.SdkClient))
	err = ipv6AddressBandwidthService.Dispatcher.Update(ipv6AddressBandwidthService, d, ResourceByteplusIpv6AddressBandwidth())
	if err != nil {
		return fmt.Errorf("error on updating Ipv6AddressBandwidth %q, %w", d.Id(), err)
	}
	return resourceByteplusIpv6AddressBandwidthRead(d, meta)
}

func resourceByteplusIpv6AddressBandwidthDelete(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6AddressBandwidthService := NewIpv6AddressBandwidthService(meta.(*bp.SdkClient))
	err = ipv6AddressBandwidthService.Dispatcher.Delete(ipv6AddressBandwidthService, d, ResourceByteplusIpv6AddressBandwidth())
	if err != nil {
		return fmt.Errorf("error on deleting Ipv6AddressBandwidth %q, %w", d.Id(), err)
	}
	return err
}
