package goauthsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetUser 根据用户 ID 获取用户详情
// 调用 GET /api/v1/users/{id} 接口，需传入有效的 access_token（通常为 client_credentials 模式获取）
//
// 参数:
//   - ctx: 上下文，用于控制请求超时等
//   - accessToken: 有效的访问令牌（需包含 profile scope）
//   - userID: 用户 ID
//
// 返回:
//   - *UserDetail: 用户详情
//   - error: 失败时返回错误；若为 401/403/404 等，返回的 error 可断言为 *ProblemDetails
//
// 示例用法:
//
//	// 先获取 client_credentials token（需包含 profile scope）
//	tokenResp, err := client.ClientCredentialsToken(context.Background(), "profile")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 获取用户详情
//	user, err := client.GetUser(context.Background(), tokenResp.AccessToken, 123)
//	if err != nil {
//	    // 可尝试断言为 *goauthsdk.ProblemDetails 获取具体错误码
//	    if pd, ok := err.(*goauthsdk.ProblemDetails); ok {
//	        fmt.Printf("错误码: %s, 详情: %s\n", pd.Title, pd.Detail)
//	    }
//	    log.Fatal(err)
//	}
//	fmt.Printf("用户ID: %d, 用户名: %s, 昵称: %s\n", user.ID, user.Username, user.Nickname)
func (c *Client) GetUser(ctx context.Context, accessToken string, userID uint64) (*UserDetail, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access_token is required")
	}
	if userID == 0 {
		return nil, fmt.Errorf("user_id is required")
	}

	// 构建请求
	req, err := buildGetUserRequest(ctx, c, accessToken, userID)
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
func buildGetUserRequest(ctx context.Context, c *Client, accessToken string, userID uint64) (*http.Request, error) {
	// 构建请求 URL
	userURL := fmt.Sprintf("%s/api/v1/users/%d", c.cfg.BackendBaseURL, userID)

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
	// 发送请求
	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("send get user request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read get user response body: %w", err)
	}

	return resp, body, nil
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
