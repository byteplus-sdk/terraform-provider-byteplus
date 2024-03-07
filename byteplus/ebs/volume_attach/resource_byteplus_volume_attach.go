package volume_attach

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
VolumeAttach can be imported using the id, e.g.
```
$ terraform import byteplus_volume_attach.default vol-abc12345:i-abc12345
```

*/

func ResourceByteplusVolumeAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusVolumeAttachCreate,
		Read:   resourceByteplusVolumeAttachRead,
		Update: resourceByteplusVolumeAttachUpdate,
		Delete: resourceByteplusVolumeAttachDelete,
		Importer: &schema.ResourceImporter{
			State: importVolumeAttach,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Id of Volume.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Id of Instance.",
			},
			"delete_with_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Delete Volume with Attached Instance." +
					"It is not recommended to use this field. If used, please ensure that the value of this field is consistent with the value of `delete_with_instance` in byteplus_volume.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of Volume.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of Volume.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time of Volume.",
			},
		},
	}
}

func resourceByteplusVolumeAttachCreate(d *schema.ResourceData, meta interface{}) (err error) {
	volumeAttachService := NewVolumeAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(volumeAttachService, d, ResourceByteplusVolumeAttach())
	if err != nil {
		return fmt.Errorf("error on attach volume %q, %w", d.Id(), err)
	}
	return resourceByteplusVolumeAttachRead(d, meta)
}

func resourceByteplusVolumeAttachRead(d *schema.ResourceData, meta interface{}) (err error) {
	volumeAttachService := NewVolumeAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(volumeAttachService, d, ResourceByteplusVolumeAttach())
	if err != nil {
		return fmt.Errorf("error on reading volume %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusVolumeAttachUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return resourceByteplusVolumeAttachRead(d, meta)
}

func resourceByteplusVolumeAttachDelete(d *schema.ResourceData, meta interface{}) (err error) {
	volumeAttachService := NewVolumeAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(volumeAttachService, d, ResourceByteplusVolumeAttach())
	if err != nil {
		return fmt.Errorf("error on detach volume %q, %w", d.Id(), err)
	}
	return err
}
