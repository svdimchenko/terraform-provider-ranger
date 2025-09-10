resource "ranger_policy" "example" {
  service    = "service-name"
  name       = "test-policy-name"
  definition = file("${path.module}/policy.json")
}
