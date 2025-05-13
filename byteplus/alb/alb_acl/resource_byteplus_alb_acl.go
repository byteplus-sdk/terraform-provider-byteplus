package alb_acl

import (
	"fmt"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Acl can be imported using the id, e.g.
```
$ terraform import byteplus_alb_acl.default acl-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceByteplusAlbAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusAclCreate,
		Read:   resourceByteplusAclRead,
		Update: resourceByteplusAclUpdate,
		Delete: resourceByteplusAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The project name of the Acl.",
			},
			"acl_entries": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The acl entry set of the Acl.",
				Set:         AclEntryHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The content of the AclEntry.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the AclEntry.",
						},
					},
				},
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
	err = aclService.Dispatcher.Create(aclService, d, ResourceByteplusAlbAcl())
	if err != nil {
		return fmt.Errorf("error on creating acl %q, %w", d.Id(), err)
	}
	return resourceByteplusAclRead(d, meta)
}

func resourceByteplusAclRead(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*bp.SdkClient))
	err = aclService.Dispatcher.Read(aclService, d, ResourceByteplusAlbAcl())
	if err != nil {
		return fmt.Errorf("error on reading acl %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusAclUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*bp.SdkClient))
	err = aclService.Dispatcher.Update(aclService, d, ResourceByteplusAlbAcl())
	if err != nil {
		return fmt.Errorf("error on updating acl %q, %w", d.Id(), err)
	}
	return resourceByteplusAclRead(d, meta)
}

func resourceByteplusAclDelete(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewAclService(meta.(*bp.SdkClient))
	err = aclService.Dispatcher.Delete(aclService, d, ResourceByteplusAlbAcl())
	if err != nil {
		return fmt.Errorf("error on deleting acl %q, %w", d.Id(), err)
	}
	return err
}
