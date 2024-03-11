package kubeconfig

import (
	"fmt"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

/*

Import
VkeKubeconfig can be imported using the id, e.g.
```
$ terraform import byteplus_vke_kubeconfig.default kce8simvqtofl0l6u4qd0
```

*/

func ResourceByteplusVkeKubeconfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusVkeKubeconfigCreate,
		Read:   resourceByteplusVkeKubeconfigRead,
		Delete: resourceByteplusVkeKubeconfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The cluster id of the Kubeconfig.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the Kubeconfig, the value of type should be Public or Private.",
			},
			"valid_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     26280,
				Description: "The ValidDuration of the Kubeconfig, the range of the ValidDuration is 1 hour to 43800 hour.",
			},
		},
	}
}

func resourceByteplusVkeKubeconfigCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kubeconfigService := NewVkeKubeconfigService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(kubeconfigService, d, ResourceByteplusVkeKubeconfig())
	if err != nil {
		return fmt.Errorf("error on creating cluster  %q, %w", d.Id(), err)
	}
	return resourceByteplusVkeKubeconfigRead(d, meta)
}

func resourceByteplusVkeKubeconfigRead(d *schema.ResourceData, meta interface{}) (err error) {
	kubeconfigService := NewVkeKubeconfigService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(kubeconfigService, d, ResourceByteplusVkeKubeconfig())
	if err != nil {
		return fmt.Errorf("error on reading cluster %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusVkeKubeconfigDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kubeconfigService := NewVkeKubeconfigService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(kubeconfigService, d, ResourceByteplusVkeKubeconfig())
	if err != nil {
		return fmt.Errorf("error on deleting cluster %q, %w", d.Id(), err)
	}
	return err
}
