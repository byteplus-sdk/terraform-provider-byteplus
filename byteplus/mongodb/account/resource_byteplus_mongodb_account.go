package account

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"sort"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
MongodbAccount can be imported using the instance_id:account_name, e.g.
```
$ terraform import byteplus_mongodb_account.default resource_id
```

*/

func ResourceByteplusMongoDBAccount() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusMongodbAccountCreate,
		Read:   resourceByteplusMongodbAccountRead,
		Update: resourceByteplusMongodbAccountUpdate,
		Delete: resourceByteplusMongodbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
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
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the mongodb instance.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the mongodb account.",
			},
			"account_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password of the mongodb account. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"auth_db": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The database of the mongodb account.",
			},
			"account_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the mongodb account.",
			},
			"account_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Set:         MongoDBAccountPrivilegeHash,
				Description: "The privilege information of account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of database.",
						},
						"role_names": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The role names of the account.",
						},
					},
				},
			},

			// computed fields
			"account_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the account.",
			},
		},
	}
	return resource
}

func resourceByteplusMongodbAccountCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBAccountService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusMongoDBAccount())
	if err != nil {
		return fmt.Errorf("error on creating mongodb_account %q, %s", d.Id(), err)
	}
	return resourceByteplusMongodbAccountRead(d, meta)
}

func resourceByteplusMongodbAccountRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBAccountService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusMongoDBAccount())
	if err != nil {
		return fmt.Errorf("error on reading mongodb_account %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusMongodbAccountUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBAccountService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusMongoDBAccount())
	if err != nil {
		return fmt.Errorf("error on updating mongodb_account %q, %s", d.Id(), err)
	}
	return resourceByteplusMongodbAccountRead(d, meta)
}

func resourceByteplusMongodbAccountDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBAccountService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusMongoDBAccount())
	if err != nil {
		return fmt.Errorf("error on deleting mongodb_account %q, %s", d.Id(), err)
	}
	return err
}

func mongodbAccountPrivilegeHashBase(m map[string]interface{}) (buf bytes.Buffer) {
	dbName := strings.ToLower(m["db_name"].(string))
	roleNames := m["role_names"].(*schema.Set).List()
	roleNamesArr := make([]string, 0)
	for _, v := range roleNames {
		roleNamesArr = append(roleNamesArr, v.(string))
	}
	sort.Strings(roleNamesArr)
	roleStr := strings.Join(roleNamesArr, ",")

	buf.WriteString(fmt.Sprintf("%s-", dbName))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(roleStr)))
	return buf
}

func MongoDBAccountPrivilegeHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	buf := mongodbAccountPrivilegeHashBase(m)
	return hashcode.String(buf.String())
}
