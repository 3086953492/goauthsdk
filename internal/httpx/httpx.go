package httpx

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPDoer 是发送 HTTP 请求的最小接口
// *http.Client 自动满足此接口
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Do 发送 HTTP 请求并读取响应体
// 调用方无需关心 resp.Body.Close，本函数统一处理
//
// 返回值:
//   - *http.Response: 响应（Body 已被读取并关闭，但 StatusCode/Header 等仍可用）
//   - []byte: 响应体内容
//   - error: 发送请求或读取响应体失败时返回错误
func Do(doer HTTPDoer, req *http.Request) (*http.Response, []byte, error) {
	resp, err := doer.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read response body: %w", err)
	}

	return resp, body, nil
}
