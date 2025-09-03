package bucket_policy

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Tos Bucket can be imported using the id, e.g.
```
$ terraform import byteplus_tos_bucket_policy.default bucketName:policy
```

*/

func ResourceByteplusTosBucketPolicy() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusTosBucketPolicyCreate,
		Read:   resourceByteplusTosBucketPolicyRead,
		Update: resourceByteplusTosBucketPolicyUpdate,
		Delete: resourceByteplusTosBucketPolicyDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must be of the form bucketName")
				}
				_ = data.Set("bucket_name", items[0])
				return []*schema.ResourceData{data}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the bucket.",
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if k == "policy" && d.Id() != "" {
						om := map[string]interface{}{}
						nm := map[string]interface{}{}
						_ = json.Unmarshal([]byte(old), &om)
						_ = json.Unmarshal([]byte(new), &nm)
						//暂时不支持version 这里忽略掉
						delete(om, "Version")
						delete(nm, "Version")
						//重构单一Principal Action Resource 字符串转换数组
						if _, ok := om["Statement"].([]interface{}); ok {
							for i, st := range om["Statement"].([]interface{}) {
								if _, ok1 := st.(map[string]interface{}); ok1 {
									temp := map[string]interface{}{}
									for k1, v1 := range st.(map[string]interface{}) {
										if k1 == "Principal" || k1 == "Action" || k1 == "Resource" {
											if reflect.TypeOf(v1).Kind() == reflect.String {
												temp[k1] = []string{v1.(string)}
											} else {
												temp[k1] = v1
											}
										} else {
											temp[k1] = v1
										}
									}
									om["Statement"].([]interface{})[i] = temp
								}
							}
						}

						if _, ok := nm["Statement"].([]interface{}); ok {
							for i, st := range nm["Statement"].([]interface{}) {
								if _, ok1 := st.(map[string]interface{}); ok1 {
									temp := map[string]interface{}{}
									for k1, v1 := range st.(map[string]interface{}) {
										if k1 == "Principal" || k1 == "Action" || k1 == "Resource" {
											if reflect.TypeOf(v1).Kind() == reflect.String {
												temp[k1] = []string{v1.(string)}
											} else {
												temp[k1] = v1
											}
										} else {
											temp[k1] = v1
										}
									}
									nm["Statement"].([]interface{})[i] = temp
								}
							}
						}

						o, _ := json.MarshalIndent(om, "", "")
						n, _ := json.MarshalIndent(nm, "", "")
						return string(o) == string(n)
					}
					return false
				},
				Description: "The policy document. This is a JSON formatted string. For more information about building Byteplus IAM policy documents with Terraform, see the  [Byteplus IAM Policy Document Guide](https://www.byteplus.com/docs/6349/102127).",
			},
		},
	}
	return resource
}

func resourceByteplusTosBucketPolicyCreate(d *schema.ResourceData, meta interface{}) (err error) {
	tosBucketService := NewTosBucketPolicyService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(tosBucketService, d, ResourceByteplusTosBucketPolicy())
	if err != nil {
		return fmt.Errorf("error on creating tos bucket policy %q, %s", d.Id(), err)
	}
	return resourceByteplusTosBucketPolicyRead(d, meta)
}

func resourceByteplusTosBucketPolicyRead(d *schema.ResourceData, meta interface{}) (err error) {
	tosBucketService := NewTosBucketPolicyService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(tosBucketService, d, ResourceByteplusTosBucketPolicy())
	if err != nil {
		return fmt.Errorf("error on reading tos bucket policy %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusTosBucketPolicyUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	tosBucketService := NewTosBucketPolicyService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(tosBucketService, d, ResourceByteplusTosBucketPolicy())
	if err != nil {
		return fmt.Errorf("error on updating tos bucket policy %q, %s", d.Id(), err)
	}
	return resourceByteplusTosBucketPolicyRead(d, meta)
}

func resourceByteplusTosBucketPolicyDelete(d *schema.ResourceData, meta interface{}) (err error) {
	tosBucketService := NewTosBucketPolicyService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(tosBucketService, d, ResourceByteplusTosBucketPolicy())
	if err != nil {
		return fmt.Errorf("error on deleting tos bucket policy %q, %s", d.Id(), err)
	}
	return err
}
