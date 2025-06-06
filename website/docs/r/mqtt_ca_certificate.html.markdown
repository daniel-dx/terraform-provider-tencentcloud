---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_ca_certificate"
sidebar_current: "docs-tencentcloud-resource-mqtt_ca_certificate"
description: |-
  Provides a resource to create a MQTT CA certificate
---

# tencentcloud_mqtt_ca_certificate

Provides a resource to create a MQTT CA certificate

## Example Usage

```hcl
resource "tencentcloud_mqtt_ca_certificate" "example" {
  instance_id              = "mqtt-zxjwkr98"
  ca_certificate           = <<-EOF
-----BEGIN CERTIFICATE-----
MIIDUDCCAjigAwIBAgIBATANBgkqhkiG9w0BAQsFADA/MQswCQYDVQQGEwJDTjEb
MBkGA1UEChMSTXkgQ0EgT3JnYW5pemF0aW9uMRMwEQYDVQQDEwpNeSBDQSBSb290
MB4XDTI1MDMyNTEyMzMyM1oXDTM1MDMyNTEyMzMyM1owPzELMAkGA1UEBhMCQ04x
GzAZBgNVBAoTEk15IENBIE9yZ2FuaXphdGlvbjETMBEGA1UEAxMKTXkgQ0EgUm9v
dDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMd8mTXv6SPc6+sQY/Po
zeDHMMIgNFq5wTA23nMPZ15P5PH3Hy76oeR9lPIepQXs36BjoRRDlmc0wb9zhZdt
vG9t3Tr8SzTkLC9nSjjs+TIk26/rAuP1igc+V8MbnSuDkgmhepwioXeMrn/ns7RK
mgvKm5C8tC4MlRlmn0R29EPfchvhW+Ab+mybKFSJfiPABDxDzSfPTCZH2wVTgAIF
0lG93SqrytBJzqhwyXN6bXq/52+CGfG264/fLN4vH+VEGE++ys0eZh+9+0GQ4cFp
gqeRFRYG31ChXMWcnKTLzh/o7GpdTCN31w7h1XkJTbaHNvZbuV0H/wwCVN8bsGkK
zo0CAwEAAaNXMFUwDgYDVR0PAQH/BAQDAgEGMBMGA1UdJQQMMAoGCCsGAQUFBwMB
MA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFCFYj1RkM/mf/ZIYzZPoMvusMM9Z
MA0GCSqGSIb3DQEBCwUAA4IBAQC9TRuaXBnx7OHbdAgukWr4/tbIEhVudKrjEjyV
4tYXhZB/adouWLih/2t+E5U8DdpenTDXhQmT57VVdUjE7ey3VTK9qYQ6swTrOi4A
pW9xwpJuqqQPEj9l8/iOdhiSF3XG/UcxwyctPux6Wmm+Xg0Nz3MV0FGGIi56JZlB
sEE4WHwkzoFYTJxIlBNQvcNxVjWmBUWRm7bBCu7vW3sqdb22Uh5X2E0v5sH0vskG
Bj/1ZqTpbCuNC2UIyiqMGwKjVUifKpEmjzJI/gdGq7c2/o987TYlpWMBE1J1my0l
CjJmbR+Ces1k4hZUWrHijCmLS+iWPiadoQ9xzWgaQeQIbU43
-----END CERTIFICATE-----
EOF
  verification_certificate = <<-EOF
-----BEGIN CERTIFICATE-----
MIIDhDCCAmygAwIBAgIRAOr5LwhpwBWsYLWLt5+HwqMwDQYJKoZIhvcNAQELBQAw
PzELMAkGA1UEBhMCQ04xGzAZBgNVBAoTEk15IENBIE9yZ2FuaXphdGlvbjETMBEG
A1UEAxMKTXkgQ0EgUm9vdDAeFw0yNTAzMjUxMjMzMjNaFw0yNjAzMjUxMjMzMjNa
MFoxCzAJBgNVBAYTAkNOMRwwGgYDVQQKExNDbGllbnQgT3JnYW5pemF0aW9uMS0w
KwYDVQQDEyQ0NzAzY2VmNC0wNDM4LTRkM2QtOTgwNy0zMTAwNjI0ODJkNWIwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDPOe0pjyLtj7Prt23T1Ep3pJaN
aKKiZyikPae4mHXDZ+AQ4zONA78OaJ6S8E9fAmf5tygWWFJQBwccjitlb1nKChGx
KogqeXS/6RwKPNWjaWUp6gabAmcCTF3g6F8gxjJv0eBn0i+UcS2LNp6wRECM08MI
xHc6B/jC78gkp+b4DuNXFQzGeqDTHgneF5immpjLP7ggWTFgjUOJgLAwGRcZf89K
T3TxN1tKtiKxiXIfzAcqgeAfDWjPYgb/3PEcWZj7Zyl4mJAPnAGghFkLDDeJHh6L
Gk9OpfIuAWmlitZFCujnM1MiEGyw5p6kSwAAD7I/p0yyIV/1VYs2AMmu2uehAgMB
AAGjYDBeMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYB
BQUHAwIwDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBQhWI9UZDP5n/2SGM2T6DL7
rDDPWTANBgkqhkiG9w0BAQsFAAOCAQEAcsW08LGC/uARyX2X0QZ9A7I+aIluI23b
VSbNbU1+3SVbm8Jfk63rb/Zkc98jPLds598YswY2gQtjT4+Dcpv60wS+c0Ltw1nJ
O23cp2kJ05+jh/5GywOur7gOG8L1xwUngqX84ObBIyeYv5MfANLmzqfZBs9nKokF
keeHU9Y0NYmFiPw4xNM7S55dbFxKizYd66uGc5b+cWkqg5xNlOqU9He0cBC6KYAj
GnyAz0ruWFPFMlftw5/OwNbc9X8G9wm6+T+bNikzOh3FwMCqjK6hdjtR4/HHxgr2
IpWw7p6yvzDYy4D99PLDjRWP+iLNQCzXOk+PKV+MYCwwYeD1loldOA==
-----END CERTIFICATE-----
EOF
  format                   = "PEM"
  status                   = "ACTIVE"
}
```

## Argument Reference

The following arguments are supported:

* `ca_certificate` - (Required, String, ForceNew) CA certificate.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `verification_certificate` - (Required, String, ForceNew) Verification certificate.
* `format` - (Optional, String, ForceNew) Certificate format, Default is PEM.
* `status` - (Optional, String) Certificate status, Default is ACTIVE.
  ACTIVE activation;
  INACTIVE not active.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `ca_cn` - Certificate common name.
* `ca_sn` - Certificate serial number.
* `created_time` - Certificate create time.
* `not_after_time` - Certificate expiring date.
* `not_before_time` - Certificate effective start date.
* `update_time` - Certificate update time.


## Import

MQTT CA certificate can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_ca_certificate.example mqtt-zxjwkr98#01
```

