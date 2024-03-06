package acl

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Acl can be imported using the id, e.g.
```
$ terraform import byteplus_acl.default acl-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceByteplusAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusAclCreate,
		Read:   resourceByteplusAclRead,
		Update: resourceByteplusAclUpdate,
		Delete: resourceByteplusAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of Acl.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Acl.",
			},
			"acl_entries": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The acl entry set of the Acl.",
				Set:         bp.ClbAclEntryHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the AclEntry.",
						},
						"entry": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The content of the AclEntry.",
						},
					},
				},
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ProjectName of the Acl.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of Acl.",
			},
		},
	}
}

func resourceByteplusAclCreate(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(aclService, d, ResourceByteplusAcl())
	if err != nil {
		return fmt.Errorf("error on creating acl %q, %w", d.Id(), err)
	}
	return resourceByteplusAclRead(d, meta)
}

func resourceByteplusAclRead(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(aclService, d, ResourceByteplusAcl())
	if err != nil {
		return fmt.Errorf("error on reading acl %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusAclUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(aclService, d, ResourceByteplusAcl())
	if err != nil {
		return fmt.Errorf("error on updating acl %q, %w", d.Id(), err)
	}
	return resourceByteplusAclRead(d, meta)
}

func resourceByteplusAclDelete(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(aclService, d, ResourceByteplusAcl())
	if err != nil {
		return fmt.Errorf("error on deleting acl %q, %w", d.Id(), err)
	}
	return err
}
