package network_acl_associate

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
NetworkAcl associate can be imported using the network_acl_id:resource_id, e.g.
```
$ terraform import byteplus_network_acl_associate.default nacl-172leak37mi9s4d1w33pswqkh:subnet-637jxq81u5mon3gd6ivc7rj
```

*/

func ResourceByteplusNetworkAclAssociate() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusAclAssociateCreate,
		Read:   resourceByteplusAclAssociateRead,
		Delete: resourceByteplusAclAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: aclAssociateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"network_acl_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of Network Acl.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The resource id of Network Acl.",
			},
		},
	}
}

func resourceByteplusAclAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	aclAssociateService := NewNetworkAclAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(aclAssociateService, d, ResourceByteplusNetworkAclAssociate())
	if err != nil {
		return fmt.Errorf("error on creating acl Associate %q, %w", d.Id(), err)
	}
	return resourceByteplusAclAssociateRead(d, meta)
}

func resourceByteplusAclAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	aclAssociateService := NewNetworkAclAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(aclAssociateService, d, ResourceByteplusNetworkAclAssociate())
	if err != nil {
		return fmt.Errorf("error on reading acl Associate %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusAclAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	aclAssociateService := NewNetworkAclAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(aclAssociateService, d, ResourceByteplusNetworkAclAssociate())
	if err != nil {
		return fmt.Errorf("error on deleting acl Associate %q, %w", d.Id(), err)
	}
	return err
}
