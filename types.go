package goauthsdk

// TokenResponse 是访问令牌响应
// 与后端 dto.OAuthAccessTokenResponse 字段对齐
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// apiResponse 是后端 API 的通用响应结构
// 参考前端 request.ts 对 success / message / data 的处理方式
type apiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
