---
subcategory: "CloudBase"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cloudbase_run"
sidebar_current: "docs-tencentcloud-resource-cloudbase_run"
description: |-
  Provide a resource to create a Cloudbase Run service.
---

# tencentcloud_cloudbase_run

Provide a resource to create a Cloudbase Run service.

## Example Usage

### Using Package Deployment

```hcl
resource "tencentcloud_cloudbase_run" "example" {
  server_name = "example"
  dockerfile  = "Dockerfile"
}
```

### Using Image Deployment

```hcl
resource "tencentcloud_cloudbase_run" "example" {
  server_name = "example"
  deploy_type = "image"
  image_url   = "registry.example.com/my-image:latest"
}
```

## Argument Reference

The following arguments are supported:

* `server_name` - (Required, String, ForceNew) Name of the Cloudbase run. Only lowercase letters, numbers and hyphens(-) are allowed, must start with a lowercase letter, and cannot exceed 39 characters.
* `build_dir` - (Optional, String) Build directory. Only valid when deploy_type is package.
* `cpu` - (Optional, Float64) CPU specification of the Cloudbase run service. Unit: cores.
* `custom_logs` - (Optional, String) Log collection path.
* `deploy_type` - (Optional, String) Deployment type. Valid values: package, image. Default value: package.
* `dockerfile` - (Optional, String) Dockerfile name. Default value: Dockerfile. Only valid when deploy_type is package.
* `env_params` - (Optional, Map) Environment variables in key-value pairs.
* `image_url` - (Optional, String) Image URL. Required when deploy_type is image.
* `max_num` - (Optional, Int) Maximum number of replicas.
* `mem` - (Optional, Float64) Memory specification of the Cloudbase run service. Unit: GB.
* `min_num` - (Optional, Int) Minimum number of replicas.
* `open_access_types` - (Optional, Set: [`String`]) Access types. Valid values: "PUBLIC", "OA", "MINIAPP", "VPC". Default value: ["PUBLIC", "OA"].
* `policy_details` - (Optional, List) HPA policy details.
* `port` - (Optional, Int) Service port.
* `source_code_exclude` - (Optional, List: [`String`]) The list of files or directories to exclude. Only valid when deploy_type is package. Default value: [".git", ".terraform*", "logs", "node_modules", "__pycache__", "venv", "terraform.tfstate*"].
* `source_code_hash` - (Optional, String) Hash value used to track source code updates. Only valid when deploy_type is package.
* `source_code_include` - (Optional, List: [`String`]) The list of files or directories to include. Only valid when deploy_type is package.

The `policy_details` object supports the following:

* `policy_threshold` - (Optional, Int) Policy threshold.
* `policy_type` - (Optional, String) Policy type. Valid values: cpu, mem.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Cloudbase Run service can be imported, e.g.

```sh
$ terraform import tencentcloud_cloudbase_run.example example
```

