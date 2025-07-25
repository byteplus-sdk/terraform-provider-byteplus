package kafka_allow_list_associate

import (
	"errors"
	"fmt"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/kafka/kafka_instance"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusKafkaAllowListAssociateService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewKafkaAllowListAssociateService(c *bp.SdkClient) *ByteplusKafkaAllowListAssociateService {
	return &ByteplusKafkaAllowListAssociateService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusKafkaAllowListAssociateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusKafkaAllowListAssociateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		return data, err
	})
}

func (s *ByteplusKafkaAllowListAssociateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results     interface{}
		resultsMap  map[string]interface{}
		instanceMap map[string]interface{}
		instances   []interface{}
		ok          bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf("invalid kafka_allow_list_associate id: %v", id)
	}
	req := map[string]interface{}{
		"AllowListId": ids[1],
	}
	action := "DescribeAllowListDetail"
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
	if err != nil {
		return data, err
	}
	results, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if resultsMap, ok = results.(map[string]interface{}); !ok {
		return resultsMap, errors.New("Value is not map ")
	}
	if len(resultsMap) == 0 {
		return resultsMap, fmt.Errorf("Kafka allowlist %s not exist ", ids[1])
	}
	logger.Debug(logger.ReqFormat, action, resultsMap)
	instances, ok = resultsMap["AssociatedInstances"].([]interface{})
	if !ok {
		return data, errors.New("Value is not slice ")
	}
	logger.Debug(logger.ReqFormat, action, instances)
	for _, instance := range instances {
		if instanceMap, ok = instance.(map[string]interface{}); !ok {
			return data, errors.New("instance is not map ")
		}
		if len(instanceMap) == 0 {
			continue
		}
		if instanceMap["InstanceId"].(string) == ids[0] {
			data = resultsMap
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Kafka allowlist associate %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusKafkaAllowListAssociateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
	}
}

func (s *ByteplusKafkaAllowListAssociateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	instanceId := resourceData.Get("instance_id").(string)
	allowListId := resourceData.Get("allow_list_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AssociateAllowList",
			ConvertMode: bp.RequestConvertIgnore,
			Convert:     map[string]bp.RequestConvert{},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["InstanceIds"] = []string{instanceId}
				(*call.SdkParam)["AllowListIds"] = []string{allowListId}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprintf("%s:%s", instanceId, allowListId))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return instanceId
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				kafka_instance.NewKafkaInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: instanceId,
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusKafkaAllowListAssociateService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusKafkaAllowListAssociateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{}
	return []bp.Callback{callback}
}

func (s *ByteplusKafkaAllowListAssociateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	instanceId := resourceData.Get("instance_id").(string)
	allowListId := resourceData.Get("allow_list_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DisassociateAllowList",
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceIds":  []string{instanceId},
				"AllowListIds": []string{allowListId},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				err := bp.CheckResourceUtilRemoved(d, s.ReadResource, 10*time.Minute)
				return err
			},
			LockId: func(d *schema.ResourceData) string {
				return instanceId
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				kafka_instance.NewKafkaInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: instanceId,
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusKafkaAllowListAssociateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusKafkaAllowListAssociateService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Kafka",
		Version:     "2022-05-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
