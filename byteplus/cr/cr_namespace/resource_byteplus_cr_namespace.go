package cr_namespace

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CR namespace can be imported using the registry:name, e.g.
```
$ terraform import byteplus_cr_namespace.default cr-basic:namespace-1
```

*/

func crNamespaceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must be 'registry:namespace'")
	}
	if err := d.Set("registry", items[0]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	if err := d.Set("name", items[1]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func ResourceByteplusCrNamespace() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCrNamespaceCreate,
		Read:   resourceByteplusCrNamespaceRead,
		Update: resourceByteplusCrNamespaceUpdate,
		Delete: resourceByteplusCrNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: crNamespaceImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The registry name.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of CrNamespace.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ProjectName of the CrNamespace.",
			},
			"repository_default_access_level": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The default access level of repository. Valid values: `Private`, `Public`. Default is `Private`.",
			},
		},
	}
	dataSource := DataSourceByteplusCrNamespaces().Schema["namespaces"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusCrNamespaceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrNamespaceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusCrNamespace())
	if err != nil {
		return fmt.Errorf("error on creating CrNamespace %q,%s", d.Id(), err)
	}
	return resourceByteplusCrNamespaceRead(d, meta)
}

func resourceByteplusCrNamespaceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrNamespaceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusCrNamespace())
	if err != nil {
		return fmt.Errorf("error on updating CrNamespace  %q, %s", d.Id(), err)
	}
	return resourceByteplusCrNamespaceRead(d, meta)
}

func resourceByteplusCrNamespaceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrNamespaceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCrNamespace())
	if err != nil {
		return fmt.Errorf("error on deleting CrNamespace %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCrNamespaceRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrNamespaceService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCrNamespace())
	if err != nil {
		return fmt.Errorf("error on reading CrNamespace %q,%s", d.Id(), err)
	}
	return err
}
