provider "ranger" {
  url             = "https://domain.com"
  username        = "admin"
  password        = "password"
  skip_tls_verify = true  # for dev purpose only
}
