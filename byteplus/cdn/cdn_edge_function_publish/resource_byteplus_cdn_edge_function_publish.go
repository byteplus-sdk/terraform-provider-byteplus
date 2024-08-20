package cdn_edge_function_publish

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnEdgeFunctionPublish can be imported using the function_id:ticket_id, e.g.
```
$ terraform import byteplus_cdn_edge_function_publish.default function_id:ticket_id
```

*/

func ResourceByteplusCdnEdgeFunctionPublish() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnEdgeFunctionPublishCreate,
		Read:   resourceByteplusCdnEdgeFunctionPublishRead,
		Delete: resourceByteplusCdnEdgeFunctionPublishDelete,
		Importer: &schema.ResourceImporter{
			State: functionPublishImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the function to which you want publish.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The description for this release.",
			},
			"publish_action": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"FullPublish", "CanaryPublish", "SnapshotPublish"}, false),
				Description: "The publish action of the edge function. Valid values: `FullPublish`, `CanaryPublish`, `SnapshotPublish`.\n" +
					"When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"publish_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{100, 200}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					action := d.Get("publish_action")
					return action.(string) != "SnapshotPublish"
				},
				Description: "The publish type of SnapshotPublish：\n200: FullPublish\n100: CanaryPublish. This field is required and valid when the `publish_action` is `SnapshotPublish`.\n" +
					"When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"version_tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					action := d.Get("publish_action")
					return action.(string) != "SnapshotPublish"
				},
				Description: "The specified version number to be published. This field is required and valid when the `publish_action` is `SnapshotPublish`.\n " +
					"When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},

			// computed fields
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content of the release record.",
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the release record.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The create time of the release record. Displayed in UNIX timestamp format.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update time of the release record. Displayed in UNIX timestamp format.",
			},
		},
	}
	return resource
}

func resourceByteplusCdnEdgeFunctionPublishCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionPublishService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnEdgeFunctionPublish())
	if err != nil {
		return fmt.Errorf("error on creating cdn_edge_function_publish %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnEdgeFunctionPublishRead(d, meta)
}

func resourceByteplusCdnEdgeFunctionPublishRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionPublishService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnEdgeFunctionPublish())
	if err != nil {
		return fmt.Errorf("error on reading cdn_edge_function_publish %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnEdgeFunctionPublishUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionPublishService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnEdgeFunctionPublish())
	if err != nil {
		return fmt.Errorf("error on updating cdn_edge_function_publish %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnEdgeFunctionPublishRead(d, meta)
}

func resourceByteplusCdnEdgeFunctionPublishDelete(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[DEBUG] deleting a byteplus_cdn_edge_function_publish resource will only remove the publish record from terraform state.")
	service := NewCdnEdgeFunctionPublishService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnEdgeFunctionPublish())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_edge_function_publish %q, %s", d.Id(), err)
	}
	return err
}

var functionPublishImporter = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("function_id", items[0]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
