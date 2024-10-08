package byteplus

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]terraform.ResourceProvider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"byteplus": testAccProvider,
	}
}

type AccTestResource struct {
	ResourceId  string
	Svc         common.ResourceService
	ClientField string
	SvcInitFunc func(*common.SdkClient) common.ResourceService
}

func GetTestAccProvider() *schema.Provider {
	return testAccProvider
}

func GetTestAccProviders() map[string]terraform.ResourceProvider {
	return testAccProviders
}

func AccTestCheckResourceExists(acc *AccTestResource) resource.TestCheckFunc {
	if acc.ClientField == "" {
		acc.ClientField = "Client"
	}
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[acc.ResourceId]
		if !ok {
			return fmt.Errorf("resource %s is not found", acc.ResourceId)
		}

		if err := initResourceService(acc); err != nil {
			return err
		}

		resourceData, err := buildResourceData(rs)
		if err != nil {
			return err
		}

		out, err := acc.Svc.ReadResource(resourceData, acc.Svc.ReadResourceId(rs.Primary.ID))
		if err != nil {
			return err
		}
		if len(out) == 0 {
			return fmt.Errorf("resource %s is not found", acc.ResourceId)
		}

		return nil
	}
}

func initResourceService(acc *AccTestResource) error {
	if acc.Svc == nil || (reflect.ValueOf(acc.Svc).Kind() == reflect.Ptr && reflect.ValueOf(acc.Svc).IsNil()) {
		if acc.SvcInitFunc != nil {
			acc.Svc = acc.SvcInitFunc(testAccProvider.Meta().(*common.SdkClient))
		} else {
			return fmt.Errorf(" neither acc.Svc nor acc.SvcInitFunc is specified ")
		}
	} else {
		it := reflect.ValueOf(acc.Svc).Elem()
		val := it.FieldByName(acc.ClientField)
		if !val.IsNil() {
			return nil
		} else {
			val.Set(reflect.ValueOf(testAccProvider.Meta()))
		}
	}
	return nil
}

func buildResourceData(rs *terraform.ResourceState) (*schema.ResourceData, error) {
	var resourceSchema map[string]*schema.Schema
	if v, ok := testAccProvider.DataSourcesMap[rs.Type]; ok {
		resourceSchema = v.Schema
	}
	if v, ok := testAccProvider.ResourcesMap[rs.Type]; ok {
		resourceSchema = v.Schema
	}
	if resourceSchema == nil {
		return nil, fmt.Errorf("cannot get schema from provider")
	}

	// build resource data
	resourceData, err := schema.InternalMap(resourceSchema).Data(rs.Primary, nil)
	if err != nil {
		return nil, err
	}
	return resourceData, nil
}

func AccTestCheckResourceRemove(acc *AccTestResource) resource.TestCheckFunc {
	if acc.ClientField == "" {
		acc.ClientField = "Client"
	}
	return func(state *terraform.State) error {
		return resource.Retry(15*time.Minute, func() *resource.RetryError {
			rs, ok := state.RootModule().Resources[acc.ResourceId]
			if !ok {
				return resource.NonRetryableError(fmt.Errorf("resource %s is not found", acc.ResourceId))
			}

			if err := initResourceService(acc); err != nil {
				return resource.NonRetryableError(fmt.Errorf("resource service initialize error %s", err.Error()))
			}

			resourceData, err := buildResourceData(rs)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("resource data build error %s", err.Error()))
			}

			out, err := acc.Svc.ReadResource(resourceData, acc.Svc.ReadResourceId(rs.Primary.ID))
			if err != nil {
				if common.ResourceNotFoundError(err) {
					return nil
				} else {
					return resource.RetryableError(fmt.Errorf("retry check reomve %s", acc.ResourceId))
				}
			}
			if len(out) == 0 {
				return nil
			} else {
				return resource.RetryableError(fmt.Errorf("retry check reomve %s", acc.ResourceId))
			}
		})
	}
}

func AccTestPreCheck(t *testing.T) {
	if v := os.Getenv("BYTEPLUS_ACCESS_KEY"); v == "" {
		t.Fatal("BYTEPLUS_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("BYTEPLUS_SECRET_KEY"); v == "" {
		t.Fatal("BYTEPLUS_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("BYTEPLUS_REGION"); v == "" {
		log.Println("[INFO] Test: Using ap-southeast-1 as test region")
		os.Setenv("BYTEPLUS_REGION", "ap-southeast-1")
	}
}

// for terraform set type check
const (
	sentinelIndex = "*"
)

func TestCheckTypeSetElemAttr(name, attr, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		is, err := primaryInstanceState(s, name)
		if err != nil {
			return err
		}

		err = testCheckTypeSetElem(is, attr, value)
		if err != nil {
			return fmt.Errorf("%q error: %s", name, err)
		}

		return nil
	}
}

// primaryInstanceState returns the primary instance state for the given
// resource name in the root module.
func primaryInstanceState(s *terraform.State, name string) (*terraform.InstanceState, error) {
	ms := s.RootModule()
	return modulePrimaryInstanceState(ms, name)
}

func modulePrimaryInstanceState(ms *terraform.ModuleState, name string) (*terraform.InstanceState, error) {
	rs, ok := ms.Resources[name]
	if !ok {
		return nil, fmt.Errorf("not found: %s in %s", name, ms.Path)
	}

	is := rs.Primary
	if is == nil {
		return nil, fmt.Errorf("no primary instance: %s in %s", name, ms.Path)
	}

	return is, nil
}

func testCheckTypeSetElem(is *terraform.InstanceState, attr, value string) error {
	attrParts := strings.Split(attr, ".")
	if attrParts[len(attrParts)-1] != sentinelIndex {
		return fmt.Errorf("%q does not end with the special value %q", attr, sentinelIndex)
	}
	for stateKey, stateValue := range is.Attributes {
		if stateValue == value {
			stateKeyParts := strings.Split(stateKey, ".")
			if len(stateKeyParts) == len(attrParts) {
				for i := range attrParts {
					if attrParts[i] != stateKeyParts[i] && attrParts[i] != sentinelIndex {
						break
					}
					if i == len(attrParts)-1 {
						return nil
					}
				}
			}
		}
	}
	return fmt.Errorf("no TypeSet element %q, with value %q in state: %#v", attr, value, is.Attributes)
}

func TestCheckTypeSetElemNestedAttrs(name, attr string, values map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		is, err := primaryInstanceState(s, name)
		if err != nil {
			return err
		}

		attrParts := strings.Split(attr, ".")
		if attrParts[len(attrParts)-1] != sentinelIndex {
			return fmt.Errorf("%q does not end with the special value %q", attr, sentinelIndex)
		}
		// account for cases where the user is trying to see if the value is unset/empty
		// there may be ambiguous scenarios where a field was deliberately unset vs set
		// to the empty string, this will match both, which may be a false positive.
		var matchCount int
		for _, v := range values {
			if v != "" {
				matchCount++
			}
		}
		if matchCount == 0 {
			return fmt.Errorf("%#v has no non-empty values", values)
		}

		if testCheckTypeSetElemNestedAttrsInState(is, attrParts, matchCount, values) {
			return nil
		}
		return fmt.Errorf("%q no TypeSet element %q, with nested attrs %#v in state: %#v", name, attr, values, is.Attributes)
	}
}

func testCheckTypeSetElemNestedAttrsInState(is *terraform.InstanceState, attrParts []string, matchCount int, values interface{}) bool {
	matches := make(map[string]int)

	for stateKey, stateValue := range is.Attributes {
		stateKeyParts := strings.Split(stateKey, ".")
		// a Set/List item with nested attrs would have a flatmap address of
		// at least length 3
		// foo.0.name = "bar"
		if len(stateKeyParts) < 3 || len(attrParts) > len(stateKeyParts) {
			continue
		}
		var pathMatch bool
		for i := range attrParts {
			if attrParts[i] != stateKeyParts[i] && attrParts[i] != sentinelIndex {
				break
			}
			if i == len(attrParts)-1 {
				pathMatch = true
			}
		}
		if !pathMatch {
			continue
		}
		id := stateKeyParts[len(attrParts)-1]
		nestedAttr := strings.Join(stateKeyParts[len(attrParts):], ".")

		var match bool
		switch t := values.(type) {
		case map[string]string:
			if v, keyExists := t[nestedAttr]; keyExists && v == stateValue {
				match = true
			}
		case map[string]*regexp.Regexp:
			if v, keyExists := t[nestedAttr]; keyExists && v != nil && v.MatchString(stateValue) {
				match = true
			}
		}
		if match {
			matches[id] = matches[id] + 1
			if matches[id] == matchCount {
				return true
			}
		}
	}
	return false
}
