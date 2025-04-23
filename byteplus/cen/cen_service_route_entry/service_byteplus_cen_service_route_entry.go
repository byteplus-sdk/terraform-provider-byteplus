package cen_service_route_entry

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCenServiceRouteEntryService struct {
	Client *bp.SdkClient
}

func NewCenServiceRouteEntryService(c *bp.SdkClient) *ByteplusCenServiceRouteEntryService {
	return &ByteplusCenServiceRouteEntryService{
		Client: c,
	}
}

func (s *ByteplusCenServiceRouteEntryService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCenServiceRouteEntryService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeCenServiceRouteEntries"
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
		results, err = bp.ObtainSdkValue("Result.ServiceRouteEntries", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ServiceRouteEntries is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusCenServiceRouteEntryService) ReadResource(resourceData *schema.ResourceData, tmpId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if tmpId == "" {
		tmpId = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(tmpId, ":")
	req := map[string]interface{}{
		"CenId":                ids[0],
		"DestinationCidrBlock": ids[1],
		"ServiceRegionId":      ids[2],
		"ServiceVpcId":         ids[3],
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
		return data, fmt.Errorf("cen service route entry %s not exist ", tmpId)
	}
	return data, err
}

func (s *ByteplusCenServiceRouteEntryService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo       map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Error")
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("cen service route entry status error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (ByteplusCenServiceRouteEntryService) WithResourceResponseHandlers(v map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return v, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusCenServiceRouteEntryService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateCenServiceRouteEntry",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"publish_to_instances": {
					ConvertType: bp.ConvertListN,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				d.SetId(fmt.Sprintf("%v:%v:%v:%v", resourceData.Get("cen_id"), resourceData.Get("destination_cidr_block"),
					resourceData.Get("service_region_id"), resourceData.Get("service_vpc_id")))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("cen_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				cen.NewCenService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("cen_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *ByteplusCenServiceRouteEntryService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyCenServiceRouteEntryAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"cen_id": {
					ForceGet: true,
				},
				"destination_cidr_block": {
					ForceGet: true,
				},
				"service_region_id": {
					ForceGet: true,
				},
				"description": {
					ConvertType: bp.ConvertDefault,
				},
				"publish_mode": {
					ForceGet:    true,
					ConvertType: bp.ConvertDefault,
				},
				"publish_to_instances": {
					ForceGet:    true,
					ConvertType: bp.ConvertListN,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				mode := d.Get("publish_mode").(string)
				instances := d.Get("publish_to_instances").(*schema.Set).List()
				if mode == "Custom" && len(instances) == 0 {
					return false, fmt.Errorf("public_to_instances must exist when publish_mode is Custom")
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("cen_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				cen.NewCenService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("cen_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCenServiceRouteEntryService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteCenServiceRouteEntry",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"CenId":                ids[0],
				"DestinationCidrBlock": ids[1],
				"ServiceRegionId":      ids[2],
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading cen service route entry on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("cen_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCenServiceRouteEntryService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "service_route_entries",
	}
}

func (s *ByteplusCenServiceRouteEntryService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "cen",
		Action:      actionName,
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		RegionType:  bp.Global,
	}
}
