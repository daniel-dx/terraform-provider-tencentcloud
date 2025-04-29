package cloudbase

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

// UploadOptions 定义上传文件的选项
type UploadOptions struct {
	Action         string            // 上传地址
	File           string            // 文件路径
	Method         string            // HTTP 方法
	Headers        map[string]string // 请求头
	WithCredential bool              // 是否携带认证信息
}

// UploadResult 定义上传结果
type UploadResult struct {
	Data    interface{}            // 响应数据
	Context map[string]interface{} // 上下文信息
}

// ZipOptions 定义压缩文件的选项
type ZipOptions struct {
	IncludeList []string // 指定要包含的文件或目录的匹配模式列表，为空时包含所有文件
	ExcludeList []string // 指定要排除的文件或目录的匹配模式列表，匹配的文件或目录将被忽略
}

// zipSourceCode 将源代码目录压缩为zip文件，并返回临时zip文件的路径
// sourceDir: 源代码目录路径
// opts: 压缩选项，可以为 nil，使用默认配置
func zipSourceCode(sourceDir string, opts *ZipOptions) (string, error) {
	if opts == nil {
		opts = &ZipOptions{}
	}

	// 创建临时zip文件
	tmpFile, err := os.CreateTemp("", "cloudbase-*.zip")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer tmpFile.Close()

	zipWriter := zip.NewWriter(tmpFile)
	defer zipWriter.Close()

	// 遍历源目录并压缩
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("获取相对路径失败: %v", err)
		}

		// 跳过根目录
		if relPath == "." {
			return nil
		}

		// 检查是否匹配 exclude 规则
		for _, pattern := range opts.ExcludeList {
			matched, m_err := doublestar.Match(pattern, relPath)
			if m_err != nil {
				return fmt.Errorf("匹配排除规则失败: %v", m_err)
			}
			if matched {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// 如果有 include 规则，检查是否匹配
		if len(opts.IncludeList) > 0 {
			isIncluded := false
			for _, pattern := range opts.IncludeList {
				matched, m_err := doublestar.Match(pattern, relPath)
				if m_err != nil {
					return fmt.Errorf("匹配包含规则失败: %v", m_err)
				}
				if matched {
					isIncluded = true
					break
				}
			}
			if !isIncluded {
				if info.IsDir() {
					// 目录不匹配但仍需继续遍历
					return nil
				}
				return nil
			}
		}

		// 创建文件条目
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("创建文件头失败: %v", err)
		}
		header.Name = relPath
		if info.IsDir() {
			header.Name += "/"
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("创建zip条目失败: %v", err)
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("打开文件失败: %v", err)
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return fmt.Errorf("写入文件失败: %v", err)
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("压缩文件失败: %v", err)
	}

	return tmpFile.Name(), nil
}

// uploadFile 上传文件到指定地址
func uploadFile(opts UploadOptions) (*UploadResult, error) {
	// 打开文件
	file, err := os.Open(opts.File)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 创建请求
	method := opts.Method
	if method == "" {
		method = "POST"
	}
	req, err := http.NewRequest(method, opts.Action, file)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 构造返回结果
	result := &UploadResult{
		Context: map[string]interface{}{
			"statusCode": resp.StatusCode,
			"status":     resp.Status,
			"headers":    resp.Header,
			"file":       opts.File,
		},
	}

	// 尝试解析 JSON
	if len(body) > 0 {
		var data interface{}
		if err := json.Unmarshal(body, &data); err == nil {
			result.Data = data
		}
	}

	// 检查响应状态
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP 错误: %s", resp.Status)
	}

	return result, nil
}
