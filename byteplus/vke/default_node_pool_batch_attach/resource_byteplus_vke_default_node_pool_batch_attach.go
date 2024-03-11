package default_node_pool_batch_attach

import (
	"fmt"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/default_node_pool"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

The resource not support import

*/

func ResourceByteplusDefaultNodePoolBatchAttach() *schema.Resource {
	m := map[string]*schema.Schema{
		"cluster_id": default_node_pool.ResourceByteplusDefaultNodePool().Schema["cluster_id"],
		"default_node_pool_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The default NodePool ID.",
		},
		"instances": default_node_pool.ResourceByteplusDefaultNodePool().Schema["instances"],
		"kubernetes_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"labels": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:        schema.TypeString,
									Required:    true,
									ForceNew:    true,
									Description: "The Key of Labels.",
								},
								"value": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
									Description: "The Value of Labels.",
								},
							},
						},
						Description: "The Labels of KubernetesConfig.",
					},
					"taints": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:        schema.TypeString,
									Required:    true,
									ForceNew:    true,
									Description: "The Key of Taints.",
								},
								"value": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
									Description: "The Value of Taints.",
								},
								"effect": {
									Type:     schema.TypeString,
									Optional: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										"NoSchedule",
										"NoExecute",
										"PreferNoSchedule",
									}, false),
									Description: "The Effect of Taints. The value can be one of the following: `NoSchedule`, `NoExecute`, `PreferNoSchedule`, default value is `NoSchedule`.",
								},
							},
						},
						Description: "The Taints of KubernetesConfig.",
					},
					"cordon": {
						Type:        schema.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Description: "The Cordon of KubernetesConfig.",
					},
				},
			},
			Description: "The KubernetesConfig of NodeConfig. Please note that this field is the configuration of the node. The same key is subject to the config of the node pool. Different keys take effect together.",
		},
	}
	bp.MergeDateSourceToResource(default_node_pool.ResourceByteplusDefaultNodePool().Schema, &m)

	// logger.Debug(logger.RespFormat, "ATTACH_TEST", m)

	return &schema.Resource{
		Create: resourceByteplusDefaultNodePoolBatchAttachCreate,
		Update: resourceByteplusDefaultNodePoolBatchAttachUpdate,
		Read:   resourceByteplusDefaultNodePoolBatchAttachUpdate,
		Delete: resourceByteplusNodePoolBatchAttachDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				return nil, fmt.Errorf("The resource not support import ")
			},
		},
		Schema: m,
	}
}

func resourceByteplusDefaultNodePoolBatchAttachCreate(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewByteplusVkeDefaultNodePoolBatchAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(nodePoolService, d, ResourceByteplusDefaultNodePoolBatchAttach())
	if err != nil {
		return fmt.Errorf("error on creating DefaultNodePoolBatchAttach  %q, %w", d.Id(), err)
	}
	return resourceByteplusDefaultNodePoolBatchAttachRead(d, meta)
}

func resourceByteplusDefaultNodePoolBatchAttachRead(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewByteplusVkeDefaultNodePoolBatchAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(nodePoolService, d, ResourceByteplusDefaultNodePoolBatchAttach())
	if err != nil {
		return fmt.Errorf("error on reading DefaultNodePoolBatchAttach %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusDefaultNodePoolBatchAttachUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewByteplusVkeDefaultNodePoolBatchAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(nodePoolService, d, ResourceByteplusDefaultNodePoolBatchAttach())
	if err != nil {
		return fmt.Errorf("error on updating DefaultNodePoolBatchAttach  %q, %w", d.Id(), err)
	}
	return resourceByteplusDefaultNodePoolBatchAttachRead(d, meta)
}

func resourceByteplusNodePoolBatchAttachDelete(d *schema.ResourceData, meta interface{}) (err error) {
	nodePoolService := NewByteplusVkeDefaultNodePoolBatchAttachService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(nodePoolService, d, ResourceByteplusDefaultNodePoolBatchAttach())
	if err != nil {
		return fmt.Errorf("error on deleting DefaultNodePoolBatchAttach %q, %w", d.Id(), err)
	}
	return err
}
