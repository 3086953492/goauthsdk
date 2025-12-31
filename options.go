package goauthsdk

import (
	"github.com/3086953492/goauthsdk/internal/configx"
	"github.com/3086953492/goauthsdk/internal/httpx"
)

// ClientOption 用于配置 Client 的可选参数
type ClientOption func(*configx.Config)

// WithHTTPClient 设置自定义 HTTP 客户端
// 传入的 client 必须满足 HTTPDoer 接口（*http.Client 自动满足）
func WithHTTPClient(client httpx.HTTPDoer) ClientOption {
	return func(cfg *configx.Config) {
		cfg.HTTPClient = client
	}
}

// WithAccessTokenSecret 设置访问令牌签名密钥，用于离线验签
func WithAccessTokenSecret(secret string) ClientOption {
	return func(cfg *configx.Config) {
		cfg.AccessTokenSecret = secret
	}
}

// WithRefreshTokenSecret 设置刷新令牌签名密钥，用于离线验签
func WithRefreshTokenSecret(secret string) ClientOption {
	return func(cfg *configx.Config) {
		cfg.RefreshTokenSecret = secret
	}
}

// WithJWTSecrets 同时设置访问令牌和刷新令牌的签名密钥
func WithJWTSecrets(accessSecret, refreshSecret string) ClientOption {
	return func(cfg *configx.Config) {
		cfg.AccessTokenSecret = accessSecret
		cfg.RefreshTokenSecret = refreshSecret
	}
}
