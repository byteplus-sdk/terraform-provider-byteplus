package cdn_kv_namespace

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnKvNamespace can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_kv_namespace.default resource_id
```

*/

func ResourceByteplusCdnKvNamespace() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnKvNamespaceCreate,
		Read:   resourceByteplusCdnKvNamespaceRead,
		Update: resourceByteplusCdnKvNamespaceUpdate,
		Delete: resourceByteplusCdnKvNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Set a recognizable name for the namespace. The input requirements are as follows:\nLength should be between 2 and 64 characters.\nIt can only contain English letters, numbers, hyphens (-), and underscores (_).",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Set a description for the namespace. The input requirements are as follows:\nAny characters are allowed.\nThe length should not exceed 80 characters.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the project to which the namespace belongs, defaulting to `default`.",
			},
		},
	}
	return resource
}

func resourceByteplusCdnKvNamespaceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnKvNamespaceService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnKvNamespace())
	if err != nil {
		return fmt.Errorf("error on creating cdn_kv_namespace %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnKvNamespaceRead(d, meta)
}

func resourceByteplusCdnKvNamespaceRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnKvNamespaceService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnKvNamespace())
	if err != nil {
		return fmt.Errorf("error on reading cdn_kv_namespace %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnKvNamespaceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnKvNamespaceService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnKvNamespace())
	if err != nil {
		return fmt.Errorf("error on updating cdn_kv_namespace %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnKvNamespaceRead(d, meta)
}

func resourceByteplusCdnKvNamespaceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnKvNamespaceService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnKvNamespace())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_kv_namespace %q, %s", d.Id(), err)
	}
	return err
}
