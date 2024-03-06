package server_group_server

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
ServerGroupServer can be imported using the id, e.g.
```
$ terraform import byteplus_server_group_server.default rsp-274xltv2*****8tlv3q3s:rs-3ciynux6i1x4w****rszh49sj
```

*/

func ResourceByteplusServerGroupServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusServerGroupServerCreate,
		Read:   resourceByteplusServerGroupServerRead,
		Update: resourceByteplusServerGroupServerUpdate,
		Delete: resourceByteplusServerGroupServerDelete,
		Importer: &schema.ResourceImporter{
			State: serverGroupServerImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the ServerGroup.",
			},
			"server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server id of instance in ServerGroup.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of ecs instance or the network card bound to ecs instance.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The type of instance. Optional choice contains `ecs`, `eni`.",
				ValidateFunc: validation.StringInSlice([]string{"ecs", "eni"}, false),
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The weight of the instance, range in 0~100.",
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The private ip of the instance.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The port receiving request.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the instance.",
			},
		},
	}
}

func resourceByteplusServerGroupServerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupServerService := NewServerGroupServerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(serverGroupServerService, d, ResourceByteplusServerGroupServer())
	if err != nil {
		return fmt.Errorf("error on creating serverGroupServer  %q, %w", d.Id(), err)
	}
	return resourceByteplusServerGroupServerRead(d, meta)
}

func resourceByteplusServerGroupServerRead(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupServerService := NewServerGroupServerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(serverGroupServerService, d, ResourceByteplusServerGroupServer())
	if err != nil {
		return fmt.Errorf("error on reading serverGroupServer %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusServerGroupServerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupServerService := NewServerGroupServerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(serverGroupServerService, d, ResourceByteplusServerGroupServer())
	if err != nil {
		return fmt.Errorf("error on updating serverGroupServer  %q, %w", d.Id(), err)
	}
	return resourceByteplusServerGroupServerRead(d, meta)
}

func resourceByteplusServerGroupServerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupServerService := NewServerGroupServerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(serverGroupServerService, d, ResourceByteplusServerGroupServer())
	if err != nil {
		return fmt.Errorf("error on deleting serverGroupServer %q, %w", d.Id(), err)
	}
	return err
}
