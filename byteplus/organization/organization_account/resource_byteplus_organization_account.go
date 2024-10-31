package organization_account

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
OrganizationAccount can be imported using the id, e.g.
```
$ terraform import byteplus_organization_account.default resource_id
```

*/

func ResourceByteplusOrganizationAccount() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusOrganizationAccountCreate,
		Read:   resourceByteplusOrganizationAccountRead,
		Update: resourceByteplusOrganizationAccountUpdate,
		Delete: resourceByteplusOrganizationAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account.",
			},
			"show_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The show name of the account.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the account.",
			},
			"org_unit_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of the organization unit. Default is root organization.",
			},
			"verification_relation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The verification relation id of the account. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"tags": bp.TagsSchema(),

			// computed fields
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner id of the account.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization.",
			},
			"org_unit_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the organization unit.",
			},
			"org_verification_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the organization verification.",
			},
			"iam_role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the iam role.",
			},
		},
	}
	return resource
}

func resourceByteplusOrganizationAccountCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewOrganizationAccountService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusOrganizationAccount())
	if err != nil {
		return fmt.Errorf("error on creating organization_account %q, %s", d.Id(), err)
	}
	return resourceByteplusOrganizationAccountRead(d, meta)
}

func resourceByteplusOrganizationAccountRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewOrganizationAccountService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusOrganizationAccount())
	if err != nil {
		return fmt.Errorf("error on reading organization_account %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusOrganizationAccountUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewOrganizationAccountService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusOrganizationAccount())
	if err != nil {
		return fmt.Errorf("error on updating organization_account %q, %s", d.Id(), err)
	}
	return resourceByteplusOrganizationAccountRead(d, meta)
}

func resourceByteplusOrganizationAccountDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewOrganizationAccountService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusOrganizationAccount())
	if err != nil {
		return fmt.Errorf("error on deleting organization_account %q, %s", d.Id(), err)
	}
	return err
}
