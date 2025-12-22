# goauthsdk

`goauthsdk` 是一个 Go 语言 SDK，用于对接 **goauth** 服务的 OAuth2 授权码模式（Authorization Code）常用流程：

- 构建用户授权跳转 URL（前端授权确认页）
- 使用授权码（code）交换访问令牌（access token）
- 使用刷新令牌（refresh token）刷新访问令牌
- 内省（introspect）令牌有效性（RFC 7662）
- 撤销（revoke）令牌（RFC 7009）
- 获取用户信息（userinfo）

## 安装

```bash
go get github.com/3086953492/goauthsdk
```

要求：Go 1.21+

## 基本概念与 URL 说明

SDK 的配置分为两类 BaseURL：

- **FrontendBaseURL**：goauth 的前端站点地址（用于拼接用户授权确认页 `GET /oauth/authorize`）。
- **BackendBaseURL**：goauth 的后端服务地址（用于调用实际接口）：
  - `POST /api/v1/oauth/token` - 换取/刷新访问令牌
  - `POST /api/v1/oauth/introspect` - 令牌内省（RFC 7662）
  - `POST /api/v1/oauth/revoke` - 令牌撤销（RFC 7009）
  - `GET /api/v1/oauth/userinfo` - 获取用户信息

### 接口响应格式

Token 接口返回格式为 `{ "code": 0, "message": "...", "data": {...} }`（`code == 0` 表示成功）：

```json
{
  "code": 0,
  "message": "交换访问令牌成功",
  "data": {
    "access_token": {
      "access_token": "xxx",
      "expires_in": 3600
    },
    "refresh_token": {
      "refresh_token": "xxx",
      "expires_in": 604800
    },
    "token_type": "Bearer",
    "scope": "read write"
  }
}
```

Introspect 接口返回格式（data 部分符合 RFC 7662）：

```json
{
  "code": 0,
  "message": "内省成功",
  "data": {
    "active": true,
    "scope": "read write",
    "client_id": "xxx",
    "username": "user@example.com",
    "token_type": "Bearer",
    "exp": 1703232000,
    "sub": "user_id"
  }
}
```

### 典型授权码流程

1. 你的应用把用户浏览器重定向到 `BuildAuthorizationURL(...)` 生成的地址（goauth 前端授权确认页）。
2. 用户在授权页面确认授权后，goauth 会重定向回你的 **RedirectURI**，并携带 `code`（以及可选的 `state`）。
3. 你的后端服务拿到 `code` 后调用 `ExchangeToken(...)` 交换访问令牌。

## 快速开始

### 1) 初始化 Client

```go
package main

import (
	"log"

	"github.com/3086953492/goauthsdk"
)

func main() {
	client, err := goauthsdk.NewClient(goauthsdk.Config{
		FrontendBaseURL: "https://portal.example.com",
		BackendBaseURL:  "https://auth.example.com",
		ClientID:        "your-client-id",
		ClientSecret:    "your-client-secret",
		RedirectURI:     "https://yourapp.com/callback",
	})
	if err != nil {
		log.Fatal(err)
	}

	_ = client
}
```

### 2) 生成授权 URL 并重定向用户

```go
authURL, err := client.BuildAuthorizationURL(
	"random-state-string", // 建议传，用于防 CSRF
	"read write",          // 可选，多个 scope 用空格分隔
)
if err != nil {
	// handle error
}

// 将用户浏览器重定向到 authURL（示例：net/http）
// http.Redirect(w, r, authURL, http.StatusFound)
```

### 3) 回调接口：用 code 交换 Token

你的回调地址（RedirectURI）会收到 `code` 参数，拿到后调用 `ExchangeToken`：

```go
import (
	"context"
	"net/http"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	token, err := client.ExchangeToken(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// token.AccessToken.AccessToken
	// token.AccessToken.ExpiresIn
	// token.RefreshToken.RefreshToken
	// token.RefreshToken.ExpiresIn
	// token.TokenType
	// token.Scope
	_ = token
}
```

### 4) 刷新访问令牌

```go
newToken, err := client.RefreshToken(context.Background(), token.RefreshToken.RefreshToken)
if err != nil {
	// handle error
}
_ = newToken
```

### 5) 内省令牌（RFC 7662）

```go
// 内省访问令牌
resp, err := client.IntrospectToken(context.Background(), token.AccessToken.AccessToken)
if err != nil {
	// handle error
}

if resp.Active {
	// resp.Exp / resp.Scope / resp.ClientID / resp.Username / resp.Sub ...
} else {
	// 无效或已过期
}

// 也可以使用 token_type_hint 指定令牌类型
resp, err = client.IntrospectTokenWithHint(context.Background(), refreshToken, "refresh_token")
```

### 6) 撤销令牌（RFC 7009）

```go
// 撤销访问令牌
err := client.RevokeToken(context.Background(), token.AccessToken.AccessToken)
if err != nil {
	// handle error
}

// 也可以使用 token_type_hint 指定令牌类型
err = client.RevokeTokenWithHint(context.Background(), refreshToken, "refresh_token")
```

### 7) 获取用户信息

```go
// 使用访问令牌获取用户信息
info, err := client.UserInfo(context.Background(), token.AccessToken.AccessToken)
if err != nil {
	// 可尝试断言为 *goauthsdk.ProblemDetails 获取具体错误码
	if pd, ok := err.(*goauthsdk.ProblemDetails); ok {
		// pd.Code 可能为 INVALID_TOKEN、INSUFFICIENT_SCOPE、USER_NOT_FOUND 等
		fmt.Printf("错误码: %s, 详情: %s\n", pd.Code, pd.Detail)
	}
	// handle error
}

// info.Sub       - 用户唯一标识（用户ID）
// info.Nickname  - 用户昵称
// info.Picture   - 用户头像URL
// info.UpdatedAt - 用户信息更新时间（Unix 时间戳）
fmt.Printf("用户ID: %s, 昵称: %s\n", info.Sub, info.Nickname)
```

## 自定义 HTTPClient（可选）

你可以传入自定义 `HTTPClient`（例如设置超时、代理、TLS 等）：

```go
import (
	"net/http"
	"time"
)

client, err := goauthsdk.NewClient(goauthsdk.Config{
	FrontendBaseURL: "https://portal.example.com",
	BackendBaseURL:  "https://auth.example.com",
	ClientID:        "your-client-id",
	ClientSecret:    "your-client-secret",
	RedirectURI:     "https://yourapp.com/callback",
	HTTPClient: &http.Client{
		Timeout: 10 * time.Second,
	},
})
```

> 说明：SDK 调用 token、introspect 和 revoke 接口时会自动使用 Basic Auth（`client_id` / `client_secret`）。

## 运行本仓库的手工测试服务（可选）

仓库自带一个用于开发/测试的手工验证服务：`test/main.go`，包含完整流程的路由（`/auth`、`/callback`、`/introspect`、`/refresh`、`/revoke`）。

```bash
go run ./test
```

启动后访问 `http://localhost:7000/` 查看说明，然后：

- 访问 `http://localhost:7000/auth` 发起授权
- 授权完成后会自动回跳到 `/callback` 并展示 token
- 用 `/introspect?token=xxx` 内省 token
- 用 `/introspect?token=xxx&token_type_hint=refresh_token` 内省刷新令牌
- 用 `/refresh?refresh_token=xxx` 刷新 token
- 用 `/revoke?token=xxx` 撤销 token
- 用 `/revoke?token=xxx&token_type_hint=refresh_token` 撤销刷新令牌
- 用 `/userinfo?token=xxx` 获取用户信息

## 常见注意事项

- **state 建议必传**：用于防止 CSRF，回调时校验 `state` 是否与发起时一致。
- **RedirectURI 必须完全一致**：需与 goauth 后台注册的回调地址匹配。
- **生产环境务必使用 HTTPS**：避免 code/token 在传输过程中被窃取。
- **不要在日志中打印完整 token/secret**：如需排查，建议仅打印前缀（仓库测试代码已做了前缀截断示例）。
