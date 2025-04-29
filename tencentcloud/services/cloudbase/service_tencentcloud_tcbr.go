package cloudbase

import (
	"context"

	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	tcbr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcbr/v20220217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TcbrService struct {
	client *connectivity.TencentCloudClient
}

// DescribeCloudRunServerDetail 查询云托管服务详情
// 参数:
//   - ctx: 上下文信息
//   - serverName: 服务名称
//   - envId: 环境ID
//
// 返回值:
//   - serverInfo: 服务详情信息
//   - errRet: 错误信息
func (s *TcbrService) DescribeCloudRunServerDetail(ctx context.Context, serverName string, envId string) (serverInfo *tcbr.DescribeCloudRunServerDetailResponse, errRet error) {
	request := tcbr.NewDescribeCloudRunServerDetailRequest()

	request.ServerName = helper.String(serverName)
	request.EnvId = helper.String(envId)

	ratelimit.Check(request.GetAction())
	response, err := s.client.UseTcbrClient().DescribeCloudRunServerDetail(request)
	if err != nil {
		errRet = err
		return
	}

	serverInfo = response
	return
}

// CreateCloudRunServer 创建云托管服务
// 参数:
//   - ctx: 上下文信息
//   - request: 创建服务请求参数
//
// 返回值:
//   - response: 创建服务响应
//   - errRet: 错误信息
func (s *TcbrService) CreateCloudRunServer(ctx context.Context, request *tcbr.CreateCloudRunServerRequest) (response *tcbr.CreateCloudRunServerResponse, errRet error) {
	ratelimit.Check(request.GetAction())
	response, err := s.client.UseTcbrClient().CreateCloudRunServer(request)
	if err != nil {
		errRet = err
		return
	}
	return response, nil
}

// UpdateCloudRunServer 更新云托管服务
// 参数:
//   - ctx: 上下文信息
//   - request: 更新服务请求参数
//
// 返回值:
//   - response: 更新服务响应
//   - errRet: 错误信息
func (s *TcbrService) UpdateCloudRunServer(ctx context.Context, request *tcbr.UpdateCloudRunServerRequest) (response *tcbr.UpdateCloudRunServerResponse, errRet error) {
	ratelimit.Check(request.GetAction())
	response, err := s.client.UseTcbrClient().UpdateCloudRunServer(request)
	if err != nil {
		errRet = err
		return
	}
	return response, nil
}

// DeleteCloudRunServer 删除云托管服务
// 参数:
//   - ctx: 上下文信息
//   - envId: 环境ID
//   - serverName: 服务名称
//
// 返回值:
//   - *struct{RequestId *string}: 包含请求ID的响应结构
//   - error: 错误信息
func (s *TcbrService) DeleteCloudRunServer(ctx context.Context, envId, serverName string) (*struct{ RequestId *string }, error) {
	// 手动构造请求
	request := tchttp.NewCommonRequest("tcbr", "2022-02-17", "DeleteCloudRunServer")
	request.SetActionParameters(map[string]interface{}{
		"EnvId":      envId,
		"ServerName": serverName,
	})

	// 创建响应结构
	response := &struct {
		tchttp.BaseResponse
		Response struct {
			RequestId string `json:"RequestId"`
		} `json:"Response"`
	}{}

	// 发送请求
	err := s.client.UseTcbrClient().Send(request, response)
	if err != nil {
		return nil, err
	}

	return &struct{ RequestId *string }{
		RequestId: &response.Response.RequestId,
	}, nil
}
