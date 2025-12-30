package goauthsdk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RevokeToken 撤销令牌（RFC 7009）
// 调用服务端撤销接口使令牌失效
//
// 参数:
//   - ctx: 上下文，用于控制请求超时等
//   - token: 需要撤销的令牌（access_token 或 refresh_token）
//
// 注意：按 RFC 7009 规范，无论令牌是否存在或已失效，只要服务端返回 HTTP 200 即视为成功
//
// 示例用法:
//
//	// 撤销访问令牌
//	err := client.RevokeToken(context.Background(), accessToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Token revoked successfully")
func (c *Client) RevokeToken(ctx context.Context, token string) error {
	return c.RevokeTokenWithHint(ctx, token, "")
}

// RevokeTokenWithHint 撤销令牌，支持 token_type_hint 参数（RFC 7009）
// 调用服务端撤销接口使令牌失效
//
// 参数:
//   - ctx: 上下文，用于控制请求超时等
//   - token: 需要撤销的令牌
//   - tokenTypeHint: 可选的令牌类型提示，可为 "access_token" 或 "refresh_token"，空字符串表示不传
//
// 注意：按 RFC 7009 规范，无论令牌是否存在或已失效，只要服务端返回 HTTP 200 即视为成功
//
// 示例用法:
//
//	// 撤销刷新令牌
//	err := client.RevokeTokenWithHint(context.Background(), refreshToken, "refresh_token")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Refresh token revoked successfully")
func (c *Client) RevokeTokenWithHint(ctx context.Context, token, tokenTypeHint string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}

	// 构建并发送请求
	req, err := buildRevokeRequest(ctx, c, token, tokenTypeHint)
	if err != nil {
		return err
	}

	// 发送请求
	resp, body, err := doRevokeRequest(c, req)
	if err != nil {
		return err
	}

	// 按 RFC 7009，HTTP 200 即成功（不关心令牌是否存在）
	if resp.StatusCode != http.StatusOK {
		return decodeAPIError(resp, body)
	}

	return nil
}

// buildRevokeRequest 构建撤销请求的 HTTP 请求
func buildRevokeRequest(ctx context.Context, c *Client, token, tokenTypeHint string) (*http.Request, error) {
	// 构建请求 URL
	revokeURL := c.cfg.BackendBaseURL + "/api/v1/oauth/revoke"

	// 构建表单参数
	formData := url.Values{}
	formData.Set("token", token)
	if tokenTypeHint != "" {
		formData.Set("token_type_hint", tokenTypeHint)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", revokeURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create revoke request: %w", err)
	}

	// 设置 Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 设置 Basic Auth（client_id 和 client_secret）
	req.SetBasicAuth(c.cfg.ClientID, c.cfg.ClientSecret)

	return req, nil
}

// doRevokeRequest 发送撤销请求并返回响应与响应体
func doRevokeRequest(c *Client, req *http.Request) (*http.Response, []byte, error) {
	// 发送请求
	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("send revoke request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read revoke response body: %w", err)
	}

	return resp, body, nil
}
