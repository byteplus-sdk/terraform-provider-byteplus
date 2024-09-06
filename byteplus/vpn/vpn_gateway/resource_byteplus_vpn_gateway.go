package vpn_gateway

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
VpnGateway can be imported using the id, e.g.
```
$ terraform import byteplus_vpn_gateway.default vgw-273zkshb2qayo7fap8t2****
```

*/

func ResourceByteplusVpnGateway() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusVpnGatewayCreate,
		Read:   resourceByteplusVpnGatewayRead,
		Update: resourceByteplusVpnGatewayUpdate,
		Delete: resourceByteplusVpnGatewayDelete,
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
				Description: "The ID of the VPC where you want to create the VPN gateway.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the subnet where you want to create the VPN gateway.",
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{5, 10, 20, 50, 100, 200, 500, 1000}),
				Description:  "The bandwidth of the VPN gateway. Unit: Mbps. Values: 5, 10, 20, 50, 100, 200, 500.",
			},
			"vpn_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the VPN gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the VPN gateway.",
			},
			"billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PrePaid",
				ValidateFunc: validation.StringInSlice([]string{"PrePaid"}, false),
				Description: "The BillingType of the VPN gateway. Only support `PrePaid`.\n" +
					"Terraform will only remove the PrePaid VPN gateway from the state file, not actually remove.",
			},
			//"period_unit": {
			//	Type:         schema.TypeString,
			//	Optional:     true,
			//	ForceNew:     true,
			//	Default:      "Month",
			//	ValidateFunc: validation.StringInSlice([]string{"Month"}, false),
			//	DiffSuppressFunc: periodDiffSuppress,
			//	Description:  "The PeriodUnit of the VPN gateway.",
			//},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          12,
				DiffSuppressFunc: periodDiffSuppress,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36})),
				Description: "The Period of the VPN gateway. Default value is 12. This parameter is only useful when creating vpn gateway. Default period unit is Month.\n" +
					"Value range: 1~9, 12, 24, 36. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"renew_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The renew type of the VPN gateway.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the VPN gateway.",
			},
			"tags": bp.TagsSchema(),
			"ipsec_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether ipsec is enabled.",
			},
			"ssl_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether ssl is enabled.",
			},
			"ssl_max_connections": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "The max connections of ssl. " +
					"This parameter can only be passed in when ssl_enabled is true. Default is 5.",
			},
		},
	}
	dataSource := DataSourceByteplusVpnGateways().Schema["vpn_gateways"].Elem.(*schema.Resource).Schema
	delete(dataSource, "id")
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusVpnGatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpnGatewayService := NewVpnGatewayService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(vpnGatewayService, d, ResourceByteplusVpnGateway())
	if err != nil {
		return fmt.Errorf("error on creating Vpn Gateway %q, %s", d.Id(), err)
	}
	return resourceByteplusVpnGatewayRead(d, meta)
}

func resourceByteplusVpnGatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpnGatewayService := NewVpnGatewayService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(vpnGatewayService, d, ResourceByteplusVpnGateway())
	if err != nil {
		return fmt.Errorf("error on reading Vpn Gateway %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpnGatewayService := NewVpnGatewayService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(vpnGatewayService, d, ResourceByteplusVpnGateway())
	if err != nil {
		return fmt.Errorf("error on updating Vpn Gateway %q, %s", d.Id(), err)
	}
	return resourceByteplusVpnGatewayRead(d, meta)
}

func resourceByteplusVpnGatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpnGatewayService := NewVpnGatewayService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(vpnGatewayService, d, ResourceByteplusVpnGateway())
	if err != nil {
		return fmt.Errorf("error on deleting Vpn Gateway %q, %s", d.Id(), err)
	}
	return err
}
