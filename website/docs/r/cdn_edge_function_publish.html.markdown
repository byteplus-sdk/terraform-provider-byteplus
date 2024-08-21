---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_edge_function_publish"
sidebar_current: "docs-byteplus-resource-cdn_edge_function_publish"
description: |-
  Provides a resource to manage cdn edge function publish
---
# byteplus_cdn_edge_function_publish
Provides a resource to manage cdn edge function publish
## Example Usage
```hcl
resource "byteplus_cdn_edge_function" "foo" {
  name         = "acc-test-function"
  remark       = "tf-test"
  project_name = "default"
  source_code  = base64encode("hello world")
  envs {
    key   = "k1"
    value = "v1"
  }
  canary_countries = ["China", "Japan", "United Kingdom"]
}

resource "byteplus_cdn_edge_function_publish" "foo" {
  function_id    = byteplus_cdn_edge_function.foo.id
  description    = "test publish"
  publish_action = "FullPublish"
}
```
## Argument Reference
The following arguments are supported:
* `function_id` - (Required, ForceNew) The ID of the function to which you want publish.
* `publish_action` - (Required, ForceNew) The publish action of the edge function. Valid values: `FullPublish`, `CanaryPublish`, `SnapshotPublish`.
When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.
* `description` - (Optional, ForceNew) The description for this release.
* `publish_type` - (Optional, ForceNew) The publish type of SnapshotPublish: 
200: FullPublish
100: CanaryPublish. This field is required and valid when the `publish_action` is `SnapshotPublish`.
When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.
* `version_tag` - (Optional, ForceNew) The specified version number to be published. This field is required and valid when the `publish_action` is `SnapshotPublish`.
 When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `content` - The content of the release record.
* `create_time` - The create time of the release record. Displayed in UNIX timestamp format.
* `creator` - The creator of the release record.
* `update_time` - The update time of the release record. Displayed in UNIX timestamp format.


## Import
CdnEdgeFunctionPublish can be imported using the function_id:ticket_id, e.g.
```
$ terraform import byteplus_cdn_edge_function_publish.default function_id:ticket_id
```

