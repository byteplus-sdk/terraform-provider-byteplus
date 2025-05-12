package alb_server_group_server

import (
	"fmt"
	"strings"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
AlbServerGroupServer can be imported using the server_group_id:server_id, e.g.
```
$ terraform import byteplus_alb_server_group_server.default rsp-274xltv2*****8tlv3q3s:rs-3ciynux6i1x4w****rszh49sj
```

*/

func ResourceByteplusAlbServerGroupServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusServerGroupServerCreate,
		Read:   resourceByteplusServerGroupServerRead,
		Update: resourceByteplusServerGroupServerUpdate,
		Delete: resourceByteplusServerGroupServerDelete,
		Importer: &schema.ResourceImporter{
			State: serverGroupServerImporter,
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of instance. Optional choice contains `ecs`, `eni`.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The weight of the instance, range in 0~100.",
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
	err = serverGroupServerService.Dispatcher.Create(serverGroupServerService, d, ResourceByteplusAlbServerGroupServer())
	if err != nil {
		return fmt.Errorf("error on creating serverGroupServer  %q, %w", d.Id(), err)
	}
	return resourceByteplusServerGroupServerRead(d, meta)
}

func resourceByteplusServerGroupServerRead(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupServerService := NewServerGroupServerService(meta.(*bp.SdkClient))
	err = serverGroupServerService.Dispatcher.Read(serverGroupServerService, d, ResourceByteplusAlbServerGroupServer())
	if err != nil {
		return fmt.Errorf("error on reading serverGroupServer %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusServerGroupServerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupServerService := NewServerGroupServerService(meta.(*bp.SdkClient))
	err = serverGroupServerService.Dispatcher.Update(serverGroupServerService, d, ResourceByteplusAlbServerGroupServer())
	if err != nil {
		return fmt.Errorf("error on updating serverGroupServer  %q, %w", d.Id(), err)
	}
	return resourceByteplusServerGroupServerRead(d, meta)
}

func resourceByteplusServerGroupServerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	serverGroupServerService := NewServerGroupServerService(meta.(*bp.SdkClient))
	err = serverGroupServerService.Dispatcher.Delete(serverGroupServerService, d, ResourceByteplusAlbServerGroupServer())
	if err != nil {
		return fmt.Errorf("error on deleting serverGroupServer %q, %w", d.Id(), err)
	}
	return err
}

var serverGroupServerImporter = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("server_group_id", items[0]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	if err := data.Set("server_id", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
