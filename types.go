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

// apiResponse 是后端 API 的通用响应结构
// 参考前端 request.ts 对 success / message / data 的处理方式
type apiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
