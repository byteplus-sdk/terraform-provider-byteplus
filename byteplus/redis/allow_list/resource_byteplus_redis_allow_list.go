package allow_list

import (
	"fmt"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Redis AllowList can be imported using the id, e.g.
```
$ terraform import byteplus_redis_allow_list.default acl-cn03wk541s55c376xxxx
```

*/

func ResourceByteplusRedisAllowList() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusRedisAllowListCreate,
		Read:   resourceByteplusRedisAllowListRead,
		Update: resourceByteplusRedisAllowListUpdate,
		Delete: resourceByteplusRedisAllowListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"allow_list_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of allow list.",
			},
			"allow_list_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of allow list.",
			},
			"allow_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "Ip list of allow list.",
			},
		},
	}
	bp.MergeDateSourceToResource(DataSourceByteplusRedisAllowLists().Schema["allow_lists"].Elem.(*schema.Resource).Schema, &resource.Schema)
	return resource
}

func resourceByteplusRedisAllowListCreate(d *schema.ResourceData, meta interface{}) (err error) {
	redisAllowListService := NewRedisAllowListService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(redisAllowListService, d, ResourceByteplusRedisAllowList())
	if err != nil {
		return fmt.Errorf("error on creating redis allowlist %v, %v", d.Id(), err)
	}
	return resourceByteplusRedisAllowListRead(d, meta)
}

func resourceByteplusRedisAllowListUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	redisAllowListService := NewRedisAllowListService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(redisAllowListService, d, ResourceByteplusRedisAllowList())
	if err != nil {
		return fmt.Errorf("error on updating redis allowlist  %q, %s", d.Id(), err)
	}
	return resourceByteplusRedisAllowListRead(d, meta)
}

func resourceByteplusRedisAllowListDelete(d *schema.ResourceData, meta interface{}) (err error) {
	redisAllowListService := NewRedisAllowListService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(redisAllowListService, d, ResourceByteplusRedisAllowList())
	if err != nil {
		return fmt.Errorf("error on deleting redis allowlist %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusRedisAllowListRead(d *schema.ResourceData, meta interface{}) (err error) {
	redisAllowListService := NewRedisAllowListService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(redisAllowListService, d, ResourceByteplusRedisAllowList())
	if err != nil {
		return fmt.Errorf("error on reading redis allowlist %q,%s", d.Id(), err)
	}
	return err
}
