package support_resource_types

import (
	"errors"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusVkeSupportResourceTypeService struct {
	Client *bp.SdkClient
}

func NewService(c *bp.SdkClient) *ByteplusVkeSupportResourceTypeService {
	return &ByteplusVkeSupportResourceTypeService{
		Client: c,
	}
}

func (s *ByteplusVkeSupportResourceTypeService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusVkeSupportResourceTypeService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "ListSupportedResourceTypes"
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
		results, err = bp.ObtainSdkValue("Result.Items", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Items is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusVkeSupportResourceTypeService) ReadResource(resourceData *schema.ResourceData, clusterId string) (data map[string]interface{}, err error) {
	return data, err
}

func (s *ByteplusVkeSupportResourceTypeService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (ByteplusVkeSupportResourceTypeService) WithResourceResponseHandlers(cluster map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *ByteplusVkeSupportResourceTypeService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}

}

func (s *ByteplusVkeSupportResourceTypeService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusVkeSupportResourceTypeService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusVkeSupportResourceTypeService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"zone_ids": {
				TargetField: "Filter.ZoneIds",
				ConvertType: bp.ConvertJsonArray,
			},
			"resource_types": {
				TargetField: "Filter.ResourceTypes",
				ConvertType: bp.ConvertJsonArray,
			},
		},
		ContentType:  bp.ContentTypeJson,
		CollectField: "resources",
	}
}

func (s *ByteplusVkeSupportResourceTypeService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vke",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
