package goauthsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// UserInfo 获取当前访问令牌对应的用户信息
// 调用 GET /api/v1/oauth/userinfo 接口，需传入有效的 access_token
//
// 参数:
//   - ctx: 上下文，用于控制请求超时等
//   - accessToken: 有效的访问令牌
//
// 返回:
//   - *UserInfo: 用户信息（sub、nickname、picture、updated_at）
//   - error: 失败时返回错误；若为 401/403/404 等，返回的 error 可断言为 *ProblemDetails
//
// 示例用法:
//
//	info, err := client.UserInfo(context.Background(), token.AccessToken.AccessToken)
//	if err != nil {
//	    // 可尝试断言为 *goauthsdk.ProblemDetails 获取具体错误码
//	    if pd, ok := err.(*goauthsdk.ProblemDetails); ok {
//	        fmt.Printf("错误码: %s, 详情: %s\n", pd.Code, pd.Detail)
//	    }
//	    log.Fatal(err)
//	}
//	fmt.Printf("用户ID: %s, 昵称: %s\n", info.Sub, info.Nickname)
func (c *Client) UserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("accessToken is required")
	}

	// 构建请求
	req, err := buildUserInfoRequest(ctx, c, accessToken)
	if err != nil {
		return nil, err
	}

	// 发送请求
	resp, body, err := doUserInfoRequest(c, req)
	if err != nil {
		return nil, err
	}

	// 解析响应
	return parseUserInfoResponse(resp, body)
}

// buildUserInfoRequest 构建获取用户信息的 HTTP 请求
func buildUserInfoRequest(ctx context.Context, c *Client, accessToken string) (*http.Request, error) {
	// 构建请求 URL
	userInfoURL := c.cfg.BackendBaseURL + "/api/v1/oauth/userinfo"

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置 Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	return req, nil
}

// doUserInfoRequest 发送用户信息请求并返回响应与响应体
func doUserInfoRequest(c *Client, req *http.Request) (*http.Response, []byte, error) {
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

// parseUserInfoResponse 解析用户信息响应
// 成功时响应格式：{ "code": 0, "message": "...", "data": {...} }
// 错误时响应格式：{ "type": "...", "title": "...", "status": ..., "code": "...", "detail": "..." }
func parseUserInfoResponse(resp *http.Response, body []byte) (*UserInfo, error) {
	// 成功响应
	if resp.StatusCode == http.StatusOK {
		var apiResp apiCodeResponse[UserInfo]
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("failed to parse userinfo response: %w (body: %s)", err, truncateBody(body))
		}

		// 检查业务是否成功（code == 0 表示成功）
		if apiResp.Code != 0 {
			return nil, fmt.Errorf("userinfo request failed with code %d: %s", apiResp.Code, apiResp.Message)
		}

		return &apiResp.Data, nil
	}

	// 错误响应：尝试解析为 ProblemDetails（兼容 code 或 title 任一存在的情况）
	var pd ProblemDetails
	if err := json.Unmarshal(body, &pd); err == nil && (pd.Code != "" || pd.Title != "") {
		return nil, &pd
	}

	// 回退：返回通用错误
	return nil, fmt.Errorf("userinfo request failed with HTTP %d: %s", resp.StatusCode, truncateBody(body))
}
