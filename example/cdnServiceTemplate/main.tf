resource "byteplus_cdn_service_template" "foo" {
  title = "tf-test2"
  message = "test2"
  project = "test"
  service_template_config = jsonencode(
    {
      OriginIpv6 = "followclient"
      ConditionalOrigin = {
        OriginRules = []
      }
      Origin = [{
        OriginAction = {
          OriginLines = [
            {
              Address = "10.10.10.10"
              HttpPort = "80"
              HttpsPort = "443"
              InstanceType = "ip"
              OriginType = "primary"
              Weight = "1"
            }
          ]
        }
      }]
      OriginHost = ""
      OriginProtocol = "http"
      OriginHost = ""
    }
  )
}