Provide a resource to create a Cloudbase Run service.

Example Usage

Using Package Deployment

```hcl
resource "tencentcloud_cloudbase_run" "example" {
  server_name = "example"
  dockerfile  = "Dockerfile"
}
```

Using Image Deployment

```hcl
resource "tencentcloud_cloudbase_run" "example" {
  server_name = "example"
  deploy_type = "image"
  image_url   = "registry.example.com/my-image:latest"
}
```

Import

Cloudbase Run service can be imported, e.g.

```sh
$ terraform import tencentcloud_cloudbase_run.example example
```