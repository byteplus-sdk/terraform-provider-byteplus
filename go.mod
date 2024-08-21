module github.com/byteplus-sdk/terraform-provider-byteplus

go 1.12

require (
	github.com/byteplus-sdk/byteplus-go-sdk v0.0.0
	github.com/fatih/color v1.7.0
	github.com/google/uuid v1.3.0
	github.com/hashicorp/hcl/v2 v2.0.0
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/mitchellh/copystructure v1.0.0
	github.com/stretchr/testify v1.8.2
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
)

replace github.com/byteplus-sdk/byteplus-go-sdk v0.0.0 => ./byteplus-go-sdk
