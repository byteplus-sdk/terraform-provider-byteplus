package ssl_vpn_server

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
SSL VPN server can be imported using the id, e.g.
```
$ terraform import byteplus_ssl_vpn_server.default vss-zm55pqtvk17oq32zd****
```

*/

func ResourceByteplusSslVpnServer() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusSslVpnServerCreate,
		Read:   resourceByteplusSslVpnServerRead,
		Update: resourceByteplusSslVpnServerUpdate,
		Delete: resourceByteplusSslVpnServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The vpn gateway id.",
			},
			"local_subnets": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The local network segment of the SSL server. The local network segment is the address segment that the client accesses through the SSL VPN connection.",
			},
			"client_ip_pool": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SSL client network segment.",
			},
			"ssl_vpn_server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the SSL server.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the ssl server.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Default:      "UDP",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"UDP", "TCP"}, false),
				Description:  "The protocol used by the SSL server. Valid values are `TCP`, `UDP`. Default Value: `UDP`.",
			},
			"cipher": {
				Type:         schema.TypeString,
				Default:      "AES-128-CBC",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AES-128-CBC", "AES-192-CBC", "AES-256-CBC", "None"}, false),
				Description:  "The encryption algorithm of the SSL server.\nValues:\n`AES-128-CBC` (default)\n`AES-192-CBC`\n`AES-256-CBC`\n`None` (do not use encryption).",
			},
			"auth": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "SHA1",
				ValidateFunc: validation.StringInSlice([]string{"SHA1", "MD5", "None"}, false),
				Description:  "The authentication algorithm of the SSL server.\nValues:\n`SHA1` (default)\n`MD5`\n`None` (do not use encryption).",
			},
			"compress": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Whether to compress the transmitted data. The default value is false.",
			},
			"ssl_vpn_server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the ssl vpn server.",
			},
		},
	}
	return resource
}

func resourceByteplusSslVpnServerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	SslVpnServerService := NewSslVpnServerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(SslVpnServerService, d, ResourceByteplusSslVpnServer())
	if err != nil {
		return fmt.Errorf("error on creating SSL Vpn Server %q, %s", d.Id(), err)
	}
	return resourceByteplusSslVpnServerRead(d, meta)
}

func resourceByteplusSslVpnServerRead(d *schema.ResourceData, meta interface{}) (err error) {
	SslVpnServerService := NewSslVpnServerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(SslVpnServerService, d, ResourceByteplusSslVpnServer())
	if err != nil {
		return fmt.Errorf("error on reading SSL Vpn Server %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusSslVpnServerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	SslVpnServerService := NewSslVpnServerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(SslVpnServerService, d, ResourceByteplusSslVpnServer())
	if err != nil {
		return fmt.Errorf("error on updating SSL Vpn Server %q, %s", d.Id(), err)
	}
	return resourceByteplusSslVpnServerRead(d, meta)
}

func resourceByteplusSslVpnServerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	SslVpnServerService := NewSslVpnServerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(SslVpnServerService, d, ResourceByteplusSslVpnServer())
	if err != nil {
		return fmt.Errorf("error on deleting SSL Vpn Server %q, %s", d.Id(), err)
	}
	return err
}
