package direct_connect_gateway_route

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
DirectConnectGatewayRoute can be imported using the id, e.g.
```
$ terraform import byteplus_direct_connect_gateway_route.default resource_id
```

*/

func ResourceByteplusDirectConnectGatewayRoute() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusDirectConnectGatewayRouteCreate,
		Read:   resourceByteplusDirectConnectGatewayRouteRead,
		Delete: resourceByteplusDirectConnectGatewayRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"direct_connect_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of direct connect gateway.",
			},
			"destination_cidr_block": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The cidr block.",
			},
			"next_hop_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of next hop.",
			},
		},
	}
	dataSource := DataSourceByteplusDirectConnectGatewayRoutes().Schema["direct_connect_gateway_routes"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusDirectConnectGatewayRouteCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayRouteService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusDirectConnectGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on creating direct_connect_gateway_route %q, %s", d.Id(), err)
	}
	return resourceByteplusDirectConnectGatewayRouteRead(d, meta)
}

func resourceByteplusDirectConnectGatewayRouteRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayRouteService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusDirectConnectGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on reading direct_connect_gateway_route %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusDirectConnectGatewayRouteDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayRouteService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusDirectConnectGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on deleting direct_connect_gateway_route %q, %s", d.Id(), err)
	}
	return err
}
