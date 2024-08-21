package cdn_kv

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnKv can be imported using the namespace_id:namespace:key, e.g.
```
$ terraform import byteplus_cdn_kv.default namespace_id:namespace:key
```

*/

func ResourceByteplusCdnKv() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnKvCreate,
		Read:   resourceByteplusCdnKvRead,
		Update: resourceByteplusCdnKvUpdate,
		Delete: resourceByteplusCdnKvDelete,
		Importer: &schema.ResourceImporter{
			State: cdnKvImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the kv namespace.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the kv namespace.",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The key of the kv namespace.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of the kv namespace key. Single Value upload data does not exceed 128KB. This field must be encrypted with base64.",
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Set the data storage time. Unit: second. After the data expires, the Value in the Key will be inaccessible.\nIf this parameter is not specified or the parameter value is 0, it is stored permanently by default.\nThe storage time cannot be less than 60s.\n" +
					"When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},

			// computed fields
			"key_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The status of the kv namespace key.",
			},
			"ddl": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Data expiration time. After the data expires, the Value in the Key will be inaccessible.\nDisplayed in UNIX timestamp format.\n0: Permanent storage.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The creation time of the kv namespace key. Displayed in UNIX timestamp format.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update time of the kv namespace key. Displayed in UNIX timestamp format.",
			},
		},
	}
	return resource
}

func resourceByteplusCdnKvCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnKvService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnKv())
	if err != nil {
		return fmt.Errorf("error on creating cdn_kv %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnKvRead(d, meta)
}

func resourceByteplusCdnKvRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnKvService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnKv())
	if err != nil {
		return fmt.Errorf("error on reading cdn_kv %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnKvUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnKvService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnKv())
	if err != nil {
		return fmt.Errorf("error on updating cdn_kv %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnKvRead(d, meta)
}

func resourceByteplusCdnKvDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnKvService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnKv())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_kv %q, %s", d.Id(), err)
	}
	return err
}

var cdnKvImporter = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 3 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("namespace_id", items[0]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	if err := data.Set("namespace", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	if err := data.Set("key", items[2]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
