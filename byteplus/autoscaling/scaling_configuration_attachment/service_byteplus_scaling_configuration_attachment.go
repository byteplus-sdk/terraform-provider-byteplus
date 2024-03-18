package scaling_configuration_attachment

import (
	"errors"
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusScalingConfigurationAttachmentService struct {
	Client *bp.SdkClient
}

func (s *ByteplusScalingConfigurationAttachmentService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1,
		func(condition map[string]interface{}) ([]interface{}, error) {
			client := s.Client.UniversalClient
			action := "DescribeScalingConfigurations"
			logger.Debug(logger.ReqFormat, action, condition)
			if condition == nil {
				resp, err = client.DoCall(getUniversalInfo(action), nil)
				if err != nil {
					return data, err
				}
			} else {
				resp, err = client.DoCall(getUniversalInfo(action), &condition)
				if err != nil {
					return data, err
				}
			}
			logger.Debug(logger.RespFormat, action, condition, *resp)
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

func (s *ByteplusScalingConfigurationAttachmentService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, errors.New("Invalid ScalingConfigurationAttachment Id ")
	}
	req := map[string]interface{}{
		"ScalingConfigurationIds.1": ids[1],
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
		return data, fmt.Errorf("The ScalingConfiguration %s does not exist ", ids[1])
	}
	return data, err
}

func (s *ByteplusScalingConfigurationAttachmentService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, s2 string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusScalingConfigurationAttachmentService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *ByteplusScalingConfigurationAttachmentService) CreateResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var (
		readData map[string]interface{}
		err      error
		configId string
		groupId  string
	)
	configId = data.Get("scaling_configuration_id").(string)
	readData, err = s.ReadResource(data, fmt.Sprintf("enable:%s", configId))
	if err != nil {
		logger.DebugInfo("Failed to read scaling configuration resource", false)
		return []bp.Callback{}
	}
	groupId = readData["ScalingGroupId"].(string)
	logger.Debug(logger.RespFormat, "Read ScalingGroupId", configId, groupId)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "EnableScalingConfiguration",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"ScalingGroupId":         groupId,
				"ScalingConfigurationId": configId,
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprintf("enable:%s", (*call.SdkParam)["ScalingConfigurationId"]))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusScalingConfigurationAttachmentService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusScalingConfigurationAttachmentService) RemoveResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusScalingConfigurationAttachmentService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusScalingConfigurationAttachmentService) ReadResourceId(id string) string {
	return id
}

func NewScalingConfigurationAttachmentService(client *bp.SdkClient) *ByteplusScalingConfigurationAttachmentService {
	return &ByteplusScalingConfigurationAttachmentService{
		Client: client,
	}
}

func (s *ByteplusScalingConfigurationAttachmentService) GetClient() *bp.SdkClient {
	return s.Client
}

func getUniversalInfo(action string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "auto_scaling",
		Action:      action,
		Version:     "2020-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}
