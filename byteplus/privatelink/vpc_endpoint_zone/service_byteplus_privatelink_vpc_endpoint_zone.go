package vpc_endpoint_zone

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/privatelink/vpc_endpoint"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusVpcEndpointZoneService struct {
	Client *bp.SdkClient
}

func NewVpcEndpointZoneService(c *bp.SdkClient) *ByteplusVpcEndpointZoneService {
	return &ByteplusVpcEndpointZoneService{
		Client: c,
	}
}

func (s *ByteplusVpcEndpointZoneService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusVpcEndpointZoneService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 10, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeVpcEndpointZones"
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
		results, err = bp.ObtainSdkValue("Result.EndpointZones", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.EndpointZones is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusVpcEndpointZoneService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, errors.New("Invalid vpc endpoint zone id ")
	}
	endpointId := ids[0]
	subnetId := ids[1]
	req := map[string]interface{}{
		"EndpointId": endpointId,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		var zoneMap map[string]interface{}
		if zoneMap, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
		if subnetId == zoneMap["SubnetId"].(string) {
			data = zoneMap
			data["PrivateIpAddress"] = data["NetworkInterfaceIP"]
			data["EndpointId"] = endpointId
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Vpc endpoint zone %s not exist ", id)
	}

	return data, err
}

func (s *ByteplusVpcEndpointZoneService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusVpcEndpointZoneService) WithResourceResponseHandlers(data map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return data, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusVpcEndpointZoneService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	endpointId := resourceData.Get("endpoint_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddZoneToVpcEndpoint",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				// 查询子网所属的 zone_id
				zoneId, err := s.getZoneIdBySubnet(d.Get("subnet_id").(string))
				if err != nil {
					return false, err
				}
				if zoneId == "" {
					return false, fmt.Errorf(" Failed to obtain zone from subnet id: %v", d.Get("subnet_id"))
				}
				(*call.SdkParam)["ZoneId"] = zoneId
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id := fmt.Sprintf("%s:%s", endpointId, d.Get("subnet_id"))
				d.SetId(id)
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc_endpoint.NewVpcEndpointService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: endpointId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return endpointId
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *ByteplusVpcEndpointZoneService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusVpcEndpointZoneService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "RemoveZoneFromVpcEndpoint",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				if len(ids) != 2 {
					return false, errors.New("Invalid vpc endpoint zone id ")
				}
				endpointId := ids[0]
				subnetId := ids[1]

				// 查询子网所属的 zone_id
				zoneId, err := s.getZoneIdBySubnet(subnetId)
				if err != nil {
					return false, err
				}
				if zoneId == "" {
					return false, fmt.Errorf(" Failed to obtain zone from subnet id: %v", d.Get("subnet_id"))
				}

				(*call.SdkParam)["EndpointId"] = endpointId
				(*call.SdkParam)["ZoneId"] = zoneId
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				vpc_endpoint.NewVpcEndpointService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutDelete),
					ResourceId: resourceData.Get("endpoint_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("endpoint_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusVpcEndpointZoneService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		IdField:      "ZoneId",
		CollectField: "vpc_endpoint_zones",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ZoneId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"NetworkInterfaceIP": {
				TargetField: "network_interface_ip",
			},
		},
	}
}

func (s *ByteplusVpcEndpointZoneService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusVpcEndpointZoneService) getZoneIdBySubnet(subnetId string) (zoneId string, err error) {
	action := "DescribeSubnets"
	req := map[string]interface{}{
		"SubnetIds.1": subnetId,
	}
	resp, err := s.Client.UniversalClient.DoCall(getVpcUniversalInfo(action), &req)
	if err != nil {
		return "", err
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	results, err := bp.ObtainSdkValue("Result.Subnets", *resp)
	if err != nil {
		return "", err
	}
	if results == nil {
		results = []interface{}{}
	}
	subnets, ok := results.([]interface{})
	if !ok {
		return "", errors.New("Result.Subnets is not Slice")
	}
	if len(subnets) == 0 {
		return "", fmt.Errorf("subnet %s not exist", subnetId)
	}
	zoneId = subnets[0].(map[string]interface{})["ZoneId"].(string)
	return zoneId, nil
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

func getVpcUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
