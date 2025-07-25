package kafka_consumed_partition

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusKafkaConsumedPartitions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusKafkaConsumedPartitionsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of kafka instance.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of kafka group.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of kafka topic.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of query.",
			},
			"consumed_partitions": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"partition_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The index number of partition.",
						},
						"accumulation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total amount of message accumulation in this topic partition for the consumer group.",
						},
						"consumed_client": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The consumed client info of partition.",
						},
						"consumed_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The consumed offset of partition.",
						},
						"start_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The start offset of partition.",
						},
						"end_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The end offset of partition.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusKafkaConsumedPartitionsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewKafkaConsumedPartitionService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(service, d, DataSourceByteplusKafkaConsumedPartitions())
}
