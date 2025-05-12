package alb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusAlbService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewAlbService(c *bp.SdkClient) *ByteplusAlbService {
	return &ByteplusAlbService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusAlbService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusAlbService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	data, err = bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeLoadBalancers"

		bytes, _ := json.Marshal(condition)
		logger.Debug(logger.ReqFormat, action, string(bytes))
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
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))
		results, err = bp.ObtainSdkValue("Result.LoadBalancers", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.LoadBalancers is not Slice")
		}
		data, err = removeSystemTags(data)
		return data, err
	})
	if err != nil {
		return data, err
	}

	for _, value := range data {
		alb, ok := value.(map[string]interface{})
		if !ok {
			return data, fmt.Errorf("Alb is not map ")
		}

		detailAction := "DescribeLoadBalancerAttributes"
		req := map[string]interface{}{
			"LoadBalancerId": alb["LoadBalancerId"],
		}
		logger.Debug(logger.ReqFormat, detailAction, req)
		detailResp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(detailAction), &req)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, detailAction, *detailResp)

		listeners, err := bp.ObtainSdkValue("Result.Listeners", *detailResp)
		if err != nil {
			return data, err
		}
		alb["Listeners"] = listeners

		accessLog, err := bp.ObtainSdkValue("Result.AccessLog", *detailResp)
		if err != nil {
			return data, err
		}
		alb["AccessLog"] = accessLog

		tlsAccessLog, err := bp.ObtainSdkValue("Result.TLSAccessLog", *detailResp)
		if err != nil {
			return data, err
		}
		alb["TLSAccessLog"] = tlsAccessLog

		healthLog, err := bp.ObtainSdkValue("Result.HealthLog", *detailResp)
		if err != nil {
			return data, err
		}
		alb["HealthLog"] = healthLog
	}

	return data, err
}

func (s *ByteplusAlbService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"LoadBalancerIds.1": id,
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
		return data, fmt.Errorf("alb %s not exist ", id)
	}

	var subnetIds []interface{}
	zoneMappings, ok := data["ZoneMappings"].([]interface{})
	if !ok {
		return data, fmt.Errorf("ZoneMappings is not slice ")
	}
	for _, v := range zoneMappings {
		zoneMap, ok := v.(map[string]interface{})
		if !ok {
			return data, fmt.Errorf("Zone is not map ")
		}
		subnetIds = append(subnetIds, zoneMap["SubnetId"])
		addresses, ok := zoneMap["LoadBalancerAddresses"].([]interface{})
		if !ok || len(addresses) == 0 {
			return data, fmt.Errorf("LoadBalancerAddresses is not slice ")
		}
		addressMap, ok := addresses[0].(map[string]interface{})
		if !ok {
			return data, fmt.Errorf("LoadBalancerAddresse is not map ")
		}

		if _, exist := addressMap["Eip"]; exist && addressMap["Eip"] != nil {
			eip, ok := addressMap["Eip"].(map[string]interface{})
			if !ok {
				return data, fmt.Errorf("Eip of LoadBalancerAddresse is not map ")
			}
			data["EipBillingConfig"] = eip
		}
		if _, exist := addressMap["Ipv6Eip"]; exist && addressMap["Ipv6Eip"] != nil {
			ipv6Eip, ok := addressMap["Ipv6Eip"].(map[string]interface{})
			if !ok {
				return data, fmt.Errorf("Ipv6Eip of LoadBalancerAddresse is not map ")
			}
			data["Ipv6EipBillingConfig"] = ipv6Eip
		}
	}
	data["SubnetIds"] = subnetIds

	return data, err
}

func (ByteplusAlbService) WithResourceResponseHandlers(alb map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return alb, map[string]bp.ResponseConvert{
			"ISP": {
				TargetField: "isp",
			},
			"EipBillingType": {
				TargetField: "eip_billing_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					billingType := i.(float64)
					switch billingType {
					case 2:
						return "PostPaidByBandwidth"
					case 3:
						return "PostPaidByTraffic"
					}
					return ""
				},
			},
			"BillingType": {
				TargetField: "billing_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					billingType := i.(float64)
					switch billingType {
					case 2:
						return "PostPaidByBandwidth"
					case 3:
						return "PostPaidByTraffic"
					}
					return ""
				},
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusAlbService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d          map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "CreateFailed")
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("alb status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (s *ByteplusAlbService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateLoadBalancer",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"subnet_ids": {
					Ignore: true,
				},
				"eip_billing_config": {
					TargetField: "EipBillingConfig",
					ConvertType: bp.ConvertListUnique,
					NextLevelConvert: map[string]bp.RequestConvert{
						"isp": {
							TargetField: "ISP",
						},
					},
				},
				"ipv6_eip_billing_config": {
					TargetField: "Ipv6EipBillingConfig",
					ConvertType: bp.ConvertListUnique,
					NextLevelConvert: map[string]bp.RequestConvert{
						"isp": {
							TargetField: "ISP",
						},
					},
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["RegionId"] = s.Client.Region
				(*call.SdkParam)["LoadBalancerBillingType"] = 1

				subnetIds, ok := d.Get("subnet_ids").(*schema.Set)
				if !ok {
					return false, fmt.Errorf("SubnetIds is not set")
				}
				vpcId, zoneIds, err := s.getVpcIdAndZoneIdBySubnets(subnetIds.List())
				if err != nil {
					return false, err
				}
				if subnetIds.Len() != len(zoneIds) {
					return false, fmt.Errorf("The length of subnetIds and zoneIds are not equal ")
				}
				for index, subnetId := range subnetIds.List() {
					(*call.SdkParam)["ZoneMappings."+strconv.Itoa(index+1)+".SubnetId"] = subnetId.(string)
					(*call.SdkParam)["ZoneMappings."+strconv.Itoa(index+1)+".ZoneId"] = zoneIds[subnetId.(string)]
				}
				(*call.SdkParam)["VpcId"] = vpcId

				// private 类型不传 eip_billing_config / ipv6_eip_billing_config
				if (*call.SdkParam)["Type"] == "private" {
					delete(*call.SdkParam, "EipBillingConfig.ISP")
					delete(*call.SdkParam, "EipBillingConfig.EipBillingType")
					delete(*call.SdkParam, "EipBillingConfig.Bandwidth")
					delete(*call.SdkParam, "Ipv6EipBillingConfig.ISP")
					delete(*call.SdkParam, "Ipv6EipBillingConfig.BillingType")
					delete(*call.SdkParam, "Ipv6EipBillingConfig.Bandwidth")
				}

				if eipBillingType, exist := (*call.SdkParam)["EipBillingConfig.EipBillingType"]; exist {
					ty := 0
					switch eipBillingType.(string) {
					case "PostPaidByBandwidth":
						ty = 2
					case "PostPaidByTraffic":
						ty = 3
					}
					(*call.SdkParam)["EipBillingConfig.EipBillingType"] = ty
				}
				if ipv6BillingType, exist := (*call.SdkParam)["Ipv6EipBillingConfig.BillingType"]; exist {
					ty := 0
					switch ipv6BillingType.(string) {
					case "PostPaidByBandwidth":
						ty = 2
					case "PostPaidByTraffic":
						ty = 3
					}
					(*call.SdkParam)["Ipv6EipBillingConfig.BillingType"] = ty
				}

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.LoadBalancerId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyLoadBalancerAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"load_balancer_name": {
					TargetField: "LoadBalancerName",
				},
				"description": {
					TargetField: "Description",
				},
				"delete_protection": {
					TargetField: "DeleteProtection",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) > 0 {
					(*call.SdkParam)["LoadBalancerId"] = d.Id()
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	// 更新Tags
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "loadbalancer", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)

	return callbacks
}

func (s *ByteplusAlbService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteLoadBalancer",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"LoadBalancerId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				if protection, ok := d.Get("delete_protection").(string); ok && protection == "on" {
					// 开启删除保护，直接返回接口报错
					return baseErr
				}

				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading alb on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "LoadBalancerIds",
				ConvertType: bp.ConvertWithN,
			},
			"tags": {
				TargetField: "TagFilters",
				ConvertType: bp.ConvertListN,
				NextLevelConvert: map[string]bp.RequestConvert{
					"value": {
						TargetField: "Values.1",
					},
				},
			},
		},
		NameField:    "LoadBalancerName",
		IdField:      "LoadBalancerId",
		CollectField: "albs",
		ResponseConverts: map[string]bp.ResponseConvert{
			"LoadBalancerId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"ISP": {
				TargetField: "isp",
			},
			"EipBillingType": {
				TargetField: "eip_billing_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					billingType := i.(float64)
					switch billingType {
					case 2:
						return "PostPaidByBandwidth"
					case 3:
						return "PostPaidByTraffic"
					}
					return ""
				},
			},
			"BillingType": {
				TargetField: "billing_type",
				Convert: func(i interface{}) interface{} {
					if i == nil {
						return nil
					}
					billingType := i.(float64)
					switch billingType {
					case 2:
						return "PostPaidByBandwidth"
					case 3:
						return "PostPaidByTraffic"
					}
					return ""
				},
			},
		},
	}
}

func (s *ByteplusAlbService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusAlbService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "alb",
		ResourceType:         "loadbalancer",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func (s *ByteplusAlbService) getVpcIdAndZoneIdBySubnets(subnetIds []interface{}) (string, map[string]string, error) {
	// describe subnet
	req := make(map[string]interface{})
	for index, subnetId := range subnetIds {
		req["SubnetIds."+strconv.Itoa(index+1)] = subnetId.(string)
	}
	action := "DescribeSubnets"
	resp, err := s.Client.UniversalClient.DoCall(getVpcUniversalInfo(action), &req)
	if err != nil {
		return "", map[string]string{}, err
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	results, err := bp.ObtainSdkValue("Result.Subnets", *resp)
	if err != nil {
		return "", map[string]string{}, err
	}
	if results == nil {
		results = []interface{}{}
	}
	subnets, ok := results.([]interface{})
	if !ok {
		return "", map[string]string{}, errors.New("Result.Subnets is not Slice")
	}
	if len(subnets) == 0 {
		return "", map[string]string{}, fmt.Errorf("subnets %v not exist", subnetIds)
	}
	zoneIds := make(map[string]string, 0)
	for _, v := range subnets {
		subnet, ok := v.(map[string]interface{})
		if !ok {
			return "", map[string]string{}, fmt.Errorf("Result.Subnet is not map")
		}
		zoneIds[subnet["SubnetId"].(string)] = subnet["ZoneId"].(string)
	}
	vpcId := subnets[0].(map[string]interface{})["VpcId"].(string)
	return vpcId, zoneIds, nil
}

func removeSystemTags(data []interface{}) ([]interface{}, error) {
	var (
		ok      bool
		result  map[string]interface{}
		results []interface{}
		tags    []interface{}
	)
	for _, d := range data {
		if result, ok = d.(map[string]interface{}); !ok {
			return results, errors.New("The elements in data are not map ")
		}
		tags, ok = result["Tags"].([]interface{})
		if ok {
			tags = bp.FilterSystemTags(tags)
			result["Tags"] = tags
		}
		results = append(results, result)
	}
	return results, nil
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "alb",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}

func getVpcUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
