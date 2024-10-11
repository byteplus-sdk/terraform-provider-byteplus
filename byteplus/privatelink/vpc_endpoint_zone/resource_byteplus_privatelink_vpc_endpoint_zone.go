package vpc_endpoint_zone

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
VpcEndpointZone can be imported using the endpointId:subnetId, e.g.
```
$ terraform import byteplus_privatelink_vpc_endpoint_zone.default ep-3rel75r081l345zsk2i59****:subnet-2bz47q19zhx4w2dx0eevn****
```

*/

func ResourceByteplusPrivatelinkVpcEndpointZone() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusPrivateLinkVpcEndpointZoneCreate,
		Read:   resourceByteplusPrivateLinkVpcEndpointZoneRead,
		Delete: resourceByteplusPrivateLinkVpcEndpointZoneDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
				}
				if err := data.Set("endpoint_id", items[0]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				if err := data.Set("subnet_id", items[1]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				return []*schema.ResourceData{data}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"endpoint_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The endpoint id of vpc endpoint zone.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet id of vpc endpoint zone.",
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The private ip address of vpc endpoint zone.",
			},

			"zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Id of vpc endpoint zone.",
			},
			"zone_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain of vpc endpoint zone.",
			},
			"zone_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of vpc endpoint zone.",
			},
			"network_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network interface id of vpc endpoint.",
			},
		},
	}
	return resource
}

func resourceByteplusPrivateLinkVpcEndpointZoneCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcEndpointZoneService := NewVpcEndpointZoneService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(vpcEndpointZoneService, d, ResourceByteplusPrivatelinkVpcEndpointZone())
	if err != nil {
		return fmt.Errorf("error on creating vpc endpoint zone %q, %w", d.Id(), err)
	}
	return resourceByteplusPrivateLinkVpcEndpointZoneRead(d, meta)
}

func resourceByteplusPrivateLinkVpcEndpointZoneRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcEndpointZoneService := NewVpcEndpointZoneService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(vpcEndpointZoneService, d, ResourceByteplusPrivatelinkVpcEndpointZone())
	if err != nil {
		return fmt.Errorf("error on reading vpc endpoint zone %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusPrivateLinkVpcEndpointZoneDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcEndpointZoneService := NewVpcEndpointZoneService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(vpcEndpointZoneService, d, ResourceByteplusPrivatelinkVpcEndpointZone())
	if err != nil {
		return fmt.Errorf("error on deleting vpc endpoint zone %q, %w", d.Id(), err)
	}
	return err
}
