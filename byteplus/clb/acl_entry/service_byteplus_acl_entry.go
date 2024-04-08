package acl_entry

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/acl"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusAclEntryService struct {
	Client *bp.SdkClient
}

func NewAclEntryService(c *bp.SdkClient) *ByteplusAclEntryService {
	return &ByteplusAclEntryService{
		Client: c,
	}
}

func (s *ByteplusAclEntryService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusAclEntryService) ReadResources(condition map[string]interface{}) ([]interface{}, error) {
	var (
		aclEntryIdMap = make(map[string]bool)
		res           = make([]interface{}, 0)
	)

	aclEntries, err := bp.WithSimpleQuery(condition, func(m map[string]interface{}) ([]interface{}, error) {
		var (
			resp    *map[string]interface{}
			err     error
			results interface{}
		)
		action := "DescribeAclAttributes"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return []interface{}{}, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return []interface{}{}, err
			}
		}

		results, err = bp.ObtainSdkValue("Result.AclEntries", *resp)
		if err != nil {
			return []interface{}{}, err
		}
		if _, ok := results.([]interface{}); !ok {
			return []interface{}{}, errors.New("Result.AclEntries is not Slice")
		}
		return results.([]interface{}), err
	})

	if err != nil {
		return aclEntries, err
	}

	aclEntryIds, ok := condition["AclEntryIds"]
	if !ok || aclEntryIds == nil {
		return aclEntries, nil
	}

	if reflect.TypeOf(aclEntryIds).Kind() != reflect.Slice {
		return []interface{}{}, fmt.Errorf("condition[\"AclEntryIds\"] is not Slice")
	}

	for _, entry := range aclEntryIds.([]string) {
		aclEntryIdMap[entry] = true
	}

	if len(aclEntryIdMap) == 0 {
		return aclEntries, nil
	}

	for _, aclEntry := range aclEntries {
		if _, ok := aclEntryIdMap[aclEntry.(map[string]interface{})["Entry"].(string)]; ok {
			res = append(res, aclEntry)
		}
	}
	return res, nil
}

func (s *ByteplusAclEntryService) ReadResource(resourceData *schema.ResourceData, tmpId string) (map[string]interface{}, error) {
	if tmpId == "" {
		tmpId = resourceData.Id()
	}

	ids := strings.Split(tmpId, ":")
	if len(ids) != 2 {
		return map[string]interface{}{}, fmt.Errorf("invalid acl entry id")
	}

	aclId := ids[0]
	entry := ids[1]
	req := map[string]interface{}{
		"AclId":       aclId,
		"AclEntryIds": []string{entry},
	}

	aclEntries, err := s.ReadResources(req)
	if err != nil {
		return nil, err
	}

	if len(aclEntries) == 0 {
		return map[string]interface{}{}, fmt.Errorf("acl entry %s:%s not exist ", aclId, entry)
	}

	for _, v := range aclEntries {
		if _, ok := v.(map[string]interface{}); !ok {
			return map[string]interface{}{}, errors.New("Value is not map ")
		}
	}

	return aclEntries[0].(map[string]interface{}), err
}

func (s *ByteplusAclEntryService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (ByteplusAclEntryService) WithResourceResponseHandlers(aclEntry map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return aclEntry, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusAclEntryService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AddAclEntries",
			ConvertMode: bp.RequestConvertAll,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["AclId"] = d.Get("acl_id")
				(*call.SdkParam)["AclEntries.1.Entry"] = d.Get("entry")
				(*call.SdkParam)["AclEntries.1.Description"] = d.Get("description")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id := fmt.Sprintf("%s:%s", d.Get("acl_id"), d.Get("entry"))
				d.SetId(id)
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("acl_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				acl.NewAclService(s.Client): {
					Target:     []string{"Active"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("acl_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *ByteplusAclEntryService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusAclEntryService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "RemoveAclEntries",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				tmpId := d.Id()
				aclEntryId := strings.Split(tmpId, ":")
				if len(aclEntryId) != 2 {
					return false, fmt.Errorf("error acl entry id: %s", tmpId)
				}
				(*call.SdkParam)["AclId"] = aclEntryId[0]
				(*call.SdkParam)["Entries.1"] = aclEntryId[1]
				return true, nil
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
							return resource.NonRetryableError(fmt.Errorf("error on  reading acl entry on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("acl_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				acl.NewAclService(s.Client): {
					Target:     []string{"Active"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("acl_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAclEntryService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusAclEntryService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "clb",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
