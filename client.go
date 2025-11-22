package goauthsdk

import (
	"net/http"
)

// Client 是 goauth SDK 的客户端
// 封装了 OAuth 授权码模式的常用操作
type Client struct {
	cfg Config
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

	return &Client{cfg: cfg}, nil
}
