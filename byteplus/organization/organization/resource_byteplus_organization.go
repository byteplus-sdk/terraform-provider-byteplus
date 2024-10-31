package organization

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Organization can be imported using the id, e.g.
```
$ terraform import byteplus_organization.default resource_id
```

*/

func ResourceByteplusOrganization() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusOrganizationCreate,
		Read:   resourceByteplusOrganizationRead,
		Delete: resourceByteplusOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// computed fields
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner id of the organization.",
			},
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The type of the organization.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the organization.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the organization.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The status of the organization.",
			},
			"delete_uk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The delete uk of the organization.",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The account id of the organization owner.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account name of the organization owner.",
			},
			"main_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The main name of the organization owner.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The created time of the organization.",
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The updated time of the organization.",
			},
			"deleted_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deleted time of the organization.",
			},
		},
	}
	return resource
}

func resourceByteplusOrganizationCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewOrganizationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusOrganization())
	if err != nil {
		return fmt.Errorf("error on creating organization %q, %s", d.Id(), err)
	}
	return resourceByteplusOrganizationRead(d, meta)
}

func resourceByteplusOrganizationRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewOrganizationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusOrganization())
	if err != nil {
		return fmt.Errorf("error on reading organization %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusOrganizationDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewOrganizationService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusOrganization())
	if err != nil {
		return fmt.Errorf("error on deleting organization %q, %s", d.Id(), err)
	}
	return err
}
