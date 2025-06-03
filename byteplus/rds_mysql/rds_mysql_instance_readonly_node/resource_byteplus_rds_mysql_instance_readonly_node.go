package rds_mysql_instance_readonly_node

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Rds Mysql Instance Readonly Node can be imported using the instance_id:node_id, e.g.
```
$ terraform import byteplus_rds_mysql_instance_readonly_node.default mysql-72da4258c2c7:mysql-72da4258c2c7-r7f93
```

*/

func ResourceByteplusRdsMysqlInstanceReadonlyNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusRdsMysqlInstanceReadonlyNodeCreate,
		Read:   resourceByteplusRdsMysqlInstanceReadonlyNodeRead,
		Update: resourceByteplusRdsMysqlInstanceReadonlyNodeUpdate,
		Delete: resourceByteplusRdsMysqlInstanceReadonlyNodeDelete,
		Importer: &schema.ResourceImporter{
			State: rdsMysqlInstanceReadonlyNodeImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The RDS mysql instance id of the readonly node.",
			},
			"node_spec": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The specification of readonly node.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The available zone of readonly node.",
			},
			"node_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the readonly node.",
			},
		},
	}
}

func resourceByteplusRdsMysqlInstanceReadonlyNodeCreate(d *schema.ResourceData, meta interface{}) (err error) {
	rdsMysqlInstanceReadonlyNodeService := NewRdsMysqlInstanceReadonlyNodeService(meta.(*bp.SdkClient))
	err = rdsMysqlInstanceReadonlyNodeService.Dispatcher.Create(rdsMysqlInstanceReadonlyNodeService, d, ResourceByteplusRdsMysqlInstanceReadonlyNode())
	if err != nil {
		return fmt.Errorf("error on creating RDS mysql instance readonly node %q, %w", d.Id(), err)
	}
	return resourceByteplusRdsMysqlInstanceReadonlyNodeRead(d, meta)
}

func resourceByteplusRdsMysqlInstanceReadonlyNodeRead(d *schema.ResourceData, meta interface{}) (err error) {
	rdsMysqlInstanceReadonlyNodeService := NewRdsMysqlInstanceReadonlyNodeService(meta.(*bp.SdkClient))
	err = rdsMysqlInstanceReadonlyNodeService.Dispatcher.Read(rdsMysqlInstanceReadonlyNodeService, d, ResourceByteplusRdsMysqlInstanceReadonlyNode())
	if err != nil {
		return fmt.Errorf("error on reading RDS mysql instance readonly node %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusRdsMysqlInstanceReadonlyNodeUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	rdsMysqlInstanceReadonlyNodeService := NewRdsMysqlInstanceReadonlyNodeService(meta.(*bp.SdkClient))
	err = rdsMysqlInstanceReadonlyNodeService.Dispatcher.Update(rdsMysqlInstanceReadonlyNodeService, d, ResourceByteplusRdsMysqlInstanceReadonlyNode())
	if err != nil {
		return fmt.Errorf("error on updating RDS mysql instance readonly node %q, %w", d.Id(), err)
	}
	return resourceByteplusRdsMysqlInstanceReadonlyNodeRead(d, meta)
}

func resourceByteplusRdsMysqlInstanceReadonlyNodeDelete(d *schema.ResourceData, meta interface{}) (err error) {
	rdsMysqlInstanceReadonlyNodeService := NewRdsMysqlInstanceReadonlyNodeService(meta.(*bp.SdkClient))
	err = rdsMysqlInstanceReadonlyNodeService.Dispatcher.Delete(rdsMysqlInstanceReadonlyNodeService, d, ResourceByteplusRdsMysqlInstanceReadonlyNode())
	if err != nil {
		return fmt.Errorf("error on deleting RDS mysql instance readonly node %q, %w", d.Id(), err)
	}
	return err
}
