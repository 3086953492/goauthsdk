package goauthsdk

// IntrospectionResponse 内省响应结构体 (RFC 7662)
// 用于验证访问令牌的有效性和获取令牌相关信息
type IntrospectionResponse struct {
	Active    bool   `json:"active"`               // 令牌是否有效
	Scope     string `json:"scope,omitempty"`      // 令牌授权范围
	ClientID  string `json:"client_id,omitempty"`  // 客户端 ID
	Username  string `json:"username,omitempty"`   // 用户名
	TokenType string `json:"token_type,omitempty"` // 令牌类型
	Exp       int64  `json:"exp,omitempty"`        // 过期时间戳（Unix 时间戳，秒）
	Sub       string `json:"sub,omitempty"`        // 主体标识
}
