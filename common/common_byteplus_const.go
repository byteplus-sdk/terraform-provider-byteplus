package common

type RequestConvertType int

const (
	ConvertDefault RequestConvertType = iota
	ConvertWithN
	ConvertListUnique
	ConvertListN
	ConvertSingleN
	ConvertJsonObject
	ConvertJsonArray
	ConvertJsonObjectArray
)

type RequestConvertMode int

const (
	RequestConvertAll RequestConvertMode = iota
	RequestConvertInConvert
	RequestConvertIgnore
)

type RequestContentType int

const (
	ContentTypeDefault RequestContentType = iota
	ContentTypeJson
)

type ServiceCategory int

const (
	DefaultServiceCategory ServiceCategory = iota
	ServiceBypass
)

type SpecialParamType int

const (
	DomainParam SpecialParamType = iota
	HeaderParam
	PathParam
	UrlParam
	FilePathParam
)

const (
	RegionalService = "Regional"
	GlobalService   = "Global"

	ByteplusIpv4EndpointSuffix = "byteplusapi.com"
)
