package cr_repository

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
CR Repository can be imported using the registry:namespace:name, e.g.
```
$ terraform import byteplus_cr_repository.default cr-basic:namespace-1:repo-1
```

*/

func crRepositoryImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 3 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must be 'registry:namespace:name'")
	}
	if err := d.Set("registry", items[0]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	if err := d.Set("namespace", items[1]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	if err := d.Set("name", items[2]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func ResourceByteplusCrRepository() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCrRepositoryCreate,
		Read:   resourceByteplusCrRepositoryRead,
		Update: resourceByteplusCrRepositoryUpdate,
		Delete: resourceByteplusCrRepositoryDelete,
		Importer: &schema.ResourceImporter{
			State: crRepositoryImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of CrRepository.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of CrRepository.",
			},
			"access_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Private",
				ValidateFunc: validation.StringInSlice([]string{"Private", "Public"}, false),
				Description:  "The access level of CrRepository.",
			},
		},
	}
	dataSource := DataSourceByteplusCrRepositories().Schema["repositories"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusCrRepositoryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRepositoryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusCrRepository())
	if err != nil {
		return fmt.Errorf("error on creating CrRepository %q,%s", d.Id(), err)
	}
	return resourceByteplusCrRepositoryRead(d, meta)
}

func resourceByteplusCrRepositoryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRepositoryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusCrRepository())
	if err != nil {
		return fmt.Errorf("error on updating CrRepository  %q, %s", d.Id(), err)
	}
	return resourceByteplusCrRepositoryRead(d, meta)
}

func resourceByteplusCrRepositoryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRepositoryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCrRepository())
	if err != nil {
		return fmt.Errorf("error on deleting CrRepository %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCrRepositoryRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRepositoryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCrRepository())
	if err != nil {
		return fmt.Errorf("Error on reading CrRepository %q,%s", d.Id(), err)
	}
	return err
}
