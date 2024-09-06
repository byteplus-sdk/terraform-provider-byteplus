---
subcategory: "VPN"
layout: "byteplus"
page_title: "Byteplus: byteplus_ssl_vpn_client_certs"
sidebar_current: "docs-byteplus-datasource-ssl_vpn_client_certs"
description: |-
  Use this data source to query detailed information of ssl vpn client certs
---
# byteplus_ssl_vpn_client_certs
Use this data source to query detailed information of ssl vpn client certs
## Example Usage
```hcl
data "byteplus_zones" "foo" {
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.byteplus_zones.foo.zones[0].id
  vpc_id      = byteplus_vpc.foo.id
}

resource "byteplus_vpn_gateway" "foo" {
  vpc_id              = byteplus_vpc.foo.id
  subnet_id           = byteplus_subnet.foo.id
  bandwidth           = 5
  vpn_gateway_name    = "acc-test1"
  description         = "acc-test1"
  period              = 7
  project_name        = "default"
  ssl_enabled         = true
  ssl_max_connections = 5
}

resource "byteplus_ssl_vpn_server" "foo" {
  vpn_gateway_id      = byteplus_vpn_gateway.foo.id
  local_subnets       = [byteplus_subnet.foo.cidr_block]
  client_ip_pool      = "172.16.2.0/24"
  ssl_vpn_server_name = "acc-test-ssl"
  description         = "acc-test"
  protocol            = "UDP"
  cipher              = "AES-128-CBC"
  auth                = "SHA1"
  compress            = true
}

resource "byteplus_ssl_vpn_client_cert" "foo" {
  ssl_vpn_server_id        = byteplus_ssl_vpn_server.foo.id
  ssl_vpn_client_cert_name = "acc-test-client-cert-${count.index}"
  description              = "acc-test"
  count                    = 5
}

data "byteplus_ssl_vpn_client_certs" "foo" {
  ids = byteplus_ssl_vpn_client_cert.foo[*].id
}
```
## Argument Reference
The following arguments are supported:
* `ids` - (Optional) The ids list of ssl vpn client cert.
* `name_regex` - (Optional) A Name Regex of ssl vpn client cert.
* `output_file` - (Optional) File name where to save data source results.
* `ssl_vpn_client_cert_name` - (Optional) The name of the ssl vpn client cert.
* `ssl_vpn_server_id` - (Optional) The id of the ssl vpn server.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `ssl_vpn_client_certs` - The collection of of ssl vpn client certs.
    * `ca_certificate` - The CA certificate.
    * `certificate_status` - The status of the ssl vpn client cert.
    * `client_certificate` - The client certificate.
    * `client_key` - The key of the ssl vpn client.
    * `creation_time` - The creation time of the ssl vpn client cert.
    * `description` - The description of the ssl vpn client cert.
    * `expired_time` - The expired time of the ssl vpn client cert.
    * `id` - The id of the ssl vpn client cert.
    * `open_vpn_client_config` - The config of the open vpn client.
    * `ssl_vpn_client_cert_id` - The id of the ssl vpn client cert.
    * `ssl_vpn_client_cert_name` - The name of the ssl vpn client cert.
    * `ssl_vpn_server_id` - The id of the ssl vpn server.
    * `status` - The status of the ssl vpn client.
    * `update_time` - The update time of the ssl vpn client cert.
* `total_count` - The total count of ssl vpn client cert query.


