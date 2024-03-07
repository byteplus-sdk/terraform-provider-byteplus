package ssl_vpn_client_cert

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
SSL VPN Client Cert can be imported using the id, e.g.
```
$ terraform import byteplus_ssl_vpn_client_cert.default vsc-2d6b7gjrzc2yo58ozfcx2****
```

*/

func ResourceByteplusSslClientCertServer() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusSslVpnClientCertCreate,
		Read:   resourceByteplusSslVpnClientCertRead,
		Update: resourceByteplusSslVpnClientCertUpdate,
		Delete: resourceByteplusSslVpnClientCertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ssl_vpn_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the ssl vpn server.",
			},
			"ssl_vpn_client_cert_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the ssl vpn client cert.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the ssl vpn client cert.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the ssl vpn client.",
			},
			"certificate_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the ssl vpn client cert.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the ssl vpn client cert.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the ssl vpn client cert.",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expired time of the ssl vpn client cert.",
			},
			"ca_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CA certificate.",
			},
			"client_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The client certificate.",
			},
			"client_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key of the ssl vpn client.",
			},
			"open_vpn_client_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The config of the open vpn client.",
			},
		},
	}
	return resource
}

func resourceByteplusSslVpnClientCertCreate(d *schema.ResourceData, meta interface{}) (err error) {
	SslVpnClientCertService := NewSslVpnClientCertService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(SslVpnClientCertService, d, ResourceByteplusSslClientCertServer())
	if err != nil {
		return fmt.Errorf("error on creating SSL Vpn Client Cert %q, %s", d.Id(), err)
	}
	return resourceByteplusSslVpnClientCertRead(d, meta)
}

func resourceByteplusSslVpnClientCertRead(d *schema.ResourceData, meta interface{}) (err error) {
	SslVpnClientCertService := NewSslVpnClientCertService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(SslVpnClientCertService, d, ResourceByteplusSslClientCertServer())
	if err != nil {
		return fmt.Errorf("error on reading SSL Vpn Client Cert %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusSslVpnClientCertUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	SslVpnClientCertService := NewSslVpnClientCertService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(SslVpnClientCertService, d, ResourceByteplusSslClientCertServer())
	if err != nil {
		return fmt.Errorf("error on updating SSL Vpn Client Cert %q, %s", d.Id(), err)
	}
	return resourceByteplusSslVpnClientCertRead(d, meta)
}

func resourceByteplusSslVpnClientCertDelete(d *schema.ResourceData, meta interface{}) (err error) {
	SslVpnClientCertService := NewSslVpnClientCertService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(SslVpnClientCertService, d, ResourceByteplusSslClientCertServer())
	if err != nil {
		return fmt.Errorf("error on deleting SSL Vpn Client Cert %q, %s", d.Id(), err)
	}
	return err
}
