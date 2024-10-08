package vpn_connection

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
VpnConnection can be imported using the id, e.g.
```
$ terraform import byteplus_vpn_connection.default vgc-3tex2x1cwd4c6c0v****
```

*/

func ResourceByteplusVpnConnection() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusVpnConnectionCreate,
		Read:   resourceByteplusVpnConnectionRead,
		Update: resourceByteplusVpnConnectionUpdate,
		Delete: resourceByteplusVpnConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpn_connection_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the VPN connection.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the VPN connection.",
			},
			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the vpn gateway. If the `AttachType` is not passed or the passed value is `VpnGateway`, this parameter must be filled. If the value of `AttachType` is `TransitRouter`, this parameter does not need to be filled.",
			},
			"customer_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the customer gateway.",
			},
			"local_subnet": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      schema.HashString,
				MinItems: 1,
				MaxItems: 30,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Description: "The local subnet of the VPN connection. Up to 5 network segments are supported.",
			},
			"remote_subnet": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      schema.HashString,
				MinItems: 1,
				MaxItems: 30,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Description: "The remote subnet of the VPN connection. Up to 5 network segments are supported.",
			},
			"dpd_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "clear",
				ValidateFunc: validation.StringInSlice([]string{"clear", "none", "hold", "restart"}, false),
				Description:  "The dpd action of the VPN connection.",
			},
			"nat_traversal": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The nat traversal of the VPN connection.",
			},

			// ike config
			"ike_config_psk": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The psk of the ike config of the VPN connection. The length does not exceed 100 characters, and only uppercase and lowercase letters, special symbols and numbers are allowed.",
			},
			"ike_config_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ikev1",
				ValidateFunc: validation.StringInSlice([]string{"ikev1", "ikev2"}, false),
				Description:  "The version of the ike config of the VPN connection. The value can be `ikev1` or `ikev2`. The default value is `ikev1`.",
			},
			"ike_config_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "main",
				ValidateFunc: validation.StringInSlice([]string{"main", "aggressive"}, false),
				Description:  "The mode of the ike config of the VPN connection. Valid values are `main`, `aggressive`, and default value is `main`.",
			},
			"ike_config_enc_alg": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "aes",
				ValidateFunc: validation.StringInSlice([]string{"aes", "aes192", "aes256", "des", "3des", "sm4"}, false),
				Description:  "The enc alg of the ike config of the VPN connection. Valid value are `aes`, `aes192`, `aes256`, `des`, `3des`, `sm4`. The default value is `aes`.",
			},
			"ike_config_auth_alg": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "sha1",
				ValidateFunc: validation.StringInSlice([]string{"sha1", "md5", "sha256", "sha384", "sha512", "sm3"}, false),
				Description:  "The auth alg of the ike config of the VPN connection. Valid value are `sha1`, `md5`, `sha256`, `sha384`, `sha512`, `sm3`. The default value is `sha1`.",
			},
			"ike_config_dh_group": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "group2",
				ValidateFunc: validation.StringInSlice([]string{"group1", "group2", "group5", "group14"}, false),
				Description:  "The dk group of the ike config of the VPN connection. Valid value are `group1`, `group2`, `group5`, `group14`. The default value is `group2`.",
			},
			"ike_config_lifetime": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				ValidateFunc: validation.IntBetween(900, 86400),
				Description:  "The lifetime of the ike config of the VPN connection. Value: 900~86400.",
			},
			"ike_config_local_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The local_id of the ike config of the VPN connection.",
			},
			"ike_config_remote_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The remote id of the ike config of the VPN connection.",
			},

			// ipsec config
			"ipsec_config_enc_alg": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "aes",
				ValidateFunc: validation.StringInSlice([]string{"aes", "aes192", "aes256", "des", "3des", "sm4"}, false),
				Description:  "The enc alg of the ipsec config of the VPN connection. Valid value are `aes`, `aes192`, `aes256`, `des`, `3des`, `sm4`. The default value is `aes`.",
			},
			"ipsec_config_auth_alg": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "sha1",
				ValidateFunc: validation.StringInSlice([]string{"sha1", "md5", "sha256", "sha384", "sha512", "sm3"}, false),
				Description:  "The auth alg of the ipsec config of the VPN connection. Valid value are `sha1`, `md5`, `sha256`, `sha384`, `sha512`, `sm3`. The default value is `sha1`.",
			},
			"ipsec_config_dh_group": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "group2",
				ValidateFunc: validation.StringInSlice([]string{"group1", "group2", "group5", "group14", "disable"}, false),
				Description:  "The dh group of the ipsec config of the VPN connection. Valid value are `group1`, `group2`, `group5`, `group14` and `disable`. The default value is `group2`.",
			},
			"ipsec_config_lifetime": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				ValidateFunc: validation.IntBetween(900, 86400),
				Description:  "The ipsec config of the ike config of the VPN connection. Value: 900~86400.",
			},
			"negotiate_instantly": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to initiate negotiation mode immediately.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the VPN connection.",
			},
			"attach_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "VpnGateway",
				Description: "The attach type of the VPN connection, the value can be `VpnGateway` or `TransitRouter`.",
			},
			"log_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable connection logging. After enabling Connection Day, you can view and download IPsec connection logs, and use the log information to troubleshoot IPsec connection problems yourself.",
			},
		},
	}
	dataSource := DataSourceByteplusVpnConnections().Schema["vpn_connections"].Elem.(*schema.Resource).Schema
	delete(dataSource, "id")
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusVpnConnectionCreate(d *schema.ResourceData, meta interface{}) (err error) {
	connectionService := NewVpnConnectionService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(connectionService, d, ResourceByteplusVpnConnection())
	if err != nil {
		return fmt.Errorf("error on creating Vpn Connections %q, %s", d.Id(), err)
	}
	return resourceByteplusVpnConnectionRead(d, meta)
}

func resourceByteplusVpnConnectionRead(d *schema.ResourceData, meta interface{}) (err error) {
	connectionService := NewVpnConnectionService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(connectionService, d, ResourceByteplusVpnConnection())
	if err != nil {
		return fmt.Errorf("error on reading Vpn Connection %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	connectionService := NewVpnConnectionService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(connectionService, d, ResourceByteplusVpnConnection())
	if err != nil {
		return fmt.Errorf("error on updating Vpn Connection %q, %s", d.Id(), err)
	}
	return resourceByteplusVpnConnectionRead(d, meta)
}

func resourceByteplusVpnConnectionDelete(d *schema.ResourceData, meta interface{}) (err error) {
	connectionService := NewVpnConnectionService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(connectionService, d, ResourceByteplusVpnConnection())
	if err != nil {
		return fmt.Errorf("error on deleting Vpn Connection %q, %s", d.Id(), err)
	}
	return err
}
