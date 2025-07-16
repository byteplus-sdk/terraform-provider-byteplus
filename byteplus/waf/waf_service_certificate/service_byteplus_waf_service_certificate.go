package waf_service_certificate

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

type ByteplusWafServiceCertificateService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewWafServiceCertificateService(c *bp.SdkClient) *ByteplusWafServiceCertificateService {
	return &ByteplusWafServiceCertificateService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusWafServiceCertificateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusWafServiceCertificateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListWafServiceCertificate"

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
		results, err = bp.ObtainSdkValue("Result.Data", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Data is not Slice")
		}

		for _, ele := range data {
			certificate, ok := ele.(map[string]interface{})
			if !ok {
				return data, fmt.Errorf(" certificate is not Map ")
			}

			certificateId := int(certificate["Id"].(float64))

			certificate["CertificateId"] = strconv.Itoa(certificateId)

			logger.Debug(logger.ReqFormat, "CertificateId", certificate["CertificateId"])
		}
		return data, err
	})
}

func (s *ByteplusWafServiceCertificateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return nil, err
}

func (s *ByteplusWafServiceCertificateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusWafServiceCertificateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (ByteplusWafServiceCertificateService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusWafServiceCertificateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusWafServiceCertificateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusWafServiceCertificateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:        "Name",
		IdField:          "CertificateId",
		CollectField:     "data",
		ResponseConverts: map[string]bp.ResponseConvert{},
	}
}

func (s *ByteplusWafServiceCertificateService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "waf",
		Version:     "2023-12-25",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
		RegionType:  bp.Global,
	}
}
