package goauthsdk

// Config 是 goauth SDK 的配置结构
// 用于初始化 OAuth 客户端
type Config struct {
	// FrontendBaseURL 前端站点基础地址，例如 https://portal.example.com
	// SDK 将在其上拼接 /oauth/authorize 构建用户授权确认页 URL
	FrontendBaseURL string

	// BackendBaseURL goauth 后端服务基础地址，例如 https://auth.example.com
	// SDK 内部拼接 /api/v1/oauth/... 调用真正的 OAuth 接口
	BackendBaseURL string

	// ClientID OAuth 客户端 ID
	ClientID string

	// ClientSecret OAuth 客户端密钥
	ClientSecret string

	// RedirectURI OAuth 回调地址，必须在客户端注册的回调白名单中
	RedirectURI string

	// HTTPClient 可选的 HTTP 客户端，需满足 HTTPDoer 接口
	// 不传则使用 http.DefaultClient
	HTTPClient HTTPDoer

	// AccessTokenSecret 可选的访问令牌签名密钥，用于离线验证访问令牌
	// 若配置此字段，可使用 ParseAccessToken 方法进行本地验签
	AccessTokenSecret string

	// RefreshTokenSecret 可选的刷新令牌签名密钥，用于离线验证刷新令牌
	// 若配置此字段，可使用 ParseRefreshToken 方法进行本地验签
	RefreshTokenSecret string
}
