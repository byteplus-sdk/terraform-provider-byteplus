package pitr_time_period

import (
	"errors"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRedisPitrTimePeriodService struct {
	Client *bp.SdkClient
}

func (v *ByteplusRedisPitrTimePeriodService) GetClient() *bp.SdkClient {
	return v.Client
}

func (v *ByteplusRedisPitrTimePeriodService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		ids       []interface{}
		resp      *map[string]interface{}
		result    interface{}
		results   []interface{}
		resultMap map[string]interface{}
	)
	action := "DescribePitrTimeWindow"
	ids = m["Ids"].(*schema.Set).List()
	if len(ids) == 0 {
		return data, nil
	}
	for _, id := range ids {
		instanceId, ok := id.(string)
		if !ok {
			return data, errors.New("err instance id")
		}
		req := map[string]interface{}{
			"InstanceId": instanceId,
		}
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		result, err = bp.ObtainSdkValue("Result", *resp)
		if err != nil {
			return data, err
		}
		if resultMap, ok = result.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		}
		// 加个ID，方便对照
		resultMap["InstanceId"] = instanceId
		results = append(results, resultMap)
	}
	return results, nil
}

func (v *ByteplusRedisPitrTimePeriodService) ReadResource(data *schema.ResourceData, s string) (map[string]interface{}, error) {
	return nil, nil
}

func (v *ByteplusRedisPitrTimePeriodService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, s string) *resource.StateChangeConf {
	return nil
}

func (v *ByteplusRedisPitrTimePeriodService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return m, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (v *ByteplusRedisPitrTimePeriodService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (v *ByteplusRedisPitrTimePeriodService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (v *ByteplusRedisPitrTimePeriodService) RemoveResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (v *ByteplusRedisPitrTimePeriodService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		CollectField: "periods",
	}
}

func (v *ByteplusRedisPitrTimePeriodService) ReadResourceId(s string) string {
	return s
}

func NewByteplusRedisPitrTimeWindowService(c *bp.SdkClient) *ByteplusRedisPitrTimePeriodService {
	return &ByteplusRedisPitrTimePeriodService{
		Client: c,
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Redis",
		Action:      actionName,
		Version:     "2020-12-07",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
	}
}
