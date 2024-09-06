package vpn_gateway

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusVpnGatewayService struct {
	Client *bp.SdkClient
}

func NewVpnGatewayService(c *bp.SdkClient) *ByteplusVpnGatewayService {
	return &ByteplusVpnGatewayService{
		Client: c,
	}
}

func (s *ByteplusVpnGatewayService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusVpnGatewayService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
		nameSet = make(map[string]bool)
	)
	if _, ok = m["VpnGatewayNames.1"]; ok {
		i := 1
		for {
			filed := fmt.Sprintf("VpnGatewayNames.%d", i)
			tmpName, ok := m[filed]
			if !ok {
				break
			}
			nameSet[tmpName.(string)] = true
			i++
			delete(m, filed)
		}
	}
	gateways, err := bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeVpnGateways"
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
		logger.Debug(logger.RespFormat, action, *resp)
		results, err = bp.ObtainSdkValue("Result.VpnGateways", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.VpnGateways is not Slice")
		}
		return data, err
	})
	if err != nil || len(nameSet) == 0 {
		return gateways, err
	}

	res := make([]interface{}, 0)
	for _, gateway := range gateways {
		if !nameSet[gateway.(map[string]interface{})["VpnGatewayName"].(string)] {
			continue
		}
		res = append(res, gateway)
	}
	return res, nil
}

func (s *ByteplusVpnGatewayService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"VpnGatewayIds.1": id,
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
		return data, fmt.Errorf("VpnGateway %s not exist ", id)
	}

	// 计费信息
	action := "DescribeVpnGatewaysBilling"
	params := &map[string]interface{}{
		"VpnGatewayIds.1": id,
	}
	billingRes, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), params)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.AllFormat, "DescribeVpnGatewaysBilling", params, *billingRes)
	tmpRes, err := bp.ObtainSdkValue("Result.VpnGateways", *billingRes)
	if err != nil {
		return data, err
	}
	if tmpRes == nil {
		return data, errors.New("Result.VpnGateways is not nil")
	}
	tmpData, ok := tmpRes.([]interface{})
	if !ok {
		return data, errors.New("Result.VpnGateways is not Slice")
	}
	if len(tmpData) == 0 {
		return data, fmt.Errorf("VpnGatewaysBilling %s not exist ", id)
	}
	data["RenewType"] = tmpData[0].(map[string]interface{})["RenewType"]
	data["RemainRenewTimes"] = int(tmpData[0].(map[string]interface{})["RemainRenewTimes"].(float64))

	return data, err
}

func (s *ByteplusVpnGatewayService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
					return nil, "", fmt.Errorf("VpnGateway  status  error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (ByteplusVpnGatewayService) WithResourceResponseHandlers(v map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		if v["BillingType"].(float64) == 1 {
			var (
				ct time.Time
				et time.Time
			)
			if strings.Contains(v["CreationTime"].(string), "+") {
				ct, _ = time.Parse("2006-01-02T15:04:05", v["CreationTime"].(string)[0:strings.Index(v["CreationTime"].(string), "+")])
			} else {
				ct, _ = time.Parse("2006-01-02 15:04:05", v["CreationTime"].(string))
			}
			if strings.Contains(v["ExpiredTime"].(string), "+") {
				et, _ = time.Parse("2006-01-02T15:04:05", v["ExpiredTime"].(string)[0:strings.Index(v["ExpiredTime"].(string), "+")])
			} else {
				et, _ = time.Parse("2006-01-02 15:04:05", v["ExpiredTime"].(string))
			}
			y := et.Year() - ct.Year()
			m := et.Month() - ct.Month()
			v["Period"] = y*12 + int(m)
		}
		return v, map[string]bp.ResponseConvert{
			"BillingType": {
				TargetField: "billing_type",
				Convert:     billingTypeResponseConvert,
			},
			"RenewType": {
				TargetField: "renew_type",
				Convert:     renewTypeResponseConvert,
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusVpnGatewayService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)

	// 创建VpnGateway
	createVpnGateway := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateVpnGateway",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"bandwidth": {
					ConvertType: bp.ConvertDefault,
				},
				"description": {
					ConvertType: bp.ConvertDefault,
				},
				"period": {
					ConvertType: bp.ConvertDefault,
				},
				"period_unit": {
					ConvertType: bp.ConvertDefault,
				},
				"subnet_id": {
					ConvertType: bp.ConvertDefault,
				},
				"vpc_id": {
					ConvertType: bp.ConvertDefault,
				},
				"vpn_gateway_name": {
					ConvertType: bp.ConvertDefault,
				},
				"billing_type": {
					TargetField: "BillingType",
					Convert:     billingTypeRequestConvert,
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
				"project_name": {
					ConvertType: bp.ConvertDefault,
				},
				"ssl_max_connections": {
					ConvertType: bp.ConvertDefault,
				},
				"ssl_enabled": {
					ConvertType: bp.ConvertDefault,
					ForceGet:    true,
				},
				"ipsec_enabled": {
					ConvertType: bp.ConvertDefault,
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) < 1 {
					return false, nil
				}
				(*call.SdkParam)["PeriodUnit"] = "Month"
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.VpnGatewayId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, createVpnGateway)

	return callbacks

}

func (s *ByteplusVpnGatewayService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)

	// 修改vpnGateway
	modifyCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyVpnGatewayAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"vpn_gateway_name": {
					ConvertType: bp.ConvertDefault,
				},
				"description": {
					ConvertType: bp.ConvertDefault,
				},
				"bandwidth": {
					ConvertType: bp.ConvertDefault,
				},
				"ssl_enabled": {
					ConvertType: bp.ConvertDefault,
					ForceGet:    true,
				},
				"ipsec_enabled": {
					ConvertType: bp.ConvertDefault,
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) < 1 {
					return false, nil
				}
				(*call.SdkParam)["VpnGatewayId"] = d.Id()
				// 只有为 true 的时候，强制加上去
				if d.Get("ssl_enabled").(bool) {
					(*call.SdkParam)["SslMaxConnections"] = d.Get("ssl_max_connections")
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				// 修改完成后，需要等待一定时间
				time.Sleep(5 * time.Second)
				return nil
			},
		},
	}
	callbacks = append(callbacks, modifyCallback)

	// 续费时长
	if resourceData.Get("renew_type").(string) == "ManualRenew" && resourceData.HasChange("period") {
		renewVpnGateway := bp.Callback{
			Call: bp.SdkCall{
				Action:      "RenewVpnGateway",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"period": {
						ConvertType: bp.ConvertDefault,
						Convert: func(data *schema.ResourceData, i interface{}) interface{} {
							o, n := data.GetChange("period")
							return n.(int) - o.(int)
						},
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 0 {
						(*call.SdkParam)["PeriodUnit"] = "Month"
						(*call.SdkParam)["VpnGatewayId"] = d.Id()
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Available"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, renewVpnGateway)
	}

	return callbacks
}

func (s *ByteplusVpnGatewayService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteVpnGateway",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"VpnGatewayId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				// todo 打印前台提示日志
				log.Println("[WARN] Terraform will unsubscribe the resource.")
				//return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				return nil, nil
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading VpnGateway on delete %q, %w", d.Id(), callErr))
						}
					}
					resp, callErr := call.ExecuteCall(d, client, call)
					logger.Debug(logger.AllFormat, call.Action, call.SdkParam, resp, callErr)
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

func (s *ByteplusVpnGatewayService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "VpnGatewayIds",
				ConvertType: bp.ConvertWithN,
			},
			"vpn_gateway_names": {
				TargetField: "VpnGatewayNames",
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
		NameField:    "VpnGatewayName",
		IdField:      "VpnGatewayId",
		CollectField: "vpn_gateways",
		ResponseConverts: map[string]bp.ResponseConvert{
			"VpnGatewayId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"BillingType": {
				TargetField: "billing_type",
				Convert:     billingTypeResponseConvert,
			},
		},
	}
}

func (s *ByteplusVpnGatewayService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpn",
		Action:      actionName,
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}

func (s *ByteplusVpnGatewayService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "vpn",
		ResourceType:         "vpngateway",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func (s *ByteplusVpnGatewayService) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	info.NeedUnsubscribe = true
	info.Products = []string{"VPN"}
	return &info, nil
}
