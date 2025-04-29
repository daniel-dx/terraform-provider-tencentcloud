package cloudbase

import (
	"context"

	tcb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcb/v20180608"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TcbService struct {
	client *connectivity.TencentCloudClient
}

// DescribeCloudBaseBuildService 查询云开发构建服务详情
// 参数:
//   - ctx: 上下文信息
//   - name: 服务名称
//   - envId: 环境ID
//   - version: 服务版本号，可选参数
//
// 返回值:
//   - response: 构建服务详情响应
//   - errRet: 错误信息
func (s *TcbService) DescribeCloudBaseBuildService(ctx context.Context, name string, envId string, version *string) (response *tcb.DescribeCloudBaseBuildServiceResponse, errRet error) {
	request := tcb.NewDescribeCloudBaseBuildServiceRequest()

	request.ServiceName = helper.String(name)
	request.EnvId = helper.String(envId)
	if version != nil {
		request.ServiceVersion = version
	}

	ratelimit.Check(request.GetAction())
	response, err := s.client.UseTcbClient().DescribeCloudBaseBuildService(request)
	if err != nil {
		errRet = err
		return
	}

	return response, nil
}
