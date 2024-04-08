package region

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRegionService struct {
	Client *bp.SdkClient
}

func (v *ByteplusRegionService) GetClient() *bp.SdkClient {
	return v.Client
}

func (v *ByteplusRegionService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp      *map[string]interface{}
		nextToken interface{}
		results   interface{}
		next      string
		ok        bool
	)
	return bp.WithNextTokenQuery(condition, "MaxResults", "NextToken", 10, nil, func(m map[string]interface{}) ([]interface{}, string, error) {
		action := "DescribeRegions"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
		} else {
			resp, err = v.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
		}
		if err != nil {
			return nil, next, err
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.Regions", *resp)
		if err != nil {
			return nil, next, err
		}
		nextToken, err = bp.ObtainSdkValue("Result.NextToken", *resp)
		if err != nil {
			return nil, next, err
		}
		next, ok = nextToken.(string)
		if !ok {
			return nil, next, fmt.Errorf("next token must be a string")
		}
		if results == nil {
			results = make([]interface{}, 0)
		}

		if data, ok = results.([]interface{}); !ok {
			return nil, next, errors.New("Result.Regions is not Slice")
		}

		return data, next, err
	})
}

func (v *ByteplusRegionService) ReadResource(data *schema.ResourceData, s string) (map[string]interface{}, error) {
	return nil, nil
}

func (v *ByteplusRegionService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			return nil, "", err
		},
	}
}

func (v *ByteplusRegionService) WithResourceResponseHandlers(region map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return region, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (v *ByteplusRegionService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (v *ByteplusRegionService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (v *ByteplusRegionService) RemoveResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (v *ByteplusRegionService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "RegionIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "RegionId",
		IdField:      "RegionId",
		CollectField: "regions",
		ResponseConverts: map[string]bp.ResponseConvert{
			"RegionId": {
				TargetField: "id",
				KeepDefault: true,
			},
		},
	}
}

func (v *ByteplusRegionService) ReadResourceId(s string) string {
	return s
}

func NewRegionService(c *bp.SdkClient) *ByteplusRegionService {
	return &ByteplusRegionService{
		Client: c,
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ecs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
