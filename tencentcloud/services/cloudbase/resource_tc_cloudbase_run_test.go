package cloudbase_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcbr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcbr/v20220217"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// 测试名称的校验规则
func TestAccTencentCloudCloudbaseRun_nameValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				// 测试名称长度超过限制的情况
				Config: `
			resource "tencentcloud_cloudbase_run" "test" {
				server_name = "this-is-a-very-long-name-that-exceeds-the-limit-of-39-characters-test"
			}
			`,
				ExpectError: regexp.MustCompile("Name length must be between 1 and 39 characters"),
			},
			{
				// 测试名称以大写字母开头的情况
				// 期望返回错误：名称必须以小写字母开头
				Config: `
			resource "tencentcloud_cloudbase_run" "test" {
			    server_name = "Test-cloudbase"
			}
			`,
				ExpectError: regexp.MustCompile("Name must start with a lowercase letter"),
			},
			{
				// 测试名称包含下划线的情况
				// 期望返回错误：名称只能包含小写字母
				Config: `
			resource "tencentcloud_cloudbase_run" "test" {
			    server_name = "test_cloudbase"
			}
			`,
				ExpectError: regexp.MustCompile("Name can only contain lowercase letters"),
			},
			{
				// 测试合法的名称格式
				// 使用小写字母和连字符，符合命名规范
				// 期望测试通过，没有错误
				Config: `
			resource "tencentcloud_cloudbase_run" "test" {
				server_name = "test-cloudbase"
			}
			`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ExpectError:        nil,
			},
		},
	})
}

// 测试值校验
func TestAccTencentCloudCloudbaseRun_fieldValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				// 测试无效的 deploy_type 值
				Config: `
				resource "tencentcloud_cloudbase_run" "test" {
					server_name = "test-invalid"
					deploy_type = "invalid"
				}
				`,
				ExpectError: regexp.MustCompile(`expected deploy_type to be one of \[package image\]`),
			},
			{
				// 测试无效的 open_access_types 值
				Config: `
				resource "tencentcloud_cloudbase_run" "test" {
					server_name = "test-access"
					open_access_types = ["INVALID"]
				}
				`,
				ExpectError: regexp.MustCompile(`expected open_access_types.0 to be one of \[PUBLIC OA MINIAPP VPC\]`),
			},
			{
				// 测试 policy_details 中无效的 policy_type
				Config: `
				resource "tencentcloud_cloudbase_run" "test" {
					server_name = "test-policy"
					policy_details {
						policy_type = "invalid"
					}
				}
				`,
				ExpectError: regexp.MustCompile(`expected policy_details.0.policy_type to be one of \[cpu mem\]`),
			},
		},
	})
}

// 测试 package 部署类型的 CRUD 操作
func TestAccTencentCloudCloudbaseRun_crud_package(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCloudCloudbaseRunDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: `
				resource "tencentcloud_cloudbase_run" "test" {
					server_name = "test-package"
					build_dir = "testdata"
					source_code_include = ["testdata/**/*"]
				}
				`,
				ExpectError: nil,
				Check: resource.ComposeTestCheckFunc(
					// 验证资源ID
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_cloudbase_run.test"),
					// 验证资源属性
					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "server_name", "test-package"),
					// 添加延迟等待
					func(s *terraform.State) error {
						time.Sleep(90 * time.Second)
						fmt.Printf("[%s] 资源创建检查结束...\n", time.Now().Format("2006-01-02 15:04:05"))
						return nil
					},
				),
			},
			// update
			{
				Config: `
				resource "tencentcloud_cloudbase_run" "test" {
					server_name = "test-package"
					build_dir = "testdata"
					source_code_include = ["testdata/**/*"]
					cpu = 0.25
					mem = 0.5
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_cloudbase_run.test"),

					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "server_name", "test-package"),
					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "cpu", "0.25"),
					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "mem", "0.5"),
					// 添加延迟等待
					func(s *terraform.State) error {
						time.Sleep(120 * time.Second)
						fmt.Printf("[%s] 资源更新检查结束...\n", time.Now().Format("2006-01-02 15:04:05"))
						return nil
					},
				),
			},
		},
	})
}

// 测试 image 部署类型的 CRUD 操作
func TestAccTencentCloudCloudbaseRun_crud_image(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCloudCloudbaseRunDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: `
                resource "tencentcloud_cloudbase_run" "test" {
                    server_name = "test-image"
                    deploy_type = "image"
                    image_url = "registry.cn-beijing.aliyuncs.com/ijx-public/nginx:alpine"
					env_params = {
						"TEST_KEY" = "test_value"
					}
                }
                `,
				ExpectError: nil,
				Check: resource.ComposeTestCheckFunc(
					// 验证资源ID
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_cloudbase_run.test"),
					// 验证资源属性
					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "server_name", "test-image"),
					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "deploy_type", "image"),
					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "image_url", "registry.cn-beijing.aliyuncs.com/ijx-public/nginx:alpine"),
					// / 验证环境变量
					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "env_params.TEST_KEY", "test_value"),
					// 添加延迟等待
					func(s *terraform.State) error {
						time.Sleep(90 * time.Second)
						fmt.Printf("[%s] 资源创建检查结束...\n", time.Now().Format("2006-01-02 15:04:05"))
						return nil
					},
				),
			},
			// update
			{
				Config: `
                resource "tencentcloud_cloudbase_run" "test" {
                    server_name = "test-image"
					deploy_type = "image"
                    image_url = "registry.cn-beijing.aliyuncs.com/ijx-public/nginx:alpine"
					env_params = {
						"TEST_KEY" = "test_value"
					}
                    cpu = 0.25
                    mem = 0.5
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_cloudbase_run.test"),
					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "server_name", "test-image"),
					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "cpu", "0.25"),
					resource.TestCheckResourceAttr("tencentcloud_cloudbase_run.test", "mem", "0.5"),
					// 添加延迟等待
					func(s *terraform.State) error {
						time.Sleep(120 * time.Second)
						fmt.Printf("[%s] 资源更新检查结束...\n", time.Now().Format("2006-01-02 15:04:05"))
						return nil
					},
				),
			},
		},
	})
}

func testAccCheckCloudCloudbaseRunDestroy(s *terraform.State) error {

	fmt.Printf("[%s] 开始检查资源是否已被销毁...\n", time.Now().Format("2006-01-02 15:04:05"))

	// 获取客户端连接
	client := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()

	// 从环境变量中获取 env_id
	envId := os.Getenv("ENV_ID")
	if envId == "" {
		return fmt.Errorf("environment variable ENV_ID is not set")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cloudbase_run" {
			continue
		}

		// 添加重试逻辑
		retryCount := 0
		maxRetries := 10                  // 最多重试10次
		retryInterval := time.Second * 10 // 每次重试间隔10秒

		for retryCount < maxRetries {
			request := tcbr.NewDescribeCloudRunServerDetailRequest()
			request.ServerName = helper.String(rs.Primary.ID)
			request.EnvId = helper.String(envId)
			response, err := client.UseTcbrClient().DescribeCloudRunServerDetail(request)

			if err != nil {
				// 如果是资源不存在的错误，说明资源已被成功删除
				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == "ResourceNotFound" {
						return nil
					}
				}
				return err
			}

			// 如果资源仍然存在，等待一段时间后重试
			if response != nil && response.Response != nil && response.Response.ServerConfig != nil {
				retryCount++
				if retryCount < maxRetries {
					time.Sleep(retryInterval)
					continue
				}
				return fmt.Errorf("CloudBase Run Server %s still exists after %d retries", rs.Primary.ID, maxRetries)
			}

			// 如果资源不存在，退出循环
			return nil
		}
	}

	return nil
}
