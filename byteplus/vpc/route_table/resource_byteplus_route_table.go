package route_table

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Route table can be imported using the id, e.g.
```
$ terraform import byteplus_route_table.default vtb-274e0syt9av407fap8tle16kb
```

*/

func ResourceByteplusRouteTable() *schema.Resource {
	return &schema.Resource{
		Delete: resourceByteplusRouteTableDelete,
		Create: resourceByteplusRouteTableCreate,
		Read:   resourceByteplusRouteTableRead,
		Update: resourceByteplusRouteTableUpdate,
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
				Description: "The id of the VPC.",
			},
			"route_table_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the route table.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the route table.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ProjectName of the route table.",
			},
		},
	}
}

func resourceByteplusRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	routeTableService := NewRouteTableService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(routeTableService, d, ResourceByteplusRouteTable()); err != nil {
		return fmt.Errorf("error on creating route table  %q, %w", d.Id(), err)
	}
	return resourceByteplusRouteTableRead(d, meta)
}

func resourceByteplusRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	routeTableService := NewRouteTableService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(routeTableService, d, ResourceByteplusRouteTable()); err != nil {
		return fmt.Errorf("error on reading route table %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	routeTableService := NewRouteTableService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Update(routeTableService, d, ResourceByteplusRouteTable()); err != nil {
		return fmt.Errorf("error on updating route table %q, %w", d.Id(), err)
	}
	return resourceByteplusRouteTableRead(d, meta)
}

func resourceByteplusRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	routeTableService := NewRouteTableService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(routeTableService, d, ResourceByteplusRouteTable()); err != nil {
		return fmt.Errorf("error on deleting route table %q, %w", d.Id(), err)
	}
	return nil
}
