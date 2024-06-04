package support_addon

import (
	"errors"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusVkeSupportAddonService struct {
	Client *bp.SdkClient
}

func NewVkeSupportAddonService(c *bp.SdkClient) *ByteplusVkeSupportAddonService {
	return &ByteplusVkeSupportAddonService{
		Client: c,
	}
}

func (s *ByteplusVkeSupportAddonService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusVkeSupportAddonService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	if _, ok := condition["Filter"]; ok {
		if kubernetesVersions, ok := condition["Filter"].(map[string]interface{})["KubernetesVersions"]; ok {
			condition["Filter"].(map[string]interface{})["Versions.Compatibilities.KubernetesVersions"] = kubernetesVersions
			delete(condition["Filter"].(map[string]interface{}), "KubernetesVersions")
		}
	}

	action := "ListSupportedAddons"
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

	results, err = bp.ObtainSdkValue("Result.Items", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		results = []interface{}{}
	}
	if data, ok = results.([]interface{}); !ok {
		return data, errors.New("Result.Items is not Slice")
	}
	return data, err
}

func (s *ByteplusVkeSupportAddonService) ReadResource(resourceData *schema.ResourceData, clusterId string) (data map[string]interface{}, err error) {
	return data, err
}

func (s *ByteplusVkeSupportAddonService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (ByteplusVkeSupportAddonService) WithResourceResponseHandlers(cluster map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *ByteplusVkeSupportAddonService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}

}

func (s *ByteplusVkeSupportAddonService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusVkeSupportAddonService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusVkeSupportAddonService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"name": {
				TargetField: "Filter.Name",
			},
			"pod_network_modes": {
				TargetField: "Filter.PodNetworkModes",
				ConvertType: bp.ConvertJsonArray,
			},
			"deploy_modes": {
				TargetField: "Filter.DeployModes",
				ConvertType: bp.ConvertJsonArray,
			},
			"deploy_node_types": {
				TargetField: "Filter.DeployNodeTypes",
				ConvertType: bp.ConvertJsonArray,
			},
			"necessaries": {
				TargetField: "Filter.Necessaries",
				ConvertType: bp.ConvertJsonArray,
			},
			"categories": {
				TargetField: "Filter.Categories",
				ConvertType: bp.ConvertJsonArray,
			},
			"kubernetes_versions": {
				TargetField: "Filter.KubernetesVersions",
				ConvertType: bp.ConvertJsonArray,
			},
		},
		ContentType:  bp.ContentTypeJson,
		NameField:    "Name",
		CollectField: "addons",
	}
}

func (s *ByteplusVkeSupportAddonService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vke",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
