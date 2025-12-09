package goauthsdk

// AccessTokenInfo 访问令牌信息
type AccessTokenInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// RefreshTokenInfo 刷新令牌信息
type RefreshTokenInfo struct {
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// TokenResponse 是访问令牌响应
// 与后端 dto.OAuthAccessTokenResponse 字段对齐
type TokenResponse struct {
	AccessToken  AccessTokenInfo  `json:"access_token"`
	RefreshToken RefreshTokenInfo `json:"refresh_token"`
	TokenType    string           `json:"token_type"`
	Scope        string           `json:"scope"`
}

// IntrospectionResponse 内省响应结构体 (RFC 7662)
// 用于验证访问令牌的有效性和获取令牌相关信息
type IntrospectionResponse struct {
	Active    bool   `json:"active"`               // 令牌是否有效
	Scope     string `json:"scope,omitempty"`      // 令牌授权范围
	ClientID  string `json:"client_id,omitempty"`  // 客户端 ID
	Username  string `json:"username,omitempty"`   // 用户名
	TokenType string `json:"token_type,omitempty"` // 令牌类型
	Exp       int64  `json:"exp,omitempty"`        // 过期时间戳
	Sub       string `json:"sub,omitempty"`        // 主体标识
}

// apiResponse 是后端 API 的通用响应结构
// 参考前端 request.ts 对 success / message / data 的处理方式
type apiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
