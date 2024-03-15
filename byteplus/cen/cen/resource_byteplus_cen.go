package cen

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Cen can be imported using the id, e.g.
```
$ terraform import byteplus_cen.default cen-7qthudw0ll6jmc****
```

*/

func ResourceByteplusCen() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCenCreate,
		Read:   resourceByteplusCenRead,
		Update: resourceByteplusCenUpdate,
		Delete: resourceByteplusCenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"cen_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the cen.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the cen.",
			},
			"tags": bp.TagsSchema(),
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ProjectName of the cen instance.",
			},
		},
	}
	s := DataSourceByteplusCens().Schema["cens"].Elem.(*schema.Resource).Schema
	delete(s, "id")
	bp.MergeDateSourceToResource(s, &resource.Schema)
	return resource
}

func resourceByteplusCenCreate(d *schema.ResourceData, meta interface{}) (err error) {
	cenService := NewCenService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(cenService, d, ResourceByteplusCen())
	if err != nil {
		return fmt.Errorf("error on creating cen  %q, %s", d.Id(), err)
	}
	return resourceByteplusCenRead(d, meta)
}

func resourceByteplusCenRead(d *schema.ResourceData, meta interface{}) (err error) {
	cenService := NewCenService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(cenService, d, ResourceByteplusCen())
	if err != nil {
		return fmt.Errorf("error on reading cen %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCenUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	cenService := NewCenService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(cenService, d, ResourceByteplusCen())
	if err != nil {
		return fmt.Errorf("error on updating cen %q, %s", d.Id(), err)
	}
	return resourceByteplusCenRead(d, meta)
}

func resourceByteplusCenDelete(d *schema.ResourceData, meta interface{}) (err error) {
	cenService := NewCenService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(cenService, d, ResourceByteplusCen())
	if err != nil {
		return fmt.Errorf("error on deleting cen %q, %s", d.Id(), err)
	}
	return err
}
