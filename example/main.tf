provider "sumologic" {}

resource "sumologic_hosted_collector" "bsick-test" {
  name        = "bsick-test"
  description = "Test collector for verifying terraform provider"
}

resource "sumologic_http_source" "bsick-test" {
  collector_id        = "${sumologic_hosted_collector.bsick-test.id}"
  name                = "bsick-test"
  message_per_request = false
}
