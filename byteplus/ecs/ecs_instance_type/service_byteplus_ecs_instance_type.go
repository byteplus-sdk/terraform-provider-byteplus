package ecs_instance_type

import (
	"encoding/json"
	"errors"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusEcsInstanceTypeService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewEcsInstanceTypeService(c *bp.SdkClient) *ByteplusEcsInstanceTypeService {
	return &ByteplusEcsInstanceTypeService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusEcsInstanceTypeService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusEcsInstanceTypeService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp      *map[string]interface{}
		results   interface{}
		nextToken interface{}
		next      string
		ok        bool
	)
	return bp.WithNextTokenQuery(condition, "MaxResults", "NextToken", 20, nil, func(m map[string]interface{}) ([]interface{}, string, error) {
		action := "DescribeInstanceTypes"
		bytes, _ := json.Marshal(condition)
		logger.Debug(logger.ReqFormat, action, string(bytes))
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, next, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, next, err
			}
		}
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))
		results, err = bp.ObtainSdkValue("Result.InstanceTypes", *resp)
		if err != nil {
			return data, next, err
		}
		nextToken, err = bp.ObtainSdkValue("Result.NextToken", *resp)
		if err != nil {
			return data, next, err
		}
		next = nextToken.(string)
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, next, errors.New("Result.InstanceTypes is not Slice")
		}
		return data, next, err
	})
}

func (s *ByteplusEcsInstanceTypeService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return data, err
}

func (s *ByteplusEcsInstanceTypeService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusEcsInstanceTypeService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (ByteplusEcsInstanceTypeService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusEcsInstanceTypeService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusEcsInstanceTypeService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusEcsInstanceTypeService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "InstanceTypeIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		IdField:      "InstanceTypeId",
		CollectField: "instance_types",
	}
}

func (s *ByteplusEcsInstanceTypeService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ecs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
