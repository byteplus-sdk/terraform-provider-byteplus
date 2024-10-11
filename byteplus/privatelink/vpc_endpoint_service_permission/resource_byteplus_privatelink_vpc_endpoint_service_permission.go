package vpc_endpoint_service_permission

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
VpcEndpointServicePermission can be imported using the serviceId:permitAccountId, e.g.
```
$ terraform import byteplus_privatelink_vpc_endpoint_service_permission.default epsvc-2fe630gurkl37k5gfuy33****:2100000000
```

*/

func ResourceByteplusPrivatelinkVpcEndpointServicePermission() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusPrivatelinkVpcEndpointServicePermissionCreate,
		Read:   resourceByteplusPrivatelinkVpcEndpointServicePermissionRead,
		Delete: resourceByteplusPrivatelinkVpcEndpointServicePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
				}
				if err := data.Set("service_id", items[0]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				if err := data.Set("permit_account_id", items[1]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				return []*schema.ResourceData{data}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"permit_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of account.",
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of service.",
			},
		},
	}
	return resource
}

func resourceByteplusPrivatelinkVpcEndpointServicePermissionCreate(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(aclService, d, ResourceByteplusPrivatelinkVpcEndpointServicePermission())
	if err != nil {
		return fmt.Errorf("error on creating vpc endpoint service permission %q, %w", d.Id(), err)
	}
	return resourceByteplusPrivatelinkVpcEndpointServicePermissionRead(d, meta)
}

func resourceByteplusPrivatelinkVpcEndpointServicePermissionRead(d *schema.ResourceData, meta interface{}) (err error) {
	aclService := NewService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(aclService, d, ResourceByteplusPrivatelinkVpcEndpointServicePermission())
	if err != nil {
		return fmt.Errorf("error on reading vpc endpoint service permission %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusPrivatelinkVpcEndpointServicePermissionDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusPrivatelinkVpcEndpointServicePermission())
	if err != nil {
		return fmt.Errorf("error on deleting vpc endpoint service permission %q, %w", d.Id(), err)
	}
	return err
}
