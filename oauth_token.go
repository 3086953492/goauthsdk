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

// ExchangeToken 使用授权码交换访问令牌
// 在用户授权后，第三方应用会在回调 URL 中收到授权码（code 参数），
// 使用该授权码调用此方法即可获取访问令牌
//
// 参数:
//   - ctx: 上下文，用于控制请求超时等
//   - code: 从回调 URL 中获取的授权码
//
// 示例用法:
//
//	// 从回调 URL 中获取授权码
//	code := r.URL.Query().Get("code")
//
//	// 交换访问令牌
//	token, err := client.ExchangeToken(context.Background(), code)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 使用访问令牌
//	fmt.Printf("Access Token: %s\n", token.AccessToken)
//	fmt.Printf("Expires In: %d seconds\n", token.ExpiresIn)
func (c *Client) ExchangeToken(ctx context.Context, code string) (*TokenResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("code is required")
	}

	// 构建并发送请求
	req, err := buildTokenRequest(ctx, c, code)
	if err != nil {
		return nil, err
	}

	resp, body, err := doTokenRequest(c, req)
	if err != nil {
		return nil, err
	}

	// 解析响应
	token, err := parseTokenResponse(resp, body)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// buildTokenRequest 构建 token 交换的 HTTP 请求
func buildTokenRequest(ctx context.Context, c *Client, code string) (*http.Request, error) {
	// 构建请求 URL
	tokenURL := c.cfg.BackendBaseURL + "/api/v1/oauth/token"

	// 构建表单参数
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", code)
	formData.Set("redirect_uri", c.cfg.RedirectURI)

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置 Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 设置 Basic Auth（client_id 和 client_secret）
	req.SetBasicAuth(c.cfg.ClientID, c.cfg.ClientSecret)

	return req, nil
}

// doTokenRequest 发送 token 请求并返回响应与响应体
func doTokenRequest(c *Client, req *http.Request) (*http.Response, []byte, error) {
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

// parseTokenResponse 解析 token 响应，检查业务成功和 HTTP 状态码
func parseTokenResponse(resp *http.Response, body []byte) (*TokenResponse, error) {
	// 解析响应
	var apiResp apiResponse[TokenResponse]
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response (status: %d, body: %s): %w", resp.StatusCode, string(body), err)
	}

	// 检查业务是否成功
	if !apiResp.Success {
		return nil, fmt.Errorf("token exchange failed (status: %d): %s", resp.StatusCode, apiResp.Message)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, apiResp.Message)
	}

	return &apiResp.Data, nil
}
