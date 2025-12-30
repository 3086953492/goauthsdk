package goauthsdk

import (
	"fmt"
	"net/url"
)

// BuildAuthorizationURL 构建用户授权时跳转的前端 URL
// 用户浏览器应重定向到该 URL，在前端授权确认页点击"确认授权"后，
// 前端会再跳转到后端 /api/v1/oauth/authorization 完成授权码生成
//
// 参数:
//   - state: 可选的状态参数，用于防止 CSRF 攻击
//   - scope: 可选的权限范围，多个 scope 用空格分隔
//
// 示例用法:
//
//	authURL, err := client.BuildAuthorizationURL("random-state-string", "read write")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// 将用户浏览器重定向到 authURL
//	http.Redirect(w, r, authURL, http.StatusFound)
func (c *Client) BuildAuthorizationURL(state, scope string) (string, error) {
	// 构造前端授权确认页地址
	u, err := url.Parse(c.cfg.FrontendBaseURL + "/oauth/authorize")
	if err != nil {
		return "", fmt.Errorf("parse frontend base url: %w", err)
	}

	// 构建 query 参数
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", c.cfg.ClientID)
	q.Set("redirect_uri", c.cfg.RedirectURI)

	if scope != "" {
		q.Set("scope", scope)
	}

	if state != "" {
		q.Set("state", state)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
