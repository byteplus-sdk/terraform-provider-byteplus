package scaling_activity

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusScalingActivityService struct {
	Client *bp.SdkClient
}

func NewScalingActivityService(c *bp.SdkClient) *ByteplusScalingActivityService {
	return &ByteplusScalingActivityService{
		Client: c,
	}
}

func (s *ByteplusScalingActivityService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusScalingActivityService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeScalingActivities"
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
		logger.Debug(logger.RespFormat, action, resp)
		results, err = bp.ObtainSdkValue("Result.ScalingActivities", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ScalingActivities is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusScalingActivityService) ReadResource(resourceData *schema.ResourceData, activityId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if activityId == "" {
		activityId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ScalingActivityIds.1": activityId,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Scaling Activity %s not exist ", activityId)
	}
	return data, err
}

func (s *ByteplusScalingActivityService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusScalingActivityService) WithResourceResponseHandlers(data map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return data, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusScalingActivityService) CreateResource(*schema.ResourceData, *schema.Resource) []bp.Callback {
	return nil
}

func (s *ByteplusScalingActivityService) ModifyResource(*schema.ResourceData, *schema.Resource) []bp.Callback {
	return nil
}

func (s *ByteplusScalingActivityService) RemoveResource(*schema.ResourceData, *schema.Resource) []bp.Callback {
	return nil
}

func (s *ByteplusScalingActivityService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "ScalingActivityIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		IdField:      "ScalingActivityId",
		CollectField: "activities",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ScalingActivityId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (s *ByteplusScalingActivityService) ReadResourceId(id string) string {
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
