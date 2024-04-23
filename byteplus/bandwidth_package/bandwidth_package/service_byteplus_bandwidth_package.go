package bandwidth_package

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusBandwidthPackageService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewBandwidthPackageService(c *bp.SdkClient) *ByteplusBandwidthPackageService {
	return &ByteplusBandwidthPackageService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusBandwidthPackageService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusBandwidthPackageService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		next    string
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithNextTokenQuery(m, "MaxResults", "NextToken", 20, nil, func(condition map[string]interface{}) ([]interface{}, string, error) {
		action := "DescribeBandwidthPackages"
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
		results, err = bp.ObtainSdkValue("Result.BandwidthPackages", *resp)
		if err != nil {
			return data, next, err
		}
		nextToken, err := bp.ObtainSdkValue("Result.NextToken", *resp)
		if err != nil {
			return data, next, err
		}
		next = nextToken.(string)
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, next, errors.New("Result.BandwidthPackages is not Slice")
		}
		return data, next, err
	})
}

func (s *ByteplusBandwidthPackageService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"BandwidthPackageIds.1": id,
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
		return data, fmt.Errorf("bandwidth_package %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusBandwidthPackageService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			failStates = append(failStates, "Failed")
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
					return nil, "", fmt.Errorf("bandwidth_package status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (s *ByteplusBandwidthPackageService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateBandwidthPackage",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"billing_type": {
					TargetField: "BillingType",
					Convert:     billingTypeRequestConvert,
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertListN,
				},
				"security_protection_types": {
					TargetField: "SecurityProtectionTypes",
					ConvertType: bp.ConvertWithN,
				},
				"isp": {
					TargetField: "ISP",
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.BandwidthPackageId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusBandwidthPackageService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"BillingType": {
				TargetField: "billing_type",
				Convert:     billingTypeResponseConvert,
			},
			"ISP": {
				TargetField: "isp",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusBandwidthPackageService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyBandwidthPackageAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"bandwidth_package_name": {
					TargetField: "BandwidthPackageName",
				},
				"description": {
					TargetField: "Description",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["BandwidthPackageId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)
	if resourceData.HasChange("bandwidth") {
		bandwidthCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyBandwidthPackageSpec",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"bandwidth": {
						TargetField: "Bandwidth",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["BandwidthPackageId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Available"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		}
		callbacks = append(callbacks, bandwidthCallback)
	}
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "bandwidthpackage", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)
	return callbacks
}

func (s *ByteplusBandwidthPackageService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteBandwidthPackage",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"BandwidthPackageId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusBandwidthPackageService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "BandwidthPackageIds",
				ConvertType: bp.ConvertWithN,
			},
			"isp": {
				TargetField: "ISP",
			},
			"tag_filters": {
				TargetField: "TagFilters",
				ConvertType: bp.ConvertListN,
				NextLevelConvert: map[string]bp.RequestConvert{
					"key": {
						TargetField: "Key",
					},
					"values": {
						TargetField: "Values",
						ConvertType: bp.ConvertWithN,
					},
				},
			},
		},
		NameField:    "BandwidthPackageName",
		IdField:      "BandwidthPackageId",
		CollectField: "packages",
		ResponseConverts: map[string]bp.ResponseConvert{
			"BandwidthPackageId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"BillingType": {
				TargetField: "billing_type",
				Convert:     billingTypeResponseConvert,
			},
			"ISP": {
				TargetField: "isp",
			},
		},
	}
}

func (s *ByteplusBandwidthPackageService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}

func (s *ByteplusBandwidthPackageService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "vpc",
		ResourceType:         "bandwidthpackage",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func (s *ByteplusBandwidthPackageService) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	if resourceData.Get("billing_type").(string) == "PrePaid" {
		info.NeedUnsubscribe = true
		info.Products = []string{"BandwidthPackage"}
	}
	return &info, nil
}
