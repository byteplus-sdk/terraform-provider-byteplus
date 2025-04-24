package allowlist_associate

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
RDS AllowList Associate can be imported using the instance id and allow list id, e.g.
```
$ terraform import byteplus_rds_mysql_allowlist_associate.default rds-mysql-h441603c68aaa:acl-d1fd76693bd54e658912e7337d5b****
```

*/

func ResourceByteplusRdsMysqlAllowlistAssociate() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusRdsMysqlAllowlistAssociateCreate,
		Read:   resourceByteplusRdsMysqlAllowlistAssociateRead,
		Delete: resourceByteplusRdsMysqlAllowlistAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: importAllowListAssociate,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the mysql instance.",
			},
			"allow_list_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the allow list.",
			},
		},
	}
}

func resourceByteplusRdsMysqlAllowlistAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlAllowListAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusRdsMysqlAllowlistAssociate())
	if err != nil {
		return fmt.Errorf("error creating RDS Mysql Allowlist Associate service: %q, %w", d.Id(), err)
	}
	return resourceByteplusRdsMysqlAllowlistAssociateRead(d, meta)
}

func resourceByteplusRdsMysqlAllowlistAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlAllowListAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusRdsMysqlAllowlistAssociate())
	if err != nil {
		return fmt.Errorf("error reading RDS Mysql Allowlist Associate service: %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusRdsMysqlAllowlistAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlAllowListAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusRdsMysqlAllowlistAssociate())
	if err != nil {
		return fmt.Errorf("error deleting RDS Mysql Allowlist Associate service: %q, %w", d.Id(), err)
	}
	return err
}

func importAllowListAssociate(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must be of the form InstanceId:AllowListId")
	}
	err = data.Set("instance_id", items[0])
	if err != nil {
		return []*schema.ResourceData{data}, err
	}
	err = data.Set("allow_list_id", items[1])
	if err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
