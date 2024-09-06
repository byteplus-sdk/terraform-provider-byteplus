package customer_gateway

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
CustomerGateway can be imported using the id, e.g.
```
$ terraform import byteplus_customer_gateway.default cgw-2byswc356dybk2dx0eed2****
```

*/

func ResourceByteplusCustomerGateway() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCustomerGatewayCreate,
		Read:   resourceByteplusCustomerGatewayRead,
		Update: resourceByteplusCustomerGatewayUpdate,
		Delete: resourceByteplusCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "The IP address of the customer gateway.",
			},
			"customer_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the customer gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the customer gateway.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the VPN customer gateway.",
			},
		},
	}
	dataSource := DataSourceByteplusCustomerGateways().Schema["customer_gateways"].Elem.(*schema.Resource).Schema
	delete(dataSource, "id")
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	customerGatewayService := NewCustomerGatewayService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(customerGatewayService, d, ResourceByteplusCustomerGateway())
	if err != nil {
		return fmt.Errorf("error on creating Customer Gateway %q, %s", d.Id(), err)
	}
	return resourceByteplusCustomerGatewayRead(d, meta)
}

func resourceByteplusCustomerGatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	customerGatewayService := NewCustomerGatewayService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(customerGatewayService, d, ResourceByteplusCustomerGateway())
	if err != nil {
		return fmt.Errorf("error on reading Customer Gateway %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	customerGatewayService := NewCustomerGatewayService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(customerGatewayService, d, ResourceByteplusCustomerGateway())
	if err != nil {
		return fmt.Errorf("error on updating Customer Gateway %q, %s", d.Id(), err)
	}
	return resourceByteplusCustomerGatewayRead(d, meta)
}

func resourceByteplusCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	customerGatewayService := NewCustomerGatewayService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(customerGatewayService, d, ResourceByteplusCustomerGateway())
	if err != nil {
		return fmt.Errorf("error on deleting Customer Gateway %q, %s", d.Id(), err)
	}
	return err
}
