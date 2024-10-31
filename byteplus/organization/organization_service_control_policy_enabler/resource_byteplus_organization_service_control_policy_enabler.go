package organization_service_control_policy_enabler

import (
	"fmt"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
ServiceControlPolicy enabler can be imported using the default_id (organization:service_control_policy_enable) , e.g.
```
$ terraform import byteplus_organization_service_control_policy_enabler.default organization:service_control_policy_enable
```

*/

func ResourceByteplusOrganizationServiceControlPolicyEnabler() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusOrganizationServiceControlPolicyEnablerCreate,
		Read:   resourceByteplusOrganizationServiceControlPolicyEnablerRead,
		Delete: resourceByteplusOrganizationServiceControlPolicyEnablerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{},
	}
	return resource
}

func resourceByteplusOrganizationServiceControlPolicyEnablerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusOrganizationServiceControlPolicyEnabler())
	if err != nil {
		return fmt.Errorf("error on creating organization_service_control_policy_enabler: %q, %s", d.Id(), err)
	}
	return resourceByteplusOrganizationServiceControlPolicyEnablerRead(d, meta)
}

func resourceByteplusOrganizationServiceControlPolicyEnablerRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusOrganizationServiceControlPolicyEnabler())
	if err != nil {
		return fmt.Errorf("error on reading organization_service_control_policy_enabler: %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusOrganizationServiceControlPolicyEnablerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusOrganizationServiceControlPolicyEnabler())
	if err != nil {
		return fmt.Errorf("erron on deleting organization_service_control_policy_enabler: %q, %s", d.Id(), err)
	}
	return err
}
