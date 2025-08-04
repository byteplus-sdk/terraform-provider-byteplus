package allow_list

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
mongodb allow list can be imported using the allowListId, e.g.
```
$ terraform import byteplus_mongodb_allow_list.default acl-d1fd76693bd54e658912e7337d5b****
```

*/

func ResourceByteplusMongoDBAllowList() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusMongoDBAllowListCreate,
		Read:   resourceByteplusMongoDBAllowListRead,
		Update: resourceByteplusMongoDBAllowListUpdate,
		Delete: resourceByteplusMongoDBAllowListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allow_list_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of allow list.",
			},
			"allow_list_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of allow list.",
			},
			"allow_list_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPv4"}, false),
				Description:  "The IP address type of allow list, valid value contains `IPv4`.",
			},
			"allow_list": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: mongoDBAllowListImportDiffSuppress,
				Description:      "IP address or IP address segment in CIDR format. Duplicate addresses are not allowed, multiple addresses should be separated by commas (,) in English.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the allow list.",
			},

			// computed fields
			"allow_list_ip_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of allow list IPs.",
			},
			"associated_instance_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of instances bound under the allow list.",
			},
			"associated_instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of associated instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance id that bound to the allow list.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance name that bound to the allow list.",
						},
						"vpc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC ID.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project name of the instance.",
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceByteplusMongoDBAllowListCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBAllowListService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusMongoDBAllowList())
	if err != nil {
		return fmt.Errorf("Error on creating mongodb allow list %q, %s ", d.Id(), err)
	}
	return resourceByteplusMongoDBAllowListRead(d, meta)
}

func resourceByteplusMongoDBAllowListUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBAllowListService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusMongoDBAllowList())
	if err != nil {
		return fmt.Errorf("Error on updating mongodb allow list %q, %s ", d.Id(), err)
	}
	return resourceByteplusMongoDBAllowListRead(d, meta)
}

func resourceByteplusMongoDBAllowListDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBAllowListService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusMongoDBAllowList())
	if err != nil {
		return fmt.Errorf("Error on deleting mongodb allow list %q, %s ", d.Id(), err)
	}
	return err
}

func resourceByteplusMongoDBAllowListRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBAllowListService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusMongoDBAllowList())
	if err != nil {
		return fmt.Errorf("Error on reading mongodb allow list %q, %s ", d.Id(), err)
	}
	return err
}

func mongoDBAllowListImportDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if len(old) != len(new) {
		return false
	}
	oldAllowLists := strings.Split(old, ",")
	newAllowLists := strings.Split(new, ",")
	sort.Strings(oldAllowLists)
	sort.Strings(newAllowLists)
	// 根据前后allow list是否相等判断是否modify
	return reflect.DeepEqual(oldAllowLists, newAllowLists)
}
