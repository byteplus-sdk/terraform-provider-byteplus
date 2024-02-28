package common

type SdkClient struct {
	Region          string
	UniversalClient *Universal
	BypassSvcClient *BypassSvc
}
