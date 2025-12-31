package goauthsdk

import "fmt"

// Client 是 goauth SDK 的客户端
// 封装了 OAuth 授权码模式的常用操作
type Client struct {
	cfg         Config
	jwtVerifier *JWTVerifier
}

// NewClient 创建一个新的 goauth SDK 客户端
//
// 示例用法:
//
//	client, err := goauthsdk.NewClient(goauthsdk.Config{
//	    FrontendBaseURL: "https://portal.example.com",
//	    BackendBaseURL:  "https://auth.example.com",
//	    ClientID:        "your-client-id",
//	    ClientSecret:    "your-client-secret",
//	    RedirectURI:     "https://yourapp.com/callback",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewClient(cfg Config) (*Client, error) {
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	normalizeConfig(&cfg)

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
