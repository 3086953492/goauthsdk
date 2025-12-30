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

// TokenResponse 是访问令牌响应（authorization_code / refresh_token 模式）
// 与后端 dto.OAuthAccessTokenResponse 字段对齐
type TokenResponse struct {
	AccessToken  AccessTokenInfo  `json:"access_token"`
	RefreshToken RefreshTokenInfo `json:"refresh_token"`
	TokenType    string           `json:"token_type"`
	Scope        string           `json:"scope"`
}

// ClientCredentialsTokenResponse 是客户端凭证模式（client_credentials）的访问令牌响应
// 该模式无用户上下文，因此无 refresh_token
type ClientCredentialsTokenResponse struct {
	AccessToken string `json:"access_token"` // JWT 访问令牌
	ExpiresIn   int    `json:"expires_in"`   // 过期时间（秒）
	TokenType   string `json:"token_type"`   // 令牌类型，通常为 "Bearer"
	Scope       string `json:"scope"`        // 授权范围；为空时返回空字符串
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

// UserInfo 用户信息结构体
// 用于 GET /api/v1/oauth/userinfo 接口返回的用户信息
type UserInfo struct {
	Sub       string `json:"sub"`        // 用户唯一标识（用户ID）
	Nickname  string `json:"nickname"`   // 用户昵称
	Picture   string `json:"picture"`    // 用户头像URL
	UpdatedAt int64  `json:"updated_at"` // 用户信息更新时间（Unix 时间戳）
}

// ProblemDetails 是 RFC 7807 Problem Details 风格的错误响应结构
// 用于解析 401/403/404 等错误响应
type ProblemDetails struct {
	Type   string `json:"type"`            // 问题类型 URI（通常为 "about:blank"）
	Title  string `json:"title,omitempty"` // 错误标题（如 UNAUTHORIZED、FORBIDDEN、USER_NOT_FOUND）
	Status int    `json:"status"`          // HTTP 状态码
	Code   string `json:"code,omitempty"`  // 业务错误码（如 INVALID_TOKEN、INSUFFICIENT_SCOPE）
	Detail string `json:"detail"`          // 错误详情描述
}

// Error 实现 error 接口，使 ProblemDetails 可作为 error 返回
func (p *ProblemDetails) Error() string {
	code := p.Code
	if code == "" {
		code = p.Title
	}
	if code == "" {
		code = "UNKNOWN_ERROR"
	}
	return code + ": " + p.Detail
}

// apiCodeResponse 是后端 API 的通用响应结构
// 接口格式：{ "code": 0, "message": "...", "data": {...} }（code == 0 表示成功）
type apiCodeResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
