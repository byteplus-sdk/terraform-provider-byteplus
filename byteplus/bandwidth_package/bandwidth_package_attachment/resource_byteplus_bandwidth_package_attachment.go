package bandwidth_package_attachment

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
BandwidthPackageAttachment can be imported using the bandwidth package id and eip id, e.g.
```
$ terraform import byteplus_bandwidth_package_attachment.default BandwidthPackageId:EipId
```

*/

func ResourceByteplusBandwidthPackageAttachment() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusBandwidthPackageAttachmentCreate,
		Read:   resourceByteplusBandwidthPackageAttachmentRead,
		Delete: resourceByteplusBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				var err error
				items := strings.Split(d.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{d}, fmt.Errorf("import id must be of the form BandwidthPackageId:EipId")
				}
				err = d.Set("bandwidth_package_id", items[0])
				if err != nil {
					return []*schema.ResourceData{d}, err
				}
				err = d.Set("allocation_id", items[1])
				if err != nil {
					return []*schema.ResourceData{d}, err
				}

				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The bandwidth package id.",
			},
			"allocation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the public IP or IPv6 public bandwidth to be added to the shared bandwidth package instance.",
			},
		},
	}
	return resource
}

func resourceByteplusBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewBandwidthPackageAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusBandwidthPackageAttachment())
	if err != nil {
		return fmt.Errorf("error on creating bandwidth_package_attachment %q, %s", d.Id(), err)
	}
	return resourceByteplusBandwidthPackageAttachmentRead(d, meta)
}

func resourceByteplusBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewBandwidthPackageAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusBandwidthPackageAttachment())
	if err != nil {
		return fmt.Errorf("error on reading bandwidth_package_attachment %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewBandwidthPackageAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusBandwidthPackageAttachment())
	if err != nil {
		return fmt.Errorf("error on deleting bandwidth_package_attachment %q, %s", d.Id(), err)
	}
	return err
}
