package instance_parameter_log

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusMongoDBInstanceParameterLogService struct {
	Client *bp.SdkClient
}

func NewMongoDBInstanceParameterLogService(c *bp.SdkClient) *ByteplusMongoDBInstanceParameterLogService {
	return &ByteplusMongoDBInstanceParameterLogService{
		Client: c,
	}
}

func (s *ByteplusMongoDBInstanceParameterLogService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusMongoDBInstanceParameterLogService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 100, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeDBInstanceParametersLog"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)
		results, err = bp.ObtainSdkValue("Result.ParameterChangeLog", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		data, ok = results.([]interface{})
		if !ok {
			return data, fmt.Errorf("DescribeDBInstanceParametersLog response is not a slice")
		}
		return data, nil
	})
}

func (s *ByteplusMongoDBInstanceParameterLogService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return data, err
}

func (s *ByteplusMongoDBInstanceParameterLogService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusMongoDBInstanceParameterLogService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusMongoDBInstanceParameterLogService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusMongoDBInstanceParameterLogService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusMongoDBInstanceParameterLogService) RemoveResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusMongoDBInstanceParameterLogService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "parameter_change_logs",
		ContentType:  bp.ContentTypeJson,
	}
}

func (s *ByteplusMongoDBInstanceParameterLogService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "mongodb",
		Action:      actionName,
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
	}
}
