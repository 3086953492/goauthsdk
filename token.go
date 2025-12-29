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
// 接口返回格式：{ "code": 0, "message": "...", "data": {...} }（code == 0 表示成功）
func parseTokenResponse(resp *http.Response, body []byte) (*TokenResponse, error) {
	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with HTTP %d: %s", resp.StatusCode, truncateBody(body))
	}

	// 解析响应
	var apiResp apiCodeResponse[TokenResponse]
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, truncateBody(body))
	}

	// 检查业务是否成功（code == 0 表示成功）
	if apiResp.Code != 0 {
		return nil, fmt.Errorf("token exchange failed with code %d: %s", apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}

// truncateBody 截断响应体用于错误输出，避免日志过长
func truncateBody(body []byte) string {
	const maxLen = 200
	if len(body) <= maxLen {
		return string(body)
	}
	return string(body[:maxLen]) + "..."
}

// RefreshToken 使用刷新令牌获取新的访问令牌
// 当访问令牌过期时，可以使用之前获取的刷新令牌来获取新的访问令牌，
// 无需用户重新授权
//
// 参数:
//   - ctx: 上下文，用于控制请求超时等
//   - refreshToken: 之前获取的刷新令牌
//
// 示例用法:
//
//	// 使用刷新令牌获取新的访问令牌
//	newToken, err := client.RefreshToken(context.Background(), token.RefreshToken.RefreshToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 使用新的访问令牌
//	fmt.Printf("New Access Token: %s\n", newToken.AccessToken.AccessToken)
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh_token is required")
	}

	// 构建并发送请求
	req, err := buildRefreshTokenRequest(ctx, c, refreshToken)
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

// buildRefreshTokenRequest 构建刷新令牌的 HTTP 请求
func buildRefreshTokenRequest(ctx context.Context, c *Client, refreshToken string) (*http.Request, error) {
	// 构建请求 URL
	tokenURL := c.cfg.BackendBaseURL + "/api/v1/oauth/token"

	// 构建表单参数
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", refreshToken)

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

// ClientCredentialsToken 使用客户端凭证模式获取访问令牌
// 该模式适用于服务端到服务端的机密通信，无用户上下文
//
// 参数:
//   - ctx: 上下文，用于控制请求超时等
//   - scope: 请求的权限范围；为空时服务端返回的 scope 也为空，不会自动赋默认值
//
// 注意事项:
//   - 该模式下 JWT 的 sub 固定为 "client:<client_id>"
//   - 使用该 token 调用 IntrospectToken 时，返回 active=true 但不含 username/sub
//   - 该 token 不适用于 UserInfo 接口（因为无用户上下文）
//
// 示例用法:
//
//	token, err := client.ClientCredentialsToken(context.Background(), "api")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 使用访问令牌调用业务 API
//	fmt.Printf("Access Token: %s\n", token.AccessToken)
//	fmt.Printf("Expires In: %d seconds\n", token.ExpiresIn)
func (c *Client) ClientCredentialsToken(ctx context.Context, scope string) (*ClientCredentialsTokenResponse, error) {
	// 构建并发送请求
	req, err := buildClientCredentialsTokenRequest(ctx, c, scope)
	if err != nil {
		return nil, err
	}

	resp, body, err := doTokenRequest(c, req)
	if err != nil {
		return nil, err
	}

	// 解析响应
	token, err := parseClientCredentialsTokenResponse(resp, body)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// buildClientCredentialsTokenRequest 构建客户端凭证模式的 HTTP 请求
func buildClientCredentialsTokenRequest(ctx context.Context, c *Client, scope string) (*http.Request, error) {
	// 构建请求 URL
	tokenURL := c.cfg.BackendBaseURL + "/api/v1/oauth/token"

	// 构建表单参数
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	if scope != "" {
		formData.Set("scope", scope)
	}

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

// parseClientCredentialsTokenResponse 解析客户端凭证模式的 token 响应
// 响应格式：{ "code": 0, "message": "...", "data": { "access_token": "...", "expires_in": ..., "token_type": "...", "scope": "..." } }
func parseClientCredentialsTokenResponse(resp *http.Response, body []byte) (*ClientCredentialsTokenResponse, error) {
	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client credentials token request failed with HTTP %d: %s", resp.StatusCode, truncateBody(body))
	}

	// 解析响应
	var apiResp apiCodeResponse[ClientCredentialsTokenResponse]
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, truncateBody(body))
	}

	// 检查业务是否成功（code == 0 表示成功）
	if apiResp.Code != 0 {
		return nil, fmt.Errorf("client credentials token request failed with code %d: %s", apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}
