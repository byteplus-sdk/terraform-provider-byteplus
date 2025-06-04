package backup

import (
	"fmt"
	"strings"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Redis Backup can be imported using the instanceId:backupId, e.g.
```
$ terraform import byteplus_redis_backup.default redis-cn02aqusft7ws****:b-cn02xmmrp751i9cdzcphjmk4****
```

*/

func ResourceByteplusRedisBackup() *schema.Resource {
	resource := &schema.Resource{
		Read:   resourceByteplusRedisBackupRead,
		Create: resourceByteplusRedisBackupCreate,
		Delete: resourceByteplusRedisBackupDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
				}
				if err := data.Set("instance_id", items[0]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				if err := data.Set("backup_point_id", items[1]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				return []*schema.ResourceData{data}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of instance to create backup.",
			},
			"backup_point_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of backup point.",
			},
		},
	}
	bp.MergeDateSourceToResource(DataSourceByteplusRedisBackups().Schema["backups"].Elem.(*schema.Resource).Schema, &resource.Schema)
	return resource
}

func resourceByteplusRedisBackupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	redisBackupService := NewRedisBackupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(redisBackupService, d, ResourceByteplusRedisBackup())
	if err != nil {
		return fmt.Errorf("error on creating redis backup %v, %v", d.Id(), err)
	}
	return resourceByteplusRedisBackupRead(d, meta)
}

func resourceByteplusRedisBackupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	return nil
}

func resourceByteplusRedisBackupRead(d *schema.ResourceData, meta interface{}) (err error) {
	redisBackupService := NewRedisBackupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(redisBackupService, d, ResourceByteplusRedisBackup())
	if err != nil {
		return fmt.Errorf("error on reading redis backup %q,%s", d.Id(), err)
	}
	return err
}
