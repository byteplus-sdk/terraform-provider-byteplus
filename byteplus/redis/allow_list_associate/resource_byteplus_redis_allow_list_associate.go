package allow_list_associate

import (
	"fmt"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Redis AllowList Association can be imported using the instanceId:allowListId, e.g.
```
$ terraform import byteplus_redis_allow_list_associate.default redis-asdljioeixxxx:acl-cn03wk541s55c376xxxx
```
*/

func ResourceByteplusRedisAllowListAssociate() *schema.Resource {
	resource := &schema.Resource{
		Read:   resourceByteplusRedisAllowListAssociateRead,
		Create: resourceByteplusRedisAllowListAssociateCreate,
		Delete: resourceByteplusRedisAllowListAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: redisAllowListAssociateImporter,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of instance to associate.",
			},
			"allow_list_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of allow list to associate.",
			},
		},
	}
	return resource
}

func resourceByteplusRedisAllowListAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	redisAllowListAssociateService := NewRedisAllowListAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(redisAllowListAssociateService, d, ResourceByteplusRedisAllowListAssociate())
	if err != nil {
		return fmt.Errorf("error on reading association %v, %v", d.Id(), err)
	}
	return err
}

func resourceByteplusRedisAllowListAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	redisAllowListAssociateService := NewRedisAllowListAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(redisAllowListAssociateService, d, ResourceByteplusRedisAllowListAssociate())
	if err != nil {
		return fmt.Errorf("error on creating redis allow list association %v, %v", d.Id(), err)
	}
	return err
}

func resourceByteplusRedisAllowListAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	redisAllowListAssociateService := NewRedisAllowListAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(redisAllowListAssociateService, d, ResourceByteplusRedisAllowListAssociate())
	if err != nil {
		return fmt.Errorf("error on deleting redis allow list association %v, %v", d.Id(), err)
	}
	return err
}
