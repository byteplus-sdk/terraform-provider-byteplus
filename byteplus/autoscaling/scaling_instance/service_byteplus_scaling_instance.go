package scaling_instance

import (
	"errors"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

type ByteplusScalingInstanceService struct {
	Client *bp.SdkClient
}

func (s *ByteplusScalingInstanceService) ReadResource(data *schema.ResourceData, s2 string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (s *ByteplusScalingInstanceService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, s2 string) *resource.StateChangeConf {
	return nil
}

func NewScalingInstanceService(c *bp.SdkClient) *ByteplusScalingInstanceService {
	return &ByteplusScalingInstanceService{
		Client: c,
	}
}

func (s *ByteplusScalingInstanceService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusScalingInstanceService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 50, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeScalingInstances"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = universalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = universalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)
		results, err = bp.ObtainSdkValue("Result.ScalingInstances", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ScalingInstances is not Slice")
		}
		return data, err
	})
}

func (ByteplusScalingInstanceService) WithResourceResponseHandlers(scalingInstance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return scalingInstance, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusScalingInstanceService) CreateResource(d *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusScalingInstanceService) ModifyResource(data *schema.ResourceData, s2 *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusScalingInstanceService) RemoveResource(data *schema.ResourceData, s2 *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusScalingInstanceService) DatasourceResources(data *schema.ResourceData, s2 *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "InstanceIds",
				ConvertType: bp.ConvertWithN,
			},
			"status": {
				TargetField: "Status",
				ConvertType: bp.ConvertDefault,
			},
		},
		ResponseConverts: map[string]bp.ResponseConvert{
			"InstanceId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
		IdField:      "InstanceId",
		CollectField: "scaling_instances",
	}
}

func (s *ByteplusScalingInstanceService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "auto_scaling",
		Action:      actionName,
		Version:     "2020-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}
