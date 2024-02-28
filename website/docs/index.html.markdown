---
layout: "byteplus"
page_title: "Provider: byteplus"
sidebar_current: "docs-byteplus-index"
description: |-
The byteplus provider is used to interact with many resources supported by Byteplus. The provider needs to be configured with the proper credentials before it can be used.
---

# Byteplus Provider

The Byteplus provider is used to interact with many resources supported by [Byteplus](https://www.byteplus.com/en).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** This guide requires an available Byteplus account or sub-account with project to create resources.

## Example Usage
```hcl
# Configure the Byteplus Provider
provider "byteplus" {
  access_key = "your ak"
  secret_key = "your sk"
  session_token = "sts token"
  region = "cn-beijing"
}

# Query Vpc
data "byteplus_vpcs" "default"{
  ids = ["vpc-mizl7m1kqccg5smt1bdpijuj"]
}

#Create vpc
resource "byteplus_vpc" "foo" {
  vpc_name = "tf-test-1"
  cidr_block = "172.16.0.0/16"
  dns_servers = ["8.8.8.8","114.114.114.114"]
}

```

## Authentication

The Byteplus provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables

### Static credentials

Static credentials can be provided by adding an `public_key` and `private_key` in-line in the
byteplus provider block:

Usage:

```hcl
provider "byteplus" {
   access_key = "your ak"
   secret_key = "your sk"
   region = "cn-beijing"
}
```

### Environment variables

You can provide your credentials via `BYTEPLUS_ACCESS_KEY` and `BYTEPLUS_SECRET_KEY`
environment variables, representing your byteplus public key and private key respectively.
`BYTEPLUS_REGION` is also used, if applicable:

```hcl
provider "byteplus" {
  
}
```

Usage:

```hcl
$ export BYTEPLUS_ACCESS_KEY="your_public_key"
$ export BYTEPLUS_SECRET_KEY="your_private_key"
$ export BYTEPLUS_REGION="cn-beijing"
$ terraform plan
```

