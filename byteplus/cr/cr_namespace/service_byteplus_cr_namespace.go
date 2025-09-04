package cr_namespace

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCrNamespaceService struct {
	Client *bp.SdkClient
}

func NewCrNamespaceService(c *bp.SdkClient) *ByteplusCrNamespaceService {
	return &ByteplusCrNamespaceService{
		Client: c,
	}
}

func (s *ByteplusCrNamespaceService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCrNamespaceService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	pageCall := func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListNamespaces"

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

		logger.Debug(logger.RespFormat, action, condition, resp)
		results, err = bp.ObtainSdkValue("Result.Items", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}

		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Results.Items is not slice")
		}
		return data, err
	}

	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, pageCall)
}

func (s *ByteplusCrNamespaceService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return data, fmt.Errorf("the id format must be 'registry:namespace'")
	}
	registry := parts[0]
	name := parts[1]

	req := map[string]interface{}{
		"Registry": registry,
		"Filter": map[string]interface{}{
			"Names": []string{name},
		},
	}

	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}

	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not a map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("CrNamespace %s is not exist", id)
	}
	return data, err
}

func (s *ByteplusCrNamespaceService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusCrNamespaceService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCrNamespaceService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateNamespace",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"registry": {
					TargetField: "Registry",
				},
				"name": {
					TargetField: "Name",
				},
				"project": {
					TargetField: "Project",
				},
				"repository_default_access_level": {
					TargetField: "RepositoryDefaultAccessLevel",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				registry := d.Get("registry").(string)
				name := d.Get("name").(string)
				id := registry + ":" + name
				d.SetId(id)
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCrNamespaceService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusCrNamespaceService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteNamespace",
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				id := s.ReadResourceId(resourceData.Id())
				parts := strings.Split(id, ":")
				if len(parts) != 2 {
					return false, fmt.Errorf("the id format must be 'registry:namespace'")
				}
				(*call.SdkParam)["Registry"] = parts[0]
				(*call.SdkParam)["Name"] = parts[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 3*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCrNamespaceService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		IdField:      "Name",
		CollectField: "namespaces",
		RequestConverts: map[string]bp.RequestConvert{
			"names": {
				TargetField: "Filter.Names",
				ConvertType: bp.ConvertJsonArray,
			},
			"projects": {
				TargetField: "Filter.Projects",
				ConvertType: bp.ConvertJsonArray,
			},
		},
	}
}

func (s *ByteplusCrNamespaceService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusCrNamespaceService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "cr",
		ResourceType:         "repository",
		ProjectResponseField: "Project",
		ProjectSchemaField:   "project",
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "cr",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
