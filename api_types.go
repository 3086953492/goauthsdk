package goauthsdk

// apiCodeResponse 是后端 API 的通用响应结构
// 接口格式：{ "code": 0, "message": "...", "data": {...} }（code == 0 表示成功）
type apiCodeResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// problemDetails 是 RFC 7807 Problem Details 风格的错误响应结构（内部使用）
// 用于解析 401/403/404 等错误响应，对外统一返回 *APIError
type problemDetails struct {
	Type   string `json:"type"`            // 问题类型 URI（通常为 "about:blank"）
	Title  string `json:"title,omitempty"` // 错误标题（如 UNAUTHORIZED、FORBIDDEN、USER_NOT_FOUND）
	Status int    `json:"status"`          // HTTP 状态码
	Code   string `json:"code,omitempty"`  // 业务错误码（如 INVALID_TOKEN、INSUFFICIENT_SCOPE）
	Detail string `json:"detail"`          // 错误详情描述
}
