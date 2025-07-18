package cr_registry

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CR Registry can be imported using the name, e.g.
```
$ terraform import byteplus_cr_registry.default enterprise-x
```

*/

func ResourceByteplusCrRegistry() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCrRegistryCreate,
		Read:   resourceByteplusCrRegistryRead,
		Update: resourceByteplusCrRegistryUpdate,
		Delete: resourceByteplusCrRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of registry.",
			},
			"delete_immediately": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether delete registry immediately. Only effected in delete action.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password of registry user.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ProjectName of the cr registry.",
			},
			"resource_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "Tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The Key of Tags.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The Value of Tags.",
						},
					},
				},
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The type of registry. Valid values: `Enterprise`, `Micro`. Default is `Enterprise`.",
			},
			"proxy_cache_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Whether to enable proxy cache.",
			},
			"proxy_cache": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "The proxy cache of registry. This field is valid when proxy_cache_enabled is true.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The type of proxy cache. Valid values: `DockerHub`, `DockerRegistry`.",
						},
						"endpoint": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "The endpoint of proxy cache.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Sensitive:   true,
							Description: "The password of proxy cache.",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "The username of proxy cache.",
						},
						"skip_ssl_verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Whether to skip ssl verify.",
						},
					},
				},
			},
		},
	}
	dataSource := DataSourceByteplusCrRegistries().Schema["registries"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusCrRegistryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRegistryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusCrRegistry())
	if err != nil {
		return fmt.Errorf("error on creating CrRegistry %q,%s", d.Id(), err)
	}
	return resourceByteplusCrRegistryRead(d, meta)
}

func resourceByteplusCrRegistryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRegistryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusCrRegistry())
	if err != nil {
		return fmt.Errorf("error on updating CrRegistry  %q, %s", d.Id(), err)
	}
	return resourceByteplusCrRegistryRead(d, meta)
}

func resourceByteplusCrRegistryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRegistryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCrRegistry())
	if err != nil {
		return fmt.Errorf("error on deleting CrRegistry %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCrRegistryRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrRegistryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCrRegistry())
	if err != nil {
		return fmt.Errorf("Error on reading CrRegistry %q,%s", d.Id(), err)
	}
	return err
}
