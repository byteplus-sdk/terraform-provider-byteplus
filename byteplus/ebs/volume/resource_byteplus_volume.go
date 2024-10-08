package volume

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Volume can be imported using the id, e.g.
```
$ terraform import byteplus_volume.default vol-mizl7m1kqccg5smt1bdpijuj
```
*/

func ResourceByteplusVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusVolumeCreate,
		Read:   resourceByteplusVolumeRead,
		Update: resourceByteplusVolumeUpdate,
		Delete: resourceByteplusVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the Zone.",
			},
			"volume_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of Volume.",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of Volume, the value is `PTSSD` or `ESSD_PL0` or `ESSD_PL1` or `ESSD_PL2` or `ESSD_FlexPL`.",
			},
			"kind": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The kind of Volume, the value is `data`.",
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The size of Volume.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Volume.",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: "The ID of the instance to which the created volume is automatically attached. " +
					"Please note this field needs to ask the system administrator to apply for a whitelist.\n" +
					"When use this field to attach ecs instance, the attached volume cannot be deleted by terraform, please use `terraform state rm byteplus_volume.resource_name` command to remove it from terraform state file and management.",
			},
			"volume_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "PostPaid",
				Description: "The charge type of the Volume, the value is `PostPaid` or `PrePaid`. " +
					"The `PrePaid` volume cannot be detached. " +
					"Please note that `PrePaid` type needs to ask the system administrator to apply for a whitelist.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of Volume.",
			},
			"trade_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of Trade.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of Volume.",
			},
			"delete_with_instance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Delete Volume with Attached Instance.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ProjectName of the Volume.",
			},
			"tags": bp.TagsSchema(),
		},
	}
}

func resourceByteplusVolumeCreate(d *schema.ResourceData, meta interface{}) (err error) {
	volumeService := NewVolumeService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(volumeService, d, ResourceByteplusVolume())
	if err != nil {
		return fmt.Errorf("error on creating volume %q, %w", d.Id(), err)
	}
	return resourceByteplusVolumeRead(d, meta)
}

func resourceByteplusVolumeRead(d *schema.ResourceData, meta interface{}) (err error) {
	volumeService := NewVolumeService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(volumeService, d, ResourceByteplusVolume())
	if err != nil {
		return fmt.Errorf("error on reading volume %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusVolumeUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	volumeService := NewVolumeService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(volumeService, d, ResourceByteplusVolume())
	if err != nil {
		return fmt.Errorf("error on updating volume %q, %w", d.Id(), err)
	}
	return resourceByteplusVolumeRead(d, meta)
}

func resourceByteplusVolumeDelete(d *schema.ResourceData, meta interface{}) (err error) {
	volumeService := NewVolumeService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(volumeService, d, ResourceByteplusVolume())
	if err != nil {
		return fmt.Errorf("error on deleting volume %q, %w", d.Id(), err)
	}
	return err
}
