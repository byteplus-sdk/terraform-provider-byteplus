package account

import (
	"fmt"
	"strings"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Redis account can be imported using the instanceId:accountName, e.g.
```
$ terraform import byteplus_redis_account.default redis-42b38c769c4b:test
```

*/

var accountImporter = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("instance_id", items[0]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	if err := data.Set("account_name", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}

func ResourceByteplusRedisAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusRedisAccountCreate,
		Read:   resourceByteplusRedisAccountRead,
		Delete: resourceByteplusRedisAccountDelete,
		Update: resourceByteplusRedisAccountUpdate,
		Importer: &schema.ResourceImporter{
			State: accountImporter,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the Redis instance.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Redis account name.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password of the redis account. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the redis account.",
			},
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Role type, the valid value can be `Administrator`, `ReadWrite`, `ReadOnly`, `NotDangerous`.",
			},
		},
	}
}

func resourceByteplusRedisAccountCreate(d *schema.ResourceData, meta interface{}) (err error) {
	redisAccountService := NewAccountService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(redisAccountService, d, ResourceByteplusRedisAccount())
	if err != nil {
		return fmt.Errorf("error on creating redis account %q, %w", d.Id(), err)
	}
	return resourceByteplusRedisAccountRead(d, meta)
}

func resourceByteplusRedisAccountRead(d *schema.ResourceData, meta interface{}) (err error) {
	redisAccountService := NewAccountService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(redisAccountService, d, ResourceByteplusRedisAccount())
	if err != nil {
		return fmt.Errorf("error on reading redis account %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusRedisAccountUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	redisAccountService := NewAccountService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(redisAccountService, d, ResourceByteplusRedisAccount())
	if err != nil {
		return fmt.Errorf("error on update redis account %q, %w", d.Id(), err)
	}
	return err
}
func resourceByteplusRedisAccountDelete(d *schema.ResourceData, meta interface{}) (err error) {
	redisAccountService := NewAccountService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(redisAccountService, d, ResourceByteplusRedisAccount())
	if err != nil {
		return fmt.Errorf("error on deleting redis account %q, %w", d.Id(), err)
	}
	return err
}
