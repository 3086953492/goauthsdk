package goauthsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// IntrospectToken 内省访问令牌
// 调用服务端内省接口验证令牌的有效性，并返回令牌相关信息
//
// 参数:
//   - ctx: 上下文，用于控制请求超时等
//   - token: 需要验证的访问令牌
//
// 示例用法:
//
//	// 验证访问令牌
//	resp, err := client.IntrospectToken(context.Background(), accessToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 检查令牌是否有效
//	if resp.Active {
//	    fmt.Printf("Token is valid, expires at: %d\n", resp.Exp)
//	} else {
//	    fmt.Println("Token is invalid or expired")
//	}
func (c *Client) IntrospectToken(ctx context.Context, token string) (*IntrospectionResponse, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// 构建并发送请求
	req, err := buildIntrospectRequest(ctx, c, token)
	if err != nil {
		return nil, err
	}

	resp, body, err := doIntrospectRequest(c, req)
	if err != nil {
		return nil, err
	}

	// 解析响应
	introspection, err := parseIntrospectResponse(resp, body)
	if err != nil {
		return nil, err
	}

	return introspection, nil
}

// buildIntrospectRequest 构建内省请求的 HTTP 请求
func buildIntrospectRequest(ctx context.Context, c *Client, token string) (*http.Request, error) {
	// 构建请求 URL
	introspectURL := c.cfg.BackendBaseURL + "/api/v1/oauth/introspect"

	// 构建表单参数
	formData := url.Values{}
	formData.Set("token", token)

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", introspectURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置 Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 设置 Basic Auth（client_id 和 client_secret）
	req.SetBasicAuth(c.cfg.ClientID, c.cfg.ClientSecret)

	return req, nil
}

// doIntrospectRequest 发送内省请求并返回响应与响应体
func doIntrospectRequest(c *Client, req *http.Request) (*http.Response, []byte, error) {
	// 发送请求
	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, body, nil
}

// parseIntrospectResponse 解析内省响应，检查业务成功和 HTTP 状态码
func parseIntrospectResponse(resp *http.Response, body []byte) (*IntrospectionResponse, error) {
	// 解析响应
	var apiResp apiResponse[IntrospectionResponse]
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response (status: %d, body: %s): %w", resp.StatusCode, string(body), err)
	}

	// 检查业务是否成功
	if !apiResp.Success {
		return nil, fmt.Errorf("token introspection failed (status: %d): %s", resp.StatusCode, apiResp.Message)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, apiResp.Message)
	}

	return &apiResp.Data, nil
}
