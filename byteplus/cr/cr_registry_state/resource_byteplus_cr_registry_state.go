package cr_registry_state

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
CR registry state can be imported using the state:registry_name, e.g.
```
$ terraform import byteplus_cr_registry.default state:cr-basic
```

*/

func crRegistryStateImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must be 'state:registryName'")
	}
	if err := d.Set("name", items[1]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func ResourceByteplusCrRegistryState() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusCrRegistryStateCreate,
		Update: resourceByteplusCrRegistryStateUpdate,
		Read:   resourceByteplusCrRegistryStateRead,
		Delete: resourceByteplusCrRegistryStateDelete,
		Importer: &schema.ResourceImporter{
			State: crRegistryStateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return false
				},
				ValidateFunc: validation.StringInSlice([]string{"Start"}, false),
				Description:  "Start cr instance action,the value must be `Start`.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The cr instance id.",
			},
			"status": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Description: "The status of cr instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"phase": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The phase status of instance.",
						},
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The condition of instance.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceByteplusCrRegistryStateCreate(d *schema.ResourceData, meta interface{}) error {
	service := NewCrRegistryStateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(service, d, ResourceByteplusCrRegistryState()); err != nil {
		return fmt.Errorf("error on creating instance state %q, %w", d.Id(), err)
	}
	return resourceByteplusCrRegistryStateRead(d, meta)
}

func resourceByteplusCrRegistryStateUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("this resource does not allow update operation")
}

func resourceByteplusCrRegistryStateRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCrRegistryStateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(service, d, ResourceByteplusCrRegistryState()); err != nil {
		return fmt.Errorf("error on reading instance state %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusCrRegistryStateDelete(d *schema.ResourceData, meta interface{}) error {
	service := NewCrRegistryStateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCrRegistryState()); err != nil {
		return fmt.Errorf("error on deleting instance state %q, %w", d.Id(), err)
	}
	return nil
}
