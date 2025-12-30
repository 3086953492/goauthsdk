package goauthsdk

import (
	"fmt"
	"net/http"

	"github.com/3086953492/gokit/jwt"
)

// Client 是 goauth SDK 的客户端
// 封装了 OAuth 授权码模式的常用操作
type Client struct {
	cfg        Config
	jwtManager *jwt.Manager
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
	// 校验必填字段
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	// 标准化 BaseURL，去掉末尾的 /
	normalizeConfig(&cfg)

	// 如果未提供 HTTPClient，使用默认客户端
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	client := &Client{cfg: cfg}

	// 若配置了 AccessTokenSecret 或 RefreshTokenSecret，初始化 jwtManager 用于离线验签
	if cfg.AccessTokenSecret != "" || cfg.RefreshTokenSecret != "" {
		jwtMgr, err := jwt.NewManager(
			jwt.WithAccessSecret(cfg.AccessTokenSecret),
			jwt.WithRefreshSecret(cfg.RefreshTokenSecret),
		)
		if err != nil {
			return nil, fmt.Errorf("create jwt manager: %w", err)
		}
		client.jwtManager = jwtMgr
	}

	return client, nil
}
