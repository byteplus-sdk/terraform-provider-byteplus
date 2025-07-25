package continuous_backup

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Redis Continuous Backup can be imported using the continuous:instanceId, e.g.
```
$ terraform import byteplus_redis_continuous_backup.default continuous:redis-asdljioeixxxx
```
*/

func ResourceByteplusRedisContinuousBackup() *schema.Resource {
	resource := &schema.Resource{
		Read:   resourceByteplusRedisContinuousBackupRead,
		Create: resourceByteplusRedisContinuousBackupCreate,
		Delete: resourceByteplusRedisContinuousBackupDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
				}
				if err := data.Set("instance_id", items[1]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				return []*schema.ResourceData{data}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Id of redis instance.",
			},
		},
	}
	return resource
}

func resourceByteplusRedisContinuousBackupRead(d *schema.ResourceData, meta interface{}) (err error) {
	redisContinuousBackupService := NewRedisContinuousBackupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(redisContinuousBackupService, d, ResourceByteplusRedisContinuousBackup())
	if err != nil {
		return fmt.Errorf("error on read continuous backup %v, %v", d.Id(), err)
	}
	return nil
}

func resourceByteplusRedisContinuousBackupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	redisContinuousBackupService := NewRedisContinuousBackupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(redisContinuousBackupService, d, ResourceByteplusRedisContinuousBackup())
	if err != nil {
		return fmt.Errorf("error on create continuous backup %v, %v", d.Id(), err)
	}
	return resourceByteplusRedisContinuousBackupRead(d, meta)
}

func resourceByteplusRedisContinuousBackupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	redisContinuousBackupService := NewRedisContinuousBackupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(redisContinuousBackupService, d, ResourceByteplusRedisContinuousBackup())
	if err != nil {
		return fmt.Errorf("error on delete continuous backup %v, %v", d.Id(), err)
	}
	return nil
}
