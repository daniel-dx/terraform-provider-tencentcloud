package cloudbase

import (
	"context"
	"fmt"
	"unicode"

	"log"
	"os"
	"strings"

	"crypto/sha256"
	"encoding/json"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcbr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcbr/v20220217"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// uploadResult 定义在文件开头的结构体部分
type uploadResult struct {
	packageName    string
	packageVersion string
}

// 在文件开头的常量定义部分添加
var defaultSourceCodeExclude = []string{".git", ".terraform*", "logs", "node_modules", "__pycache__", "venv", "terraform.tfstate*"}

func ResourceTencentCloudCloudbaseRun() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTencentCloudCloudbaseRunCreate,
		ReadContext:   resourceTencentCloudCloudbaseRunRead,
		UpdateContext: resourceTencentCloudCloudbaseRunUpdate,
		DeleteContext: resourceTencentCloudCloudbaseRunDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: resourceTencentCloudCloudbaseRunCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"server_name": {
				Type:             schema.TypeString,
				Description:      "Name of the Cloudbase run. Only lowercase letters, numbers and hyphens(-) are allowed, must start with a lowercase letter, and cannot exceed 39 characters.",
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateCloudbaseRunName(),
			},

			// 代码更新相关配置
			"source_code_hash": {
				Type:        schema.TypeString,
				Description: "Hash value used to track source code updates. Only valid when deploy_type is package.",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("deploy_type").(string) != "package"
				},
			},
			"source_code_include": {
				Type:        schema.TypeList,
				Description: "The list of files or directories to include. Only valid when deploy_type is package.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("deploy_type").(string) != "package"
				},
			},
			"source_code_exclude": {
				Type:        schema.TypeList,
				Description: fmt.Sprintf("The list of files or directories to exclude. Only valid when deploy_type is package. Default value: [\"%s\"].", strings.Join(defaultSourceCodeExclude, "\", \"")),
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("deploy_type").(string) != "package"
				},
			},

			// 部署相关配置
			"deploy_type": {
				Type:        schema.TypeString,
				Description: "Deployment type. Valid values: package, image. Default value: package.",
				Optional:    true,
				Default:     "package",
				ValidateFunc: validation.StringInSlice([]string{
					"package",
					"image",
				}, false),
			},
			"image_url": {
				Type:        schema.TypeString,
				Description: "Image URL. Required when deploy_type is image.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("deploy_type").(string) != "image"
				},
			},
			"dockerfile": {
				Type:        schema.TypeString,
				Description: "Dockerfile name. Default value: Dockerfile. Only valid when deploy_type is package.",
				Optional:    true,
				Default:     "Dockerfile",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("deploy_type").(string) != "package"
				},
			},
			"build_dir": {
				Type:        schema.TypeString,
				Description: "Build directory. Only valid when deploy_type is package.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("deploy_type").(string) != "package"
				},
			},

			// 服务相关配置
			"port": {
				Type:        schema.TypeInt,
				Description: "Service port.",
				Optional:    true,
				Computed:    true,
			},
			"cpu": {
				Type:        schema.TypeFloat,
				Description: "CPU specification of the Cloudbase run service. Unit: cores.",
				Optional:    true,
				Computed:    true,
			},
			"mem": {
				Type:        schema.TypeFloat,
				Description: "Memory specification of the Cloudbase run service. Unit: GB.",
				Optional:    true,
				Computed:    true,
			},
			"min_num": {
				Type:        schema.TypeInt,
				Description: "Minimum number of replicas.",
				Optional:    true,
				Computed:    true,
			},
			"max_num": {
				Type:        schema.TypeInt,
				Description: "Maximum number of replicas.",
				Optional:    true,
				Computed:    true,
			},
			"custom_logs": {
				Type:        schema.TypeString,
				Description: "Log collection path.",
				Optional:    true,
				Computed:    true,
			},
			"open_access_types": {
				Type:        schema.TypeSet,
				Description: `Access types. Valid values: "PUBLIC", "OA", "MINIAPP", "VPC". Default value: ["PUBLIC", "OA"].`,
				Optional:    true,
				Computed:    true, // 用户没设置，则以服务端配置为准
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"PUBLIC",
						"OA",
						"MINIAPP",
						"VPC",
					}, false),
				},
			},
			"policy_details": {
				Type:        schema.TypeList,
				Description: "HPA policy details.",
				Optional:    true,
				Computed:    true, // 用户没设置，则以服务端配置为准
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_type": {
							Type:        schema.TypeString,
							Description: "Policy type. Valid values: cpu, mem.",
							Optional:    true,
							ValidateFunc: validation.StringInSlice([]string{
								"cpu",
								"mem",
							}, false),
						},
						"policy_threshold": {
							Type:        schema.TypeInt,
							Description: "Policy threshold.",
							Optional:    true,
						},
					},
				},
			},
			"env_params": {
				Type:        schema.TypeMap,
				Description: "Environment variables in key-value pairs.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudCloudbaseRunCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return createOrUpdateCloudRunServer(ctx, d, meta, true)
}

func resourceTencentCloudCloudbaseRunUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return createOrUpdateCloudRunServer(ctx, d, meta, false)
}

func resourceTencentCloudCloudbaseRunRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// 获取服务客户端
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := TcbrService{client: client}

	// 从环境变量中获取 env_id
	envId := os.Getenv("ENV_ID")
	if envId == "" {
		return diag.FromErr(fmt.Errorf("environment variable ENV_ID is not set"))
	}

	// 使用服务名称作为资源 ID
	serverName := d.Id()

	// 查询服务详情
	serverInfo, err := service.DescribeCloudRunServerDetail(ctx, serverName, envId)
	if err != nil {
		return diag.FromErr(err)
	}

	if serverInfo == nil || serverInfo.Response == nil || serverInfo.Response.ServerConfig == nil {
		d.SetId("")
		return nil
	}

	// 更新资源数据
	// deploy_type / image_url / build_dir 以本地为准备
	detail := serverInfo.Response.ServerConfig

	if err := d.Set("server_name", detail.ServerName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dockerfile", detail.Dockerfile); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port", detail.Port); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("open_access_types", detail.OpenAccessTypes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu", detail.Cpu); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mem", detail.Mem); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("min_num", detail.MinNum); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("max_num", detail.MaxNum); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("custom_logs", detail.CustomLogs); err != nil {
		return diag.FromErr(err)
	}
	if detail.EnvParams != nil && *detail.EnvParams != "" {
		var envParams map[string]interface{}
		if err := json.Unmarshal([]byte(*detail.EnvParams), &envParams); err != nil {
			return diag.FromErr(fmt.Errorf("failed to unmarshal env_params: %v", err))
		}
		if err := d.Set("env_params", envParams); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("env_params", nil); err != nil {
			return diag.FromErr(err)
		}
	}
	if detail.PolicyDetails != nil && len(detail.PolicyDetails) > 0 {
		policyDetails := make([]map[string]interface{}, 0, len(detail.PolicyDetails))
		for _, policy := range detail.PolicyDetails {
			policyDetail := map[string]interface{}{
				"policy_type":      policy.PolicyType,
				"policy_threshold": policy.PolicyThreshold,
			}
			policyDetails = append(policyDetails, policyDetail)
		}
		if err := d.Set("policy_details", policyDetails); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceTencentCloudCloudbaseRunDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// 获取服务客户端
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := TcbrService{client: client}

	// 从环境变量中获取 env_id
	envId := os.Getenv("ENV_ID")
	if envId == "" {
		return diag.FromErr(fmt.Errorf("environment variable ENV_ID is not set"))
	}

	// 使用服务名称作为资源 ID
	serverName := d.Id()

	// 调用删除服务
	_, err := service.DeleteCloudRunServer(ctx, envId, serverName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("删除服务失败: %v", err))
	}

	return diags
}

// createOrUpdateCloudRunServer 处理创建或更新服务的通用逻辑
func createOrUpdateCloudRunServer(ctx context.Context, d *schema.ResourceData, meta interface{}, isCreate bool) diag.Diagnostics {
	var diags diag.Diagnostics

	// 获取服务客户端
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := TcbService{client: client}
	tcbrService := TcbrService{client: client}

	// 从环境变量中获取 env_id
	envId := os.Getenv("ENV_ID")
	if envId == "" {
		return diag.FromErr(fmt.Errorf("environment variable ENV_ID is not set"))
	}

	var request interface{}
	if isCreate {
		req := tcbr.NewCreateCloudRunServerRequest()
		req.EnvId = helper.String(envId)
		req.ServerName = helper.String(d.Get("server_name").(string))
		request = req
	} else {
		req := tcbr.NewUpdateCloudRunServerRequest()
		req.EnvId = helper.String(envId)
		req.ServerName = helper.String(d.Id())
		request = req
	}

	// 设置 DeployInfo 结构体的参数
	var deployInfo *tcbr.DeployParam
	if isCreate {
		request.(*tcbr.CreateCloudRunServerRequest).DeployInfo = &tcbr.DeployParam{}
		deployInfo = request.(*tcbr.CreateCloudRunServerRequest).DeployInfo
	} else {
		request.(*tcbr.UpdateCloudRunServerRequest).DeployInfo = &tcbr.DeployParam{}
		deployInfo = request.(*tcbr.UpdateCloudRunServerRequest).DeployInfo
	}
	if err := setDeployInfo(ctx, d, &service, deployInfo); err != nil {
		return diag.FromErr(err)
	}

	// 设置 ServerConfig 结构体的参数
	var serverConfig *tcbr.ServerBaseConfig
	if isCreate {
		request.(*tcbr.CreateCloudRunServerRequest).ServerConfig = &tcbr.ServerBaseConfig{}
		serverConfig = request.(*tcbr.CreateCloudRunServerRequest).ServerConfig
	} else {
		request.(*tcbr.UpdateCloudRunServerRequest).ServerConfig = &tcbr.ServerBaseConfig{}
		serverConfig = request.(*tcbr.UpdateCloudRunServerRequest).ServerConfig
	}
	if err := setServerConfig(d, serverConfig); err != nil {
		return diag.FromErr(err)
	}

	// 调用服务
	var err error
	if isCreate {
		_, err = tcbrService.CreateCloudRunServer(ctx, request.(*tcbr.CreateCloudRunServerRequest))
	} else {
		_, err = tcbrService.UpdateCloudRunServer(ctx, request.(*tcbr.UpdateCloudRunServerRequest))
	}
	if err != nil {
		action := "创建"
		if !isCreate {
			action = "更新"
		}
		return diag.FromErr(fmt.Errorf("%s服务失败: %v", action, err))
	}

	// 设置资源 ID (仅创建时需要)
	if isCreate {
		d.SetId(d.Get("server_name").(string))
	}

	return diags
}

// validateCloudbaseRunName 验证Cloudbase Run服务名称的合法性
func validateCloudbaseRunName() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		value := v.(string)
		runes := []rune(value)
		var diags diag.Diagnostics

		// 检查长度
		if len(runes) > 39 || len(runes) < 1 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid name length",
				Detail:   "Name length must be between 1 and 39 characters",
			})
			return diags
		}

		// 检查首字符是否为小写字母
		if !unicode.IsLower(runes[0]) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid name format",
				Detail:   "Name must start with a lowercase letter",
			})
			return diags
		}

		// 检查所有字符
		for _, r := range runes {
			if !unicode.IsLower(r) && !unicode.IsNumber(r) && r != '-' {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid name format",
					Detail:   "Name can only contain lowercase letters, numbers and hyphens(-)",
				})
				return diags
			}
		}
		return diags
	}
}

// prepareAndUploadCode 准备并上传代码
func prepareAndUploadCode(ctx context.Context, d *schema.ResourceData, service *TcbService, name, envId, buildDir string) (*uploadResult, error) {
	// 获取上传信息
	buildService, err := service.DescribeCloudBaseBuildService(ctx, name, envId, nil)
	if err != nil {
		return nil, fmt.Errorf("获取构建服务信息失败: %v", err)
	}

	if buildService.Response == nil || buildService.Response.UploadUrl == nil || buildService.Response.PackageName == nil {
		return nil, fmt.Errorf("构建服务返回信息不完整")
	}

	// 构建 headers
	headers := make(map[string]string)
	if buildService.Response.UploadHeaders != nil {
		for _, header := range buildService.Response.UploadHeaders {
			if header.Key != nil && header.Value != nil {
				headers[*header.Key] = *header.Value
			}
		}
	}

	// 获取排除列表，如果用户没有设置则使用默认值
	var excludeList []string
	if v, ok := d.GetOk("source_code_exclude"); ok {
		excludeList = make([]string, len(v.([]interface{})))
		for i, item := range v.([]interface{}) {
			excludeList[i] = item.(string)
		}
	} else {
		excludeList = defaultSourceCodeExclude
	}

	// 获取包含列表
	var includeList []string
	if v, ok := d.GetOk("source_code_include"); ok {
		includeList = make([]string, len(v.([]interface{})))
		for i, item := range v.([]interface{}) {
			includeList[i] = item.(string)
		}
	}

	// 压缩代码目录
	zipFile, err := zipSourceCode(buildDir, &ZipOptions{
		IncludeList: includeList,
		ExcludeList: excludeList,
	})
	if err != nil {
		return nil, fmt.Errorf("代码压缩失败: %v", err)
	}
	// 打印压缩文件路径
	log.Printf("[DEBUG] Zip file path: %s", zipFile)
	defer os.Remove(zipFile)

	// 上传文件
	uploadOpts := UploadOptions{
		Action:  *buildService.Response.UploadUrl,
		File:    zipFile,
		Method:  "PUT",
		Headers: headers,
	}

	_, err = uploadFile(uploadOpts)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %v", err)
	}

	return &uploadResult{
		packageName:    *buildService.Response.PackageName,
		packageVersion: *buildService.Response.PackageVersion,
	}, nil
}

// setDeployInfo 设置部署信息
func setDeployInfo(ctx context.Context, d *schema.ResourceData, service *TcbService, deployInfo *tcbr.DeployParam) error {
	deployType := d.Get("deploy_type").(string)
	serverName := d.Get("server_name").(string)
	envId := os.Getenv("ENV_ID")
	targetDir, _ := os.Getwd()

	// 上传代码并获取包信息
	var packageInfo *uploadResult
	if deployType == "package" {
		result, err := prepareAndUploadCode(ctx, d, service, serverName, envId, targetDir)
		if err != nil {
			return fmt.Errorf("上传代码失败: %v", err)
		}
		packageInfo = result
	}

	// Update 时需要设置 ReleaseType
	if d.Id() != "" {
		deployInfo.ReleaseType = helper.String("FULL")
	}

	deployInfo.DeployType = helper.String(deployType)
	if deployType == "package" {
		deployInfo.PackageName = &packageInfo.packageName
		deployInfo.PackageVersion = &packageInfo.packageVersion
	} else {
		deployInfo.ImageUrl = helper.String(d.Get("image_url").(string))
	}

	return nil
}

// setServerConfig 设置服务配置
func setServerConfig(d *schema.ResourceData, config *tcbr.ServerBaseConfig) error {
	// 通过检查资源 ID 判断是否为创建操作
	isCreate := d.Id() == ""

	// 设置基础配置
	config.Dockerfile = helper.String(d.Get("dockerfile").(string))
	config.BuildDir = helper.String(d.Get("build_dir").(string))

	// 设置 CPU 配置
	if v, ok := d.GetOk("cpu"); ok {
		config.Cpu = helper.Float64(v.(float64))
	} else if isCreate {
		config.Cpu = helper.Float64(0)
	}

	// 设置内存配置
	if v, ok := d.GetOk("mem"); ok {
		config.Mem = helper.Float64(v.(float64))
	} else if isCreate {
		config.Mem = helper.Float64(0)
	}

	// 设置最小副本数
	if v, ok := d.GetOk("min_num"); ok {
		config.MinNum = helper.Uint64(uint64(v.(int)))
	} else if isCreate {
		config.MinNum = helper.Uint64(0)
	}

	// 设置最大副本数
	if v, ok := d.GetOk("max_num"); ok {
		config.MaxNum = helper.Uint64(uint64(v.(int)))
	} else if isCreate {
		config.MaxNum = helper.Uint64(0)
	}

	// 设置服务端口，默认值为 80
	if v, ok := d.GetOk("port"); ok {
		config.Port = helper.Int64(int64(v.(int)))
	} else if isCreate {
		config.Port = helper.Int64(80)
	}

	// 设置日志采集路径
	if v, ok := d.GetOk("custom_logs"); ok {
		config.CustomLogs = helper.String(v.(string))
	} else if isCreate {
		config.CustomLogs = helper.String("")
	}

	// 设置访问类型
	var accessTypes []interface{}
	if v, ok := d.GetOk("open_access_types"); ok {
		accessTypesSet := v.(*schema.Set)
		accessTypes = accessTypesSet.List()
	} else {
		if isCreate {
			accessTypes = []interface{}{"PUBLIC", "OA"}
		}
	}
	config.OpenAccessTypes = make([]*string, 0, len(accessTypes))
	for _, at := range accessTypes {
		config.OpenAccessTypes = append(config.OpenAccessTypes, helper.String(at.(string)))
	}

	// 设置策略详情
	config.PolicyDetails = []*tcbr.HpaPolicy{}
	if v, ok := d.GetOk("policy_details"); ok {
		policyDetails := v.([]interface{})
		if len(policyDetails) > 0 {
			pd := policyDetails[0].(map[string]interface{})
			policy := &tcbr.HpaPolicy{
				PolicyType:      helper.String(pd["policy_type"].(string)),
				PolicyThreshold: helper.Uint64(uint64(pd["policy_threshold"].(int))),
			}
			config.PolicyDetails = append(config.PolicyDetails, policy)
		}
	}

	// 设置环境变量
	if v, ok := d.GetOk("env_params"); ok {
		envParams := v.(map[string]interface{})
		if len(envParams) > 0 {
			envParamsBytes, err := json.Marshal(envParams)
			if err != nil {
				return fmt.Errorf("序列化环境变量失败: %v", err)
			}
			config.EnvParams = helper.String(string(envParamsBytes))
		}
	}

	// 设置非用户配置的参数
	config.EnvId = helper.String(os.Getenv("ENV_ID"))
	config.ServerName = helper.String(d.Get("server_name").(string))
	config.InitialDelaySeconds = helper.Uint64(0)
	config.HasDockerfile = helper.Bool(true)
	if isCreate {
		config.CreateTime = helper.String("")
		config.Tag = helper.String("")
	}

	return nil
}

// resourceTencentCloudCloudbaseRunCustomizeDiff 在 plan 阶段计算源代码哈希值
func resourceTencentCloudCloudbaseRunCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	deployType := d.Get("deploy_type").(string)

	// 只有在部署类型为 package 时才需要计算源代码哈希
	if deployType != "package" {
		return nil
	}

	// ============ 计算目录哈希值 ============ //
	currentDir, _ := os.Getwd()

	// 获取排除列表，如果用户没有设置则使用默认值
	var excludeList []string
	if v, ok := d.GetOk("source_code_exclude"); ok {
		excludeList = make([]string, len(v.([]interface{})))
		for i, item := range v.([]interface{}) {
			excludeList[i] = item.(string)
		}
	} else {
		excludeList = defaultSourceCodeExclude
	}

	// 获取包含列表
	var includeList []string
	if v, ok := d.GetOk("source_code_include"); ok {
		includeList = make([]string, len(v.([]interface{})))
		for i, item := range v.([]interface{}) {
			includeList[i] = item.(string)
		}
	}

	// 压缩代码目录
	zipFile, err := zipSourceCode(currentDir, &ZipOptions{
		IncludeList: includeList,
		ExcludeList: excludeList,
	})
	if err != nil {
		return fmt.Errorf("计算源代码哈希值失败: %v", err)
	}
	defer os.Remove(zipFile)

	// 读取压缩文件内容
	content, err := os.ReadFile(zipFile)
	if err != nil {
		return fmt.Errorf("读取压缩文件失败: %v", err)
	}

	// 计算哈希值
	hash := fmt.Sprintf("%x", sha256.Sum256(content))

	// 设置哈希值
	d.SetNew("source_code_hash", hash)

	return nil
}
