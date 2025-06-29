package rds_mysql_endpoint

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
RdsMysqlEndpoint can be imported using the instance id and endpoint id, e.g.
```
$ terraform import byteplus_rds_mysql_endpoint.default mysql-3c25f219***:mysql-3c25f219****-custom-eeb5
```

*/

func ResourceByteplusRdsMysqlEndpoint() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusRdsMysqlEndpointCreate,
		Read:   resourceByteplusRdsMysqlEndpointRead,
		Update: resourceByteplusRdsMysqlEndpointUpdate,
		Delete: resourceByteplusRdsMysqlEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: endpointImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the mysql instance.",
			},
			"endpoint_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The id of the endpoint. Import an exist endpoint, usually for import a default endpoint generated with instance creating.",
			},
			"read_write_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ReadOnly",
				Description: "Reading and writing mode: ReadWrite, ReadOnly(Default).",
			},
			"endpoint_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the endpoint.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the endpoint.",
			},
			"nodes": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:      schema.HashString,
				Required: true,
				Description: "List of node IDs configured for the endpoint. Required when EndpointType is Custom. " +
					"To add a master node to the terminal, there is no need to fill in the master node ID, just fill in `Primary`.",
			},
			"auto_add_new_nodes": {
				Type:     schema.TypeBool,
				Computed: false,
				Optional: true,
				Description: "When the terminal type is a read-write terminal or a read-only terminal, " +
					"support is provided for setting whether new nodes are automatically added." +
					" The values are:\ntrue: Automatically add.\nfalse: Do not automatically add (default).",
			},
			"read_write_spliting": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					//当 ReadWriteMode 为读写时支持设置；当 ReadWriteMode 为只读时不支持设置。此参数仅对默认终端生效。
					return d.Get("read_write_mode").(string) == "ReadOnly" ||
						d.Get("endpoint_id").(string) == ""
				},
				Description: "Enable read-write separation. Possible values: TRUE, FALSE.\n" +
					"This setting can be configured when ReadWriteMode is set to read-write, " +
					"but cannot be configured when ReadWriteMode is set to read-only. " +
					"This parameter only applies to the default terminal.",
			},
			"read_only_node_max_delay_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					//读写类型的终端，且开通读写分离后支持设置此参数。
					return d.Get("read_write_mode").(string) == "ReadOnly" ||
						!d.Get("read_write_spliting").(bool)
				},
				Description: "The maximum delay threshold for read-only nodes, when the delay time of a read-only node exceeds this value, " +
					"the read traffic will not be sent to that node, unit: seconds. " +
					"Value range: 0~3600. Default value: 30.",
			},
			"read_only_node_distribution_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.Get("read_write_spliting").(bool)
				},
				Description: "Read weight allocation mode. This parameter is required when enabling read-write separation setting to TRUE. " +
					"Possible values:\nDefault: Automatically allocate weights based on specifications (default).\nCustom: Custom weight allocation.",
			},
			"read_only_node_weight": {
				Type: schema.TypeSet,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					//当 ReadOnlyNodeDistributionType 取值为 Custom 时，需要传入此参数。
					return d.Get("read_only_node_distribution_type").(string) != "Custom"
				},
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Read-only nodes require NodeId to be passed, while primary nodes do not require it.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The primary node needs to pass in the NodeType as Primary, while the read-only node does not need to pass it in.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The read weight of the node increases by 100, with a maximum value of 10000.",
						},
					},
				},
				Optional: true,
				Description: "Customize read weight distribution, that is, pass in the read request weight of the master node and read-only nodes. " +
					"It increases by 100 and the maximum value is 10000. " +
					"When the ReadOnlyNodeDistributionType value is Custom, " +
					"this parameter needs to be passed in.",
			},
			//"dns_visibility": {
			//	Type:        schema.TypeBool,
			//	Optional:    true,
			//	Computed:    true,
			//	Description: "Values:\nfalse: Volcano Engine private network resolution (default).\ntrue: Volcano Engine private and public network resolution.",
			//},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Connection address, Please note that the connection address can only modify the prefix." +
					" In one call, it is not possible to modify both the connection address prefix and the port at the same time.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The port. Cannot modify public network port. In one call, it is not possible to modify both the connection address prefix and the port at the same time.",
			},
		},
	}
	return resource
}

func resourceByteplusRdsMysqlEndpointCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlEndpointService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusRdsMysqlEndpoint())
	if err != nil {
		return fmt.Errorf("error on creating rds_mysql_endpoint %q, %s", d.Id(), err)
	}
	return resourceByteplusRdsMysqlEndpointRead(d, meta)
}

func resourceByteplusRdsMysqlEndpointRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlEndpointService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusRdsMysqlEndpoint())
	if err != nil {
		return fmt.Errorf("error on reading rds_mysql_endpoint %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusRdsMysqlEndpointUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlEndpointService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusRdsMysqlEndpoint())
	if err != nil {
		return fmt.Errorf("error on updating rds_mysql_endpoint %q, %s", d.Id(), err)
	}
	return resourceByteplusRdsMysqlEndpointRead(d, meta)
}

func resourceByteplusRdsMysqlEndpointDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlEndpointService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusRdsMysqlEndpoint())
	if err != nil {
		return fmt.Errorf("error on deleting rds_mysql_endpoint %q, %s", d.Id(), err)
	}
	return err
}

func endpointImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must be 'instanceId:endpointId'")
	}
	instanceId := items[0]
	endpointId := items[1]
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("endpoint_id", endpointId)

	return []*schema.ResourceData{d}, nil
}
