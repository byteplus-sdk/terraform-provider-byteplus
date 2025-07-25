package kafka_public_address

import (
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/kafka/kafka_instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusKafkaPublicAddressService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewKafkaInternetEnablerService(c *bp.SdkClient) *ByteplusKafkaPublicAddressService {
	return &ByteplusKafkaPublicAddressService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusKafkaPublicAddressService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusKafkaPublicAddressService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *ByteplusKafkaPublicAddressService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		resp *map[string]interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	req := map[string]interface{}{
		"InstanceID": ids[0],
	}
	logger.Debug(logger.ReqFormat, "DescribeInstanceDetail", req)
	resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo("DescribeInstanceDetail"), &req)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, "DescribeInstanceDetail", req, *resp)
	eip, err := bp.ObtainSdkValue("Result.BasicInstanceInfo.EipId", *resp)
	if err != nil {
		return data, err
	}
	if eip == nil || eip != resourceData.Get("eip_id") {
		return nil, fmt.Errorf("instance_id and eip not associate")
	}

	connection, err := bp.ObtainSdkValue("Result.ConnectionInfo", *resp)
	if err != nil {
		return data, err
	}
	for _, ele := range connection.([]interface{}) {
		conn := ele.(map[string]interface{})
		if conn["EndpointType"] == "SASL_SSL" {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("instance_id and eip not associate")
}

func (s *ByteplusKafkaPublicAddressService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusKafkaPublicAddressService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreatePublicAddress",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["InstanceId"] = resourceData.Get("instance_id")
				(*call.SdkParam)["EipId"] = resourceData.Get("eip_id")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprintf("%s:%s", resourceData.Get("instance_id"), resourceData.Get("eip_id")))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				kafka_instance.NewKafkaInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusKafkaPublicAddressService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusKafkaPublicAddressService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusKafkaPublicAddressService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeletePublicAddress",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": resourceData.Get("instance_id"),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				kafka_instance.NewKafkaInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusKafkaPublicAddressService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusKafkaPublicAddressService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "kafka",
		Version:     "2022-05-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
