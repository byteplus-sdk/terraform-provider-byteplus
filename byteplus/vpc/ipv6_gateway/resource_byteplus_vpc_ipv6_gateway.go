package ipv6_gateway

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Ipv6Gateway can be imported using the id, e.g.
```
$ terraform import byteplus_vpc_ipv6_gateway.default ipv6gw-12bcapllb5ukg17q7y2sd3thx
```

*/

func ResourceByteplusIpv6Gateway() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusIpv6GatewayCreate,
		Read:   resourceByteplusIpv6GatewayRead,
		Update: resourceByteplusIpv6GatewayUpdate,
		Delete: resourceByteplusIpv6GatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VPC which the Ipv6Gateway belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Ipv6Gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the Ipv6Gateway.",
			},
		},
	}
	dataSource := DataSourceByteplusIpv6Gateways().Schema["ipv6_gateways"].Elem.(*schema.Resource).Schema
	delete(dataSource, "id")
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusIpv6GatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6GatewayService := NewIpv6GatewayService(meta.(*bp.SdkClient))
	err = ipv6GatewayService.Dispatcher.Create(ipv6GatewayService, d, ResourceByteplusIpv6Gateway())
	if err != nil {
		return fmt.Errorf("error on creating Ipv6Gateway %q, %w", d.Id(), err)
	}
	return resourceByteplusIpv6GatewayRead(d, meta)
}

func resourceByteplusIpv6GatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6GatewayService := NewIpv6GatewayService(meta.(*bp.SdkClient))
	err = ipv6GatewayService.Dispatcher.Read(ipv6GatewayService, d, ResourceByteplusIpv6Gateway())
	if err != nil {
		return fmt.Errorf("error on reading Ipv6Gateway %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusIpv6GatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6GatewayService := NewIpv6GatewayService(meta.(*bp.SdkClient))
	err = ipv6GatewayService.Dispatcher.Update(ipv6GatewayService, d, ResourceByteplusIpv6Gateway())
	if err != nil {
		return fmt.Errorf("error on updating Ipv6Gateway %q, %w", d.Id(), err)
	}
	return resourceByteplusIpv6GatewayRead(d, meta)
}

func resourceByteplusIpv6GatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	ipv6GatewayService := NewIpv6GatewayService(meta.(*bp.SdkClient))
	err = ipv6GatewayService.Dispatcher.Delete(ipv6GatewayService, d, ResourceByteplusIpv6Gateway())
	if err != nil {
		return fmt.Errorf("error on deleting Ipv6Gateway %q, %w", d.Id(), err)
	}
	return err
}
