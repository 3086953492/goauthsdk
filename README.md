# goauthsdk

`goauthsdk` 是一个 Go 语言 SDK，用于对接 **goauth** 服务的 OAuth2 常用流程：

- 授权码模式（Authorization Code）：构建用户授权跳转 URL、授权码交换令牌、刷新令牌
- 客户端凭证模式（Client Credentials）：服务端到服务端的机密通信
- 内省（introspect）令牌有效性（RFC 7662）
- 撤销（revoke）令牌（RFC 7009）
- 获取用户信息（userinfo）
- 根据用户 ID 获取用户详情
- 离线验证令牌（基于 JWT 签名验签，无需调用服务端）

## 安装

```bash
go get github.com/3086953492/goauthsdk
```

要求：Go 1.23+

## 基本概念与 URL 说明

SDK 的配置分为两类 BaseURL：

- **FrontendBaseURL**：goauth 的前端站点地址（用于拼接用户授权确认页 `GET /oauth/authorize`）。
- **BackendBaseURL**：goauth 的后端服务地址（用于调用实际接口）：
  - `POST /api/v1/oauth/token` - 换取/刷新访问令牌、客户端凭证模式
  - `POST /api/v1/oauth/introspect` - 令牌内省（RFC 7662）
  - `POST /api/v1/oauth/revoke` - 令牌撤销（RFC 7009）
  - `GET /api/v1/oauth/userinfo` - 获取当前用户信息
  - `GET /api/v1/users/{id}` - 获取指定用户详情

### 接口响应格式

Token 接口返回格式为 `{ "code": 0, "message": "...", "data": {...} }`（`code == 0` 表示成功）：

**授权码模式 / 刷新令牌响应：**

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

**客户端凭证模式响应：**

```json
{
  "code": 0,
  "message": "获取访问令牌成功",
  "data": {
    "access_token": "xxx",
    "expires_in": 3600,
    "token_type": "Bearer",
    "scope": "api"
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

### 5) 客户端凭证模式

适用于服务端到服务端的机密通信，无用户上下文：

```go
// 使用客户端凭证模式获取访问令牌
token, err := client.ClientCredentialsToken(context.Background(), "api")
if err != nil {
	// handle error
}

// token.AccessToken  - JWT 访问令牌
// token.ExpiresIn    - 过期时间（秒）
// token.TokenType    - 令牌类型，通常为 "Bearer"
// token.Scope        - 授权范围
fmt.Printf("Access Token: %s\n", token.AccessToken)
```

> **注意**：
> - 该模式下 JWT 的 sub 固定为 "client:<client_id>"
> - 使用该 token 调用 IntrospectToken 时，返回 active=true 但不含 username/sub
> - 该 token 不适用于 UserInfo 接口（因为无用户上下文）

### 6) 内省令牌（RFC 7662）

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

### 7) 撤销令牌（RFC 7009）

```go
// 撤销访问令牌
err := client.RevokeToken(context.Background(), token.AccessToken.AccessToken)
if err != nil {
	// handle error
}

// 也可以使用 token_type_hint 指定令牌类型
err = client.RevokeTokenWithHint(context.Background(), refreshToken, "refresh_token")
```

### 8) 获取用户信息

使用访问令牌获取当前授权用户的基本信息：

```go
info, err := client.UserInfo(context.Background(), token.AccessToken.AccessToken)
if err != nil {
	// 可通过 errors.As 获取结构化错误信息
	var apiErr *goauthsdk.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("错误码: %s, 详情: %s\n", apiErr.Code, apiErr.Detail)
	}
	// handle error
}

// info.Sub       - 用户唯一标识（用户ID）
// info.Nickname  - 用户昵称
// info.Picture   - 用户头像URL
// info.UpdatedAt - 用户信息更新时间（Unix 时间戳）
fmt.Printf("用户ID: %s, 昵称: %s\n", info.Sub, info.Nickname)
```

### 9) 获取用户详情

根据用户 ID 获取用户的详细信息（需要 client_credentials 模式的 token）：

```go
// 先获取 client_credentials token（需包含 profile scope）
ccToken, err := client.ClientCredentialsToken(context.Background(), "profile")
if err != nil {
	log.Fatal(err)
}

// 获取用户详情
user, err := client.GetUser(context.Background(), ccToken.AccessToken, 123)
if err != nil {
	var apiErr *goauthsdk.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("错误码: %s, 详情: %s\n", apiErr.Code, apiErr.Detail)
	}
	log.Fatal(err)
}

// user.ID        - 用户主键 ID
// user.Subject   - 用户唯一标识（对外使用）
// user.Username  - 用户名
// user.Nickname  - 昵称
// user.Avatar    - 头像 URL
// user.Status    - 状态：1=正常，0=禁用
// user.Role      - 角色：user / admin
// user.CreatedAt - 创建时间（ISO 8601 格式）
// user.UpdatedAt - 更新时间（ISO 8601 格式）
fmt.Printf("用户ID: %d, 用户名: %s, 昵称: %s\n", user.ID, user.Username, user.Nickname)
```

## 错误处理

SDK 统一使用 `*APIError` 类型返回 API 错误，可通过 `errors.As` 获取结构化错误信息：

```go
import "errors"

info, err := client.UserInfo(context.Background(), accessToken)
if err != nil {
	var apiErr *goauthsdk.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("HTTP 状态码: %d\n", apiErr.Status)
		fmt.Printf("错误码: %s\n", apiErr.Code)
		fmt.Printf("错误详情: %s\n", apiErr.Detail)
	}
}
```

`APIError` 结构体字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `Status` | `int` | HTTP 状态码 |
| `Code` | `string` | 业务错误码 |
| `Detail` | `string` | 错误详情描述 |
| `Type` | `string` | RFC7807 问题类型 URI |
| `Title` | `string` | RFC7807 错误标题 |

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

仓库自带一个用于开发/测试的手工验证服务：`cmd/goauthsdk-testserver`，包含完整流程的路由。

```bash
go run ./cmd/goauthsdk-testserver
```

启动后访问 `http://localhost:7000/` 查看说明，支持的路由包括：

| 路由 | 说明 |
|------|------|
| `GET /` | 说明页 |
| `GET /auth` | 发起 OAuth 授权（可选参数: `?scope=read&state=test`） |
| `GET /callback` | OAuth 回调地址（自动接收 code 并交换 token） |
| `GET /client_credentials` | 客户端凭证模式获取令牌（可选参数: `?scope=api`） |
| `GET /introspect` | 内省令牌（必需: `?token=xxx`，可选: `&token_type_hint=access_token\|refresh_token`） |
| `GET /refresh` | 刷新访问令牌（必需参数: `?refresh_token=xxx`） |
| `GET /revoke` | 撤销令牌（必需: `?token=xxx`，可选: `&token_type_hint=access_token\|refresh_token`） |
| `GET /userinfo` | 获取用户信息（必需: `?token=xxx`） |
| `GET /user` | 获取用户详情（必需: `?token=xxx&user_id=123`） |
| `GET /parse` | 离线解析令牌（必需: `?token=xxx`，可选: `&type=access\|refresh`） |
| `GET /validate` | 离线验证令牌有效性（必需: `?token=xxx`） |

## 离线验证令牌（可选）

如果你的应用需要在本地验证 JWT 令牌（无需调用服务端 introspect 接口），可以在初始化时配置访问令牌和刷新令牌的签名密钥：

### 通过 Client 使用

```go
client, err := goauthsdk.NewClient(goauthsdk.Config{
	FrontendBaseURL:    "https://portal.example.com",
	BackendBaseURL:     "https://auth.example.com",
	ClientID:           "your-client-id",
	ClientSecret:       "your-client-secret",
	RedirectURI:        "https://yourapp.com/callback",
	AccessTokenSecret:  "your-access-token-secret",  // 访问令牌签名密钥
	RefreshTokenSecret: "your-refresh-token-secret", // 刷新令牌签名密钥
})
if err != nil {
	log.Fatal(err)
}

// 离线解析访问令牌
claims, err := client.ParseAccessToken(accessToken)
if err != nil {
	// 令牌无效、过期、签名错误等
	log.Fatal(err)
}

fmt.Printf("用户标识(Subject): %s\n", claims.Subject)
fmt.Printf("令牌类型: %s\n", claims.TokenType)
fmt.Printf("扩展字段: %v\n", claims.Extra)
fmt.Printf("过期时间: %v\n", claims.ExpiresAt)

// 离线解析刷新令牌
refreshClaims, err := client.ParseRefreshToken(refreshToken)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("刷新令牌用户标识: %s\n", refreshClaims.Subject)

// 仅验证令牌有效性（不获取 claims）
if err := client.ValidateToken(accessToken); err != nil {
	log.Fatal("令牌无效:", err)
}
```

### 独立使用 JWTVerifier

如果你只需要离线验签功能，无需完整的 OAuth 客户端，可以直接使用 `JWTVerifier`：

```go
// 创建独立的 JWTVerifier
verifier, err := goauthsdk.NewJWTVerifier(
	"your-access-token-secret",  // 访问令牌签名密钥
	"your-refresh-token-secret", // 刷新令牌签名密钥（可传空字符串）
)
if err != nil {
	log.Fatal(err)
}

// 解析访问令牌
claims, err := verifier.ParseAccessToken(accessToken)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Subject: %s\n", claims.Subject)

// 解析刷新令牌
refreshClaims, err := verifier.ParseRefreshToken(refreshToken)
if err != nil {
	log.Fatal(err)
}

// 仅验证令牌有效性
if err := verifier.ValidateToken(accessToken); err != nil {
	log.Fatal("令牌无效:", err)
}
```

也可以通过 `Client.JWTVerifier()` 获取底层的 `JWTVerifier` 实例：

```go
verifier := client.JWTVerifier()
if verifier != nil {
	claims, _ := verifier.ParseAccessToken(accessToken)
	// ...
}
```

**Claims 结构体字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `TokenType` | `TokenType` | 令牌类型：`access` 或 `refresh` |
| `Extra` | `map[string]any` | 自定义扩展字段（角色、权限等） |
| `Subject` | `string` | 用户唯一标识（来自 RegisteredClaims） |
| `Issuer` | `string` | 令牌签发者 |
| `ExpiresAt` | `*jwt.NumericDate` | 过期时间 |
| `IssuedAt` | `*jwt.NumericDate` | 签发时间 |

> **注意**：
> - `AccessTokenSecret` 必须与 goauth 服务端配置的访问令牌签名密钥一致
> - `RefreshTokenSecret` 必须与 goauth 服务端配置的刷新令牌签名密钥一致
> - 两个密钥可以只配置其中一个，但对应的解析方法需要配置相应的密钥才能使用
> - 若未配置密钥调用解析方法，将返回 `ErrJWTNotConfigured` 错误

## 常见注意事项

- **state 建议必传**：用于防止 CSRF，回调时校验 `state` 是否与发起时一致。
- **RedirectURI 必须完全一致**：需与 goauth 后台注册的回调地址匹配。
- **生产环境务必使用 HTTPS**：避免 code/token 在传输过程中被窃取。
- **不要在日志中打印完整 token/secret**：如需排查，建议仅打印前缀（仓库测试代码已做了前缀截断示例）。
- **离线验签需要正确的密钥**：`AccessTokenSecret` 和 `RefreshTokenSecret` 必须分别与服务端签发对应令牌使用的密钥一致。
