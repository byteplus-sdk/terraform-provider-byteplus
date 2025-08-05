package allow_list_associate

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
mongodb allow list associate can be imported using the instanceId:allowListId, e.g.
```
$ terraform import byteplus_mongodb_allow_list_associate.default mongo-replica-e405f8e2****:acl-d1fd76693bd54e658912e7337d5b****
```

*/

func ResourceByteplusMongodbAllowListAssociate() *schema.Resource {
	resource := &schema.Resource{
		Read:   resourceByteplusMongodbAllowListAssociateRead,
		Create: resourceByteplusMongodbAllowListAssociateCreate,
		Delete: resourceByteplusMongodbAllowListAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: mongodbAllowListAssociateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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

func resourceByteplusMongodbAllowListAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongodbAllowListAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusMongodbAllowListAssociate())
	if err != nil {
		return fmt.Errorf("error on reading mongodb allow list association %v, %v", d.Id(), err)
	}
	return err
}

func resourceByteplusMongodbAllowListAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongodbAllowListAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusMongodbAllowListAssociate())
	if err != nil {
		return fmt.Errorf("error on creating mongodb allow list association %v, %v", d.Id(), err)
	}
	return err
}

func resourceByteplusMongodbAllowListAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongodbAllowListAssociateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusMongodbAllowListAssociate())
	if err != nil {
		return fmt.Errorf("error on deleting mongodb allow list association %v, %v", d.Id(), err)
	}
	return err
}
