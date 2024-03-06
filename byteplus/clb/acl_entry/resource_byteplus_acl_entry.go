package acl_entry

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
AclEntry can be imported using the id, e.g.
```
$ terraform import byteplus_acl_entry.default ID is a string concatenated with colons(AclId:Entry)
```

*/

func ResourceByteplusAclEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusAclEntryCreate,
		Read:   resourceByteplusAclEntryRead,
		Delete: resourceByteplusAclEntryDelete,
		Importer: &schema.ResourceImporter{
			State: aclEntryImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of Acl.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The description of the AclEntry.",
			},
			"entry": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The content of the AclEntry.",
			},
		},
	}
}

func resourceByteplusAclEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	aclEntryService := NewAclEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(aclEntryService, d, ResourceByteplusAclEntry())
	if err != nil {
		return fmt.Errorf("error on creating acl entry %q, %w", d.Id(), err)
	}
	return resourceByteplusAclEntryRead(d, meta)
}

func resourceByteplusAclEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	aclEntryService := NewAclEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(aclEntryService, d, ResourceByteplusAclEntry())
	if err != nil {
		return fmt.Errorf("error on reading acl entry %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusAclEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	aclEntryService := NewAclEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(aclEntryService, d, ResourceByteplusAclEntry())
	if err != nil {
		return fmt.Errorf("error on deleting acl entry %q, %w", d.Id(), err)
	}
	return err
}
