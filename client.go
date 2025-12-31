package goauthsdk

import (
	"fmt"

	"github.com/3086953492/goauthsdk/internal/configx"
)

// Client 是 goauth SDK 的客户端
// 封装了 OAuth 授权码模式的常用操作
type Client struct {
	cfg         configx.Config
	jwtVerifier *JWTVerifier
}

// NewClient 创建一个新的 goauth SDK 客户端
//
// 必填参数:
//   - frontendBaseURL: 前端站点基础地址，例如 https://portal.example.com
//   - backendBaseURL: goauth 后端服务基础地址，例如 https://auth.example.com
//   - clientID: OAuth 客户端 ID
//   - clientSecret: OAuth 客户端密钥
//   - redirectURI: OAuth 回调地址，必须在客户端注册的回调白名单中
//
// 可选参数通过 ClientOption 传入:
//   - WithHTTPClient: 自定义 HTTP 客户端
//   - WithAccessTokenSecret: 访问令牌签名密钥（用于离线验签）
//   - WithRefreshTokenSecret: 刷新令牌签名密钥（用于离线验签）
//   - WithJWTSecrets: 同时设置访问/刷新令牌密钥
//
// 示例用法:
//
//	client, err := goauthsdk.NewClient(
//	    "https://portal.example.com",
//	    "https://auth.example.com",
//	    "your-client-id",
//	    "your-client-secret",
//	    "https://yourapp.com/callback",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewClient(
	frontendBaseURL, backendBaseURL, clientID, clientSecret, redirectURI string,
	opts ...ClientOption,
) (*Client, error) {
	cfg := configx.Config{
		FrontendBaseURL: frontendBaseURL,
		BackendBaseURL:  backendBaseURL,
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		RedirectURI:     redirectURI,
	}

	// 应用可选配置
	for _, opt := range opts {
		opt(&cfg)
	}

	if err := configx.Validate(&cfg); err != nil {
		return nil, err
	}
	configx.Normalize(&cfg)

	client := &Client{cfg: cfg}

	// 若配置了 AccessTokenSecret 或 RefreshTokenSecret，创建 JWTVerifier 用于离线验签
	if cfg.AccessTokenSecret != "" || cfg.RefreshTokenSecret != "" {
		verifier, err := NewJWTVerifier(cfg.AccessTokenSecret, cfg.RefreshTokenSecret)
		if err != nil {
			return nil, fmt.Errorf("create jwt verifier: %w", err)
		}
		client.jwtVerifier = verifier
	}

	return client, nil
}

// JWTVerifier 返回 Client 持有的 JWTVerifier 实例
// 若初始化时未配置 AccessTokenSecret 或 RefreshTokenSecret，返回 nil
func (c *Client) JWTVerifier() *JWTVerifier {
	return c.jwtVerifier
}
