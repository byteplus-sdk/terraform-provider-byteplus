package iam_role_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_role"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusIamRolesDatasourceConfig = `
resource "byteplus_iam_role" "foo1" {
	role_name = "acc-test-role1"
    display_name = "acc-test1"
	description = "acc-test1"
    trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"
	max_session_duration = 3600
}

resource "byteplus_iam_role" "foo2" {
    role_name = "acc-test-role2"
    display_name = "acc-test2"
	description = "acc-test2"
    trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"ecs\"]}}]}"
	max_session_duration = 3600
}

data "byteplus_iam_roles" "foo"{
    role_name = "${byteplus_iam_role.foo1.role_name},${byteplus_iam_role.foo2.role_name}"
}
`

func TestAccByteplusIamRolesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_iam_roles.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_role.ByteplusIamRoleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamRolesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "roles.#", "2"),
				),
			},
		},
	})
}
