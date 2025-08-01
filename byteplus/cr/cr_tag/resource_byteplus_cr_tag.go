package cr_tag

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CR tags can be imported using the registry:namespace:repository:tag, e.g.
```
$ terraform import byteplus_cr_tag.default cr-basic:namespace-1:repo-1:v1
```

*/

func CrTagImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 4 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must be 'registry:namespace:repository:tag'")
	}
	if err := d.Set("registry", items[0]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	if err := d.Set("namespace", items[1]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	if err := d.Set("repository", items[2]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	if err := d.Set("name", items[3]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func ResourceByteplusCrTag() *schema.Resource {
	resource := &schema.Resource{
		Read:   resourceByteplusCrTagRead,
		Create: resourceByteplusCrTagCreate,
		Delete: resourceByteplusCrTagDelete,
		Importer: &schema.ResourceImporter{
			State: CrTagImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The CrRegistry name.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target namespace name.",
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of repository.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of OCI product.",
			},
		},
	}
	dataSource := DataSourceByteplusCrTags().Schema["tags"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusCrTagCreate(d *schema.ResourceData, meta interface{}) (err error) {
	return fmt.Errorf("cr tag only support import")
}

func resourceByteplusCrTagDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrTagService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCrTag())
	if err != nil {
		return fmt.Errorf("error on deleting cr tag %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCrTagRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrTagService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCrTag())
	if err != nil {
		return fmt.Errorf("Error on reading cr tag %q,%s", d.Id(), err)
	}
	return err
}
