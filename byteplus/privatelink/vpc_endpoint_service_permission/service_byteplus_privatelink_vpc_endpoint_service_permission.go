package vpc_endpoint_service_permission

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/privatelink/vpc_endpoint_service"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusService struct {
	Client *bp.SdkClient
}

func NewService(c *bp.SdkClient) *ByteplusService {
	return &ByteplusService{
		Client: c,
	}
}

func (s *ByteplusService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeVpcEndpointServicePermissions"
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
		logger.Debug(logger.RespFormat, action, *resp)
		results, err = bp.ObtainSdkValue("Result.Permissions", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Permissions is not Slice")
		}

		return data, err
	})

}

func (s *ByteplusService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	ids := strings.Split(id, ":")
	serviceId := ids[0]
	accountId := ids[1]

	results, err = s.ReadResources(map[string]interface{}{
		"ServiceId":       serviceId,
		"PermitAccountId": accountId,
	})
	if err != nil {
		return data, err
	}
	if len(results) == 0 {
		return data, fmt.Errorf("Vpc endpoint service permission %s not exist ", id)
	}

	return map[string]interface{}{
		"ServiceId":       serviceId,
		"PermitAccountId": accountId,
	}, nil
}

func (s *ByteplusService) WithResourceResponseHandlers(nodePool map[string]interface{}) []bp.ResourceResponseHandler {
	return nil
}

func (s *ByteplusService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddPermissionToVpcEndpointService",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprint((*call.SdkParam)["ServiceId"], ":", (*call.SdkParam)["PermitAccountId"]))
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc_endpoint_service.NewService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("service_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("service_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	return callbacks
}

func (s *ByteplusService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(s.ReadResourceId(resourceData.Id()), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "RemovePermissionFromVpcEndpointService",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"ServiceId":       ids[0],
				"PermitAccountId": ids[1],
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				return resp, err
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("service_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc_endpoint_service.NewService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("service_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "permissions",
	}
}

func (s *ByteplusService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "privatelink",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
		ContentType: bp.Default,
	}
}
