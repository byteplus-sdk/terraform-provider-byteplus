package instance_state

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/redis/instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRedisInstanceStateService struct {
	Client *bp.SdkClient
}

func (v *ByteplusRedisInstanceStateService) GetClient() *bp.SdkClient {
	return v.Client
}

func (v *ByteplusRedisInstanceStateService) ReadResources(m map[string]interface{}) ([]interface{}, error) {
	return nil, nil
}

func (v *ByteplusRedisInstanceStateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	if id == "" {
		id = v.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, errors.New("id err")
	}
	data, err = instance.NewRedisDbInstanceService(v.Client).ReadResource(resourceData, ids[1])
	return data, err
}

func (v *ByteplusRedisInstanceStateService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, s string) *resource.StateChangeConf {
	return nil
}

func (v *ByteplusRedisInstanceStateService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return m, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (v *ByteplusRedisInstanceStateService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	action := ""
	if data.Get("action").(string) == "Restart" {
		action = "RestartDBInstance"
	}
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"action": {
					Ignore: true,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return v.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				instanceId := d.Get("instance_id").(string)
				logger.Debug(logger.RespFormat, call.Action, instanceId)
				d.SetId(fmt.Sprintf("state:%s", instanceId))
				return nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				instance.NewRedisDbInstanceService(v.Client): {
					Target:     []string{"Running"},
					Timeout:    data.Timeout(schema.TimeoutCreate),
					ResourceId: data.Get("instance_id").(string),
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (v *ByteplusRedisInstanceStateService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (v *ByteplusRedisInstanceStateService) RemoveResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return nil
}

func (v *ByteplusRedisInstanceStateService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (v *ByteplusRedisInstanceStateService) ReadResourceId(s string) string {
	return s
}

func NewRedisInstanceStateService(c *bp.SdkClient) *ByteplusRedisInstanceStateService {
	return &ByteplusRedisInstanceStateService{
		Client: c,
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Redis",
		Action:      actionName,
		Version:     "2020-12-07",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
	}
}
