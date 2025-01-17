package direct_connect_connection

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
DirectConnectConnection can be imported using the id, e.g.
```
$ terraform import byteplus_direct_connect_connection.default dcc-7qthudw0ll6jmc****
```

*/

func ResourceByteplusDirectConnectConnection() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusDirectConnectConnectionCreate,
		Read:   resourceByteplusDirectConnectConnectionRead,
		Update: resourceByteplusDirectConnectConnectionUpdate,
		Delete: resourceByteplusDirectConnectConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"direct_connect_connection_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of direct connect.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of direct connect.",
			},
			"direct_connect_access_point_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The direct connect access point id.",
			},
			"line_operator": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The physical leased line operator.valid value contains `ChinaTelecom`,`ChinaMobile`,`ChinaUnicom`,`ChinaOther`.",
			},
			"port_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The physical leased line port type and spec.valid value contains `1000Base-T`,`10GBase-T`,`1000Base`,`10GBase`,`40GBase`,`100GBase`.",
			},
			"port_spec": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The physical leased line port spec.valid value contains `1G`,`10G`.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The line band width,unit:Mbps.",
			},
			"peer_location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The local IDC address.",
			},
			"customer_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated line contact name.",
			},
			"customer_contact_phone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated line contact phone.",
			},
			"customer_contact_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dedicated line contact email.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The physical leased line tags.",
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

func resourceByteplusDirectConnectConnectionCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectConnectionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusDirectConnectConnection())
	if err != nil {
		return fmt.Errorf("error on creating direct_connect_connection %q, %s", d.Id(), err)
	}
	return resourceByteplusDirectConnectConnectionRead(d, meta)
}

func resourceByteplusDirectConnectConnectionRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectConnectionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusDirectConnectConnection())
	if err != nil {
		return fmt.Errorf("error on reading direct_connect_connection %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusDirectConnectConnectionUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectConnectionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusDirectConnectConnection())
	if err != nil {
		return fmt.Errorf("error on updating direct_connect_connection %q, %s", d.Id(), err)
	}
	return resourceByteplusDirectConnectConnectionRead(d, meta)
}

func resourceByteplusDirectConnectConnectionDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectConnectionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusDirectConnectConnection())
	if err != nil {
		return fmt.Errorf("error on deleting direct_connect_connection %q, %s", d.Id(), err)
	}
	return err
}
