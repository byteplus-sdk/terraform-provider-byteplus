package kafka_consumed_topic

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusKafkaConsumedTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusKafkaConsumedTopicsRead,
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
				Optional:    true,
				Description: "The name of kafka topic. This field supports fuzzy query.",
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
			"consumed_topics": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of kafka topic.",
						},
						"accumulation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total amount of message accumulation in this topic for the consumer group.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusKafkaConsumedTopicsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewKafkaConsumedTopicService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(service, d, DataSourceByteplusKafkaConsumedTopics())
}
