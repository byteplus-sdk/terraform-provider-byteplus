package direct_connect_virtual_interface

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
DirectConnectVirtualInterface can be imported using the id, e.g.
```
$ terraform import byteplus_direct_connect_virtual_interface.default resource_id
```

*/

func ResourceByteplusDirectConnectVirtualInterface() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusDirectConnectVirtualInterfaceCreate,
		Read:   resourceByteplusDirectConnectVirtualInterfaceRead,
		Update: resourceByteplusDirectConnectVirtualInterfaceUpdate,
		Delete: resourceByteplusDirectConnectVirtualInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"virtual_interface_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of virtual interface.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of virtual interface.",
			},
			"direct_connect_connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The direct connect connection ID which associated with.",
			},
			"direct_connect_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The direct connect gateway ID which associated with.",
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The VLAN ID used to connect to the local IDC, please ensure that this VLAN ID is not occupied, the value range: 0 ~ 2999.",
			},
			"local_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The local IP that associated with.",
			},
			"peer_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The peer IP that associated with.",
			},
			"route_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The route type of virtual interface,valid value contains `Static`,`BGP`.",
			},
			"enable_bfd": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether enable BFD detect.",
			},
			"bfd_detect_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The BFD detect interval.",
			},
			"bfd_detect_multiplier": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The BFD detect times.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The band width limit of virtual interface,in Mbps.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The tags that direct connect gateway added.",
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
			"enable_nqa": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether enable NQA detect.",
			},
			"nqa_detect_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The NQA detect interval.",
			},
			"nqa_detect_multiplier": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The NAQ detect times.",
			},
		},
	}
	return resource
}

func resourceByteplusDirectConnectVirtualInterfaceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectVirtualInterfaceService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusDirectConnectVirtualInterface())
	if err != nil {
		return fmt.Errorf("error on creating direct_connect_virtual_interface %q, %s", d.Id(), err)
	}
	return resourceByteplusDirectConnectVirtualInterfaceRead(d, meta)
}

func resourceByteplusDirectConnectVirtualInterfaceRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectVirtualInterfaceService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusDirectConnectVirtualInterface())
	if err != nil {
		return fmt.Errorf("error on reading direct_connect_virtual_interface %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusDirectConnectVirtualInterfaceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectVirtualInterfaceService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusDirectConnectVirtualInterface())
	if err != nil {
		return fmt.Errorf("error on updating direct_connect_virtual_interface %q, %s", d.Id(), err)
	}
	return resourceByteplusDirectConnectVirtualInterfaceRead(d, meta)
}

func resourceByteplusDirectConnectVirtualInterfaceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDirectConnectVirtualInterfaceService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusDirectConnectVirtualInterface())
	if err != nil {
		return fmt.Errorf("error on deleting direct_connect_virtual_interface %q, %s", d.Id(), err)
	}
	return err
}
