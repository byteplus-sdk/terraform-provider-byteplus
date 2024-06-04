package scaling_configuration

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

type ByteplusScalingConfigurationService struct {
	Client *bp.SdkClient
}

func NewScalingConfigurationService(c *bp.SdkClient) *ByteplusScalingConfigurationService {
	return &ByteplusScalingConfigurationService{
		Client: c,
	}
}

func (s *ByteplusScalingConfigurationService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusScalingConfigurationService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeScalingConfigurations"
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
		logger.Debug(logger.RespFormat, action, condition, resp)
		results, err = bp.ObtainSdkValue("Result.ScalingConfigurations", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ScalingConfigurations is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusScalingConfigurationService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ScalingConfigurationIds.1": id,
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
		return data, fmt.Errorf("ScalingConfiguration %s not exist ", id)
	}

	return data, err
}

func (s *ByteplusScalingConfigurationService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (ByteplusScalingConfigurationService) WithResourceResponseHandlers(scalingConfiguration map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return scalingConfiguration, map[string]bp.ResponseConvert{
			"Eip.Bandwidth": {
				TargetField: "eip_bandwidth",
			},
			"Eip.ISP": {
				TargetField: "eip_isp",
			},
			"Eip.BillingType": {
				TargetField: "eip_billing_type",
			},
		}, nil
	}

	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusScalingConfigurationService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	createConfigCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateScalingConfiguration",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"volumes": {
					ConvertType: bp.ConvertListN,
					ForceGet:    true,
				},
				"security_group_ids": {
					ConvertType: bp.ConvertWithN,
				},
				"instance_types": {
					ConvertType: bp.ConvertWithN,
				},
				"eip_bandwidth": {
					TargetField: "Eip.Bandwidth",
				},
				"eip_isp": {
					TargetField: "Eip.ISP",
				},
				"eip_billing_type": {
					TargetField: "Eip.BillingType",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if tags, ok := d.GetOk("tags"); ok {
					tagMap := map[string]interface{}{}
					for _, v := range tags.(*schema.Set).List() {
						if vMap, ok := v.(map[string]interface{}); ok {
							tagMap[vMap["key"].(string)] = vMap["value"]
						}
					}
					if tagsStr, err := json.Marshal(tagMap); err != nil {
						return false, err
					} else {
						(*call.SdkParam)["Tags"] = string(tagsStr)
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.ScalingConfigurationId", *resp)
				d.SetId(id.(string))
				return nil
			},
		},
	}
	callbacks = append(callbacks, createConfigCallback)

	return callbacks

}

func (s *ByteplusScalingConfigurationService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	// 修改伸缩配置
	modifyConfigurationCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyScalingConfiguration",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"scaling_configuration_name": {
					ConvertType: bp.ConvertDefault,
				},
				"image_id": {
					ConvertType: bp.ConvertDefault,
				},
				"instance_types": {
					ConvertType: bp.ConvertWithN,
				},
				"instance_name": {
					ConvertType: bp.ConvertDefault,
				},
				"instance_description": {
					ConvertType: bp.ConvertDefault,
				},
				"host_name": {
					ConvertType: bp.ConvertDefault,
				},
				"password": {
					ConvertType: bp.ConvertDefault,
				},
				"key_pair_name": {
					ConvertType: bp.ConvertDefault,
				},
				"key_pair_id": {
					ConvertType: bp.ConvertDefault,
				},
				"security_enhancement_strategy": {
					ConvertType: bp.ConvertDefault,
				},
				"user_data": {
					ConvertType: bp.ConvertDefault,
				},
				"volumes": {
					ConvertType: bp.ConvertListN,
				},
				"security_group_ids": {
					ConvertType: bp.ConvertWithN,
				},
				"eip_bandwidth": {
					TargetField: "Eip.Bandwidth",
				},
				"eip_isp": {
					TargetField: "Eip.ISP",
				},
				"eip_billing_type": {
					TargetField: "Eip.BillingType",
				},
				"project_name": {
					ConvertType: bp.ConvertDefault,
				},
				"spot_strategy": {
					ConvertType: bp.ConvertDefault,
				},
				"hpc_cluster_id": {
					ConvertType: bp.ConvertDefault,
				},
				"ipv6_address_count": {
					ConvertType: bp.ConvertDefault,
				},
			},
			RequestIdField: "ScalingConfigurationId",
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if d.HasChange("eip_bandwidth") || d.HasChange("eip_isp") || d.HasChange("eip_billing_type") {
					(*call.SdkParam)["Eip.Bandwidth"] = d.Get("eip_bandwidth")
					(*call.SdkParam)["Eip.ISP"] = d.Get("eip_isp")
					(*call.SdkParam)["Eip.BillingType"] = d.Get("eip_billing_type")
				}
				if d.HasChange("volumes") {
					for i, ele := range d.Get("volumes").([]interface{}) {
						volume := ele.(map[string]interface{})
						(*call.SdkParam)[fmt.Sprintf("Volumes.%d.DeleteWithInstance", i+1)] = volume["delete_with_instance"]
						(*call.SdkParam)[fmt.Sprintf("Volumes.%d.Size", i+1)] = volume["size"]
						(*call.SdkParam)[fmt.Sprintf("Volumes.%d.VolumeType", i+1)] = volume["volume_type"]
					}
				}
				if d.HasChange("tags") {
					if tags, ok := d.GetOk("tags"); ok {
						tagMap := map[string]interface{}{}
						for _, v := range tags.(*schema.Set).List() {
							if vMap, ok := v.(map[string]interface{}); ok {
								tagMap[vMap["key"].(string)] = vMap["value"]
							}
						}
						if tagsStr, err := json.Marshal(tagMap); err != nil {
							return false, err
						} else {
							(*call.SdkParam)["Tags"] = string(tagsStr)
						}
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, modifyConfigurationCallback)

	return callbacks
}

func (s *ByteplusScalingConfigurationService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteScalingConfiguration",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"ScalingConfigurationId": resourceData.Id(),
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
							return resource.NonRetryableError(fmt.Errorf("error on reading ScalingConfiguration on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.NonRetryableError(callErr)
				})
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusScalingConfigurationService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "ScalingConfigurationIds",
				ConvertType: bp.ConvertWithN,
			},
			"scaling_configuration_names": {
				TargetField: "ScalingConfigurationNames",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "ScalingConfigurationName",
		IdField:      "ScalingConfigurationId",
		CollectField: "scaling_configurations",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ScalingConfigurationId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"Eip.Bandwidth": {
				TargetField: "eip_bandwidth",
			},
			"Eip.ISP": {
				TargetField: "eip_isp",
			},
			"Eip.BillingType": {
				TargetField: "eip_billing_type",
			},
		},
	}
}

func (s *ByteplusScalingConfigurationService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "auto_scaling",
		Action:      actionName,
		Version:     "2020-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}
