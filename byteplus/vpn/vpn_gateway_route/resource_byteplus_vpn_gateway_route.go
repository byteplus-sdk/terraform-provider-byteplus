package vpn_gateway_route

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
VpnGatewayRoute can be imported using the id, e.g.
```
$ terraform import byteplus_vpn_gateway_route.default vgr-3tex2c6c0v844c****
```

*/

func ResourceByteplusVpnGatewayRoute() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusVpnGatewayRouteCreate,
		Read:   resourceByteplusVpnGatewayRouteRead,
		Delete: resourceByteplusVpnGatewayRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VPN gateway of the VPN gateway route.",
			},
			"destination_cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "The destination cidr block of the VPN gateway route.",
			},
			"next_hop_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The next hop id of the VPN gateway route.",
			},
		},
	}
	dataSource := DataSourceByteplusVpnGatewayRoutes().Schema["vpn_gateway_routes"].Elem.(*schema.Resource).Schema
	delete(dataSource, "id")
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusVpnGatewayRouteCreate(d *schema.ResourceData, meta interface{}) (err error) {
	routeService := NewVpnGatewayRouteService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(routeService, d, ResourceByteplusVpnGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on creating Vpn Gateway route %q, %s", d.Id(), err)
	}
	return resourceByteplusVpnGatewayRouteRead(d, meta)
}

func resourceByteplusVpnGatewayRouteRead(d *schema.ResourceData, meta interface{}) (err error) {
	routeService := NewVpnGatewayRouteService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(routeService, d, ResourceByteplusVpnGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on reading Vpn Gateway route %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusVpnGatewayRouteDelete(d *schema.ResourceData, meta interface{}) (err error) {
	routeService := NewVpnGatewayRouteService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(routeService, d, ResourceByteplusVpnGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on deleting Vpn Gateway route %q, %s", d.Id(), err)
	}
	return err
}
