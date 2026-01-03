package goauthsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/3086953492/goauthsdk/internal/httpx"
)

// GetUser 根据 sub 获取用户详情
// 调用 GET /api/v1/users/sub/{sub} 接口，需传入有效的 access_token（通常为 client_credentials 模式获取）
//
// 参数:
//   - ctx: 上下文，用于控制请求超时等
//   - accessToken: 有效的访问令牌（需包含 profile scope）
//   - sub: 用户唯一标识（subject）
//
// 返回:
//   - *UserDetail: 用户详情
//   - error: 失败时返回错误；若为 401/403/404 等，可通过 errors.As 获取 *APIError
//
// 示例用法:
//
//	// 先获取 client_credentials token（需包含 profile scope）
//	tokenResp, err := client.ClientCredentialsToken(context.Background(), "profile")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 获取用户详情（sub 可从 userinfo.Sub 或 JWT claims 中获取）
//	user, err := client.GetUser(context.Background(), tokenResp.AccessToken, "user-sub-xxx")
//	if err != nil {
//	    // 可通过 errors.As 获取结构化错误信息
//	    var apiErr *goauthsdk.APIError
//	    if errors.As(err, &apiErr) {
//	        fmt.Printf("错误码: %s, 详情: %s\n", apiErr.Code, apiErr.Detail)
//	    }
//	    log.Fatal(err)
//	}
//	fmt.Printf("用户ID: %d, 用户名: %s, 昵称: %s\n", user.ID, user.Username, user.Nickname)
func (c *Client) GetUser(ctx context.Context, accessToken string, sub string) (*UserDetail, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access_token is required")
	}
	if sub == "" {
		return nil, fmt.Errorf("sub is required")
	}

	// 构建请求
	req, err := buildGetUserRequest(ctx, c, accessToken, sub)
	if err != nil {
		return nil, err
	}

	// 发送请求
	resp, body, err := doGetUserRequest(c, req)
	if err != nil {
		return nil, err
	}

	// 解析响应
	return parseGetUserResponse(resp, body)
}

// buildGetUserRequest 构建获取用户详情的 HTTP 请求
func buildGetUserRequest(ctx context.Context, c *Client, accessToken string, sub string) (*http.Request, error) {
	// 构建请求 URL（对 sub 做 path escape 防止特殊字符破坏路径）
	userURL := fmt.Sprintf("%s/api/v1/users/sub/%s", c.cfg.BackendBaseURL, url.PathEscape(sub))

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", userURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create get user request: %w", err)
	}

	// 设置 Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return req, nil
}

// doGetUserRequest 发送获取用户详情请求并返回响应与响应体
func doGetUserRequest(c *Client, req *http.Request) (*http.Response, []byte, error) {
	return httpx.Do(c.cfg.HTTPClient, req)
}

// parseGetUserResponse 解析获取用户详情响应
// 成功时响应格式：{ "code": 0, "message": "...", "data": {...} }
// 错误时响应格式：{ "type": "...", "title": "...", "status": ..., "detail": "..." }
func parseGetUserResponse(resp *http.Response, body []byte) (*UserDetail, error) {
	// 非 2xx：统一走 decodeAPIError
	if resp.StatusCode != http.StatusOK {
		return nil, decodeAPIError(resp, body)
	}

	// 解析响应
	var apiResp apiCodeResponse[UserDetail]
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse user response: %w", err)
	}

	// 检查业务是否成功（code == 0 表示成功）
	if apiResp.Code != 0 {
		return nil, newBusinessError(resp.StatusCode, apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}
