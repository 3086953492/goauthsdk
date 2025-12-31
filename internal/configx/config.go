package configx

import (
	"github.com/3086953492/goauthsdk/internal/httpx"
)

// Config 是 goauth SDK 的内部配置结构
type Config struct {
	// FrontendBaseURL 前端站点基础地址
	FrontendBaseURL string

	// BackendBaseURL goauth 后端服务基础地址
	BackendBaseURL string

	// ClientID OAuth 客户端 ID
	ClientID string

	// ClientSecret OAuth 客户端密钥
	ClientSecret string

	// RedirectURI OAuth 回调地址
	RedirectURI string

	// HTTPClient 可选的 HTTP 客户端
	HTTPClient httpx.HTTPDoer

	// AccessTokenSecret 可选的访问令牌签名密钥
	AccessTokenSecret string

	// RefreshTokenSecret 可选的刷新令牌签名密钥
	RefreshTokenSecret string
}
