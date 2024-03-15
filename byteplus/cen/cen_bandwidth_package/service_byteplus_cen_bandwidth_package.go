package cen_bandwidth_package

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCenBandwidthPackageService struct {
	Client *bp.SdkClient
}

func NewCenBandwidthPackageService(c *bp.SdkClient) *ByteplusCenBandwidthPackageService {
	return &ByteplusCenBandwidthPackageService{
		Client: c,
	}
}

func (s *ByteplusCenBandwidthPackageService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCenBandwidthPackageService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
		nameSet = make(map[string]bool)
	)
	if _, ok = m["CenBandwidthPackageNames.1"]; ok {
		i := 1
		for {
			filed := fmt.Sprintf("CenBandwidthPackageNames.%d", i)
			tmpName, ok := m[filed]
			if !ok {
				break
			}
			nameSet[tmpName.(string)] = true
			i++
			delete(m, filed)
		}
	}
	packages, err := bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeCenBandwidthPackages"
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
		results, err = bp.ObtainSdkValue("Result.CenBandwidthPackages", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.CenBandwidthPackages is not Slice")
		}
		return data, err
	})
	if err != nil || len(nameSet) == 0 {
		return packages, err
	}

	res := make([]interface{}, 0)
	for _, v := range packages {
		if !nameSet[v.(map[string]interface{})["CenBandwidthPackageName"].(string)] {
			continue
		}
		res = append(res, v)
	}
	return res, nil
}

func (s *ByteplusCenBandwidthPackageService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"CenBandwidthPackageIds.1": id,
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
		return data, fmt.Errorf("cen bandwidth package %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCenBandwidthPackageService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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

			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				demo, err = s.ReadResource(resourceData, id)
				if err != nil {
					if bp.ResourceNotFoundError(err) {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}
				return nil
			}); err != nil {
				return nil, "", err
			}

			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("cen status error, status:%s", status.(string))
				}
			}
			project, err := bp.ObtainSdkValue("ProjectName", demo)
			if err != nil {
				return nil, "", err
			}
			if resourceData.Get("project_name") != nil && resourceData.Get("project_name").(string) != "" {
				if project != resourceData.Get("project_name") {
					return demo, "", err
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (ByteplusCenBandwidthPackageService) WithResourceResponseHandlers(v map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return v, map[string]bp.ResponseConvert{
			"BillingType": {
				TargetField: "billing_type",
				Convert:     billingTypeResponseConvert,
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusCenBandwidthPackageService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateCenBandwidthPackage",
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
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.CenBandwidthPackageId", *resp)
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

func (s *ByteplusCenBandwidthPackageService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyCenBandwidthPackageAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"bandwidth": {
					ConvertType: bp.ConvertDefault,
				},
				"cen_bandwidth_package_name": {
					ConvertType: bp.ConvertDefault,
				},
				"description": {
					ConvertType: bp.ConvertDefault,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) == 0 {
					return false, nil
				}
				(*call.SdkParam)["CenBandwidthPackageId"] = d.Id()
				delete(*call.SdkParam, "Tags")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available", "InUse"},
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	callbacks = append(callbacks, callback)
	setResourceTagsCallbacks := bp.SetResourceTags(s.Client, "TagResources", "UntagResources", "cenbandwidthpackage", resourceData, getUniversalInfo)
	callbacks = append(callbacks, setResourceTagsCallbacks...)
	if resourceData.HasChange("project_name") {
		projectCallback := s.ModifyProject(resourceData)
		callbacks = append(callbacks, projectCallback...)
	}
	return callbacks
}

func (s *ByteplusCenBandwidthPackageService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteCenBandwidthPackage",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"CenBandwidthPackageId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
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
							return resource.NonRetryableError(fmt.Errorf("error on reading cen bandwidth package on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusCenBandwidthPackageService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "CenBandwidthPackageIds",
				ConvertType: bp.ConvertWithN,
			},
			"cen_bandwidth_package_names": {
				TargetField: "CenBandwidthPackageNames",
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
		NameField:    "CenBandwidthPackageName",
		IdField:      "CenBandwidthPackageId",
		CollectField: "bandwidth_packages",
		ResponseConverts: map[string]bp.ResponseConvert{
			"CenBandwidthPackageId": {
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

func (s *ByteplusCenBandwidthPackageService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "cen",
		Action:      actionName,
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}

func (s *ByteplusCenBandwidthPackageService) ModifyProject(resourceData *schema.ResourceData) []bp.Callback {
	var call []bp.Callback
	id := s.ReadResourceId(resourceData.Id())
	if resourceData.HasChange("project_name") {
		modifyProject := bp.Callback{
			Call: bp.SdkCall{
				Action:      "MoveProjectResource",
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"project_name": {
						ConvertType: bp.ConvertDefault,
						TargetField: "TargetProjectName",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if (*call.SdkParam)["TargetProjectName"] == nil || (*call.SdkParam)["TargetProjectName"] == "" {
						return false, fmt.Errorf("Could set ProjectName to empty ")
					}
					//获取用户ID
					input := map[string]interface{}{
						"ProjectName": (*call.SdkParam)["TargetProjectName"],
					}
					logger.Debug(logger.ReqFormat, "GetProject", input)
					out, err := s.Client.UniversalClient.DoCall(s.getIAMUniversalInfo("GetProject"), &input)
					if err != nil {
						return false, err
					}
					accountId, err := bp.ObtainSdkValue("Result.AccountID", *out)
					if err != nil {
						return false, err
					}
					trnStr := fmt.Sprintf("trn:%s:%s:%d:%s/%s", "cen", "", int(accountId.(float64)),
						"cenbandwidthpackage", id)
					(*call.SdkParam)["ResourceTrn.1"] = trnStr
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(s.getIAMUniversalInfo(call.Action), call.SdkParam)
				},
			},
		}
		call = append(call, modifyProject)
	}
	return call
}

func (s *ByteplusCenBandwidthPackageService) getIAMUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Version:     "2021-08-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}

func (s *ByteplusCenBandwidthPackageService) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	if resourceData.Get("billing_type").(string) == "PrePaid" {
		info.NeedUnsubscribe = true
		info.Products = []string{"CEN"}
	}
	return &info, nil
}
