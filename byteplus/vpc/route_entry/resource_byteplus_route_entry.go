package route_entry

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
Route entry can be imported using the route_table_id:route_entry_id, e.g.
```
$ terraform import byteplus_route_entry.default vtb-274e19skkuhog7fap8u4i8ird:rte-274e1g9ei4k5c7fap8sp974fq
```

*/

func ResourceByteplusRouteEntry() *schema.Resource {
	return &schema.Resource{
		Delete: resourceByteplusRouteEntryDelete,
		Create: resourceByteplusRouteEntryCreate,
		Read:   resourceByteplusRouteEntryRead,
		Update: resourceByteplusRouteEntryUpdate,
		Importer: &schema.ResourceImporter{
			State: importRouteEntry,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the route table.",
			},
			"route_entry_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the route entry.",
			},
			"destination_cidr_block": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The destination CIDR block of the route entry.",
			},
			"next_hop_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The type of the next hop, Optional choice contains `Instance`, `NetworkInterface`, `NatGW`, `VpnGW`, `TransitRouter`.",
				ValidateFunc: validation.StringInSlice([]string{"Instance", "NetworkInterface", "NatGW", "VpnGW", "TransitRouter"}, false),
			},
			"next_hop_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the next hop.",
			},
			"route_entry_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the route entry.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the route entry.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the route entry.",
			},
		},
	}
}

func resourceByteplusRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	routeEntryService := NewRouteEntryService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(routeEntryService, d, ResourceByteplusRouteEntry()); err != nil {
		return fmt.Errorf("error on creating route entry  %q, %w", d.Id(), err)
	}
	return resourceByteplusRouteEntryRead(d, meta)
}

func resourceByteplusRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	routeEntryService := NewRouteEntryService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(routeEntryService, d, ResourceByteplusRouteEntry()); err != nil {
		return fmt.Errorf("error on reading route entry %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusRouteEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	routeEntryService := NewRouteEntryService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Update(routeEntryService, d, ResourceByteplusRouteEntry()); err != nil {
		return fmt.Errorf("error on updating route entry %q, %w", d.Id(), err)
	}
	return resourceByteplusRouteEntryRead(d, meta)
}

func resourceByteplusRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	routeEntryService := NewRouteEntryService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(routeEntryService, d, ResourceByteplusRouteEntry()); err != nil {
		return fmt.Errorf("error on deleting route entry %q, %w", d.Id(), err)
	}
	return nil
}
