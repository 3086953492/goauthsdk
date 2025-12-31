package goauthsdk

// AccessTokenInfo 访问令牌信息
type AccessTokenInfo struct {
	AccessToken string `json:"access_token"` // JWT 访问令牌
	ExpiresIn   int    `json:"expires_in"`   // 过期时间（秒）
}

// RefreshTokenInfo 刷新令牌信息
type RefreshTokenInfo struct {
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int    `json:"expires_in"`    // 过期时间（秒）
}

// TokenResponse 是访问令牌响应（authorization_code / refresh_token 模式）
// 与后端 dto.OAuthAccessTokenResponse 字段对齐
type TokenResponse struct {
	AccessToken  AccessTokenInfo  `json:"access_token"`  // 访问令牌信息
	RefreshToken RefreshTokenInfo `json:"refresh_token"` // 刷新令牌信息
	TokenType    string           `json:"token_type"`    // 令牌类型，通常为 "Bearer"
	Scope        string           `json:"scope"`         // 授权范围
}

// ClientCredentialsTokenResponse 是客户端凭证模式（client_credentials）的访问令牌响应
// 该模式无用户上下文，因此无 refresh_token
type ClientCredentialsTokenResponse struct {
	AccessToken string `json:"access_token"` // JWT 访问令牌
	ExpiresIn   int    `json:"expires_in"`   // 过期时间（秒）
	TokenType   string `json:"token_type"`   // 令牌类型，通常为 "Bearer"
	Scope       string `json:"scope"`        // 授权范围；为空时返回空字符串
}
