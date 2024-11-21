package direct_connect_gateway

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
DirectConnectGateway can be imported using the id, e.g.
```
$ terraform import byteplus_direct_connect_gateway.default resource_id
```

*/

func ResourceByteplusDirectConnectGateway() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusDirectConnectGatewayCreate,
		Read:   resourceByteplusDirectConnectGatewayRead,
		Update: resourceByteplusDirectConnectGatewayUpdate,
		Delete: resourceByteplusDirectConnectGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"direct_connect_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of direct connect gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of direct connect gateway.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The direct connect gateway tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The tag value.",
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceByteplusDirectConnectGatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on creating direct_connect_gateway %q, %s", d.Id(), err)
	}
	return resourceByteplusDirectConnectGatewayRead(d, meta)
}

func resourceByteplusDirectConnectGatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on reading direct_connect_gateway %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusDirectConnectGatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on updating direct_connect_gateway %q, %s", d.Id(), err)
	}
	return resourceByteplusDirectConnectGatewayRead(d, meta)
}

func resourceByteplusDirectConnectGatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectGatewayService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on deleting direct_connect_gateway %q, %s", d.Id(), err)
	}
	return err
}
