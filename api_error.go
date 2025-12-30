package goauthsdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// APIError 是 SDK 统一的 API 错误类型
// 调用方可以通过 errors.As(err, &apiErr) 获取结构化错误信息
type APIError struct {
	// Status HTTP 状态码
	Status int `json:"status"`

	// Code 业务错误码（字符串形式）
	// 可能来源于 RFC7807 ProblemDetails.Code、ProblemDetails.Title 或 apiCodeResponse.Code
	Code string `json:"code,omitempty"`

	// Detail 错误详情描述
	Detail string `json:"detail,omitempty"`

	// Type RFC7807 问题类型 URI（通常为 "about:blank"）
	Type string `json:"type,omitempty"`

	// Title RFC7807 错误标题
	Title string `json:"title,omitempty"`
}

// Error 实现 error 接口
func (e *APIError) Error() string {
	code := e.Code
	if code == "" {
		code = e.Title
	}
	if code == "" {
		code = http.StatusText(e.Status)
	}
	if code == "" {
		code = strconv.Itoa(e.Status)
	}

	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", code, e.Detail)
	}
	return code
}

// decodeAPIError 从 HTTP 响应解析统一的 APIError
// 优先按 RFC7807 ProblemDetails 解码，其次尝试 {code, message}，最后兜底生成基于 HTTP status 的错误
func decodeAPIError(resp *http.Response, body []byte) *APIError {
	// 尝试解析 RFC7807 ProblemDetails
	var pd ProblemDetails
	if err := json.Unmarshal(body, &pd); err == nil && (pd.Code != "" || pd.Title != "" || pd.Detail != "") {
		code := pd.Code
		if code == "" {
			code = pd.Title
		}
		return &APIError{
			Status: resp.StatusCode,
			Code:   code,
			Detail: pd.Detail,
			Type:   pd.Type,
			Title:  pd.Title,
		}
	}

	// 尝试解析 {code, message} 格式
	var codeMsg struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &codeMsg); err == nil && (codeMsg.Code != 0 || codeMsg.Message != "") {
		return &APIError{
			Status: resp.StatusCode,
			Code:   strconv.Itoa(codeMsg.Code),
			Detail: codeMsg.Message,
		}
	}

	// 兜底：基于 HTTP 状态码生成错误
	return &APIError{
		Status: resp.StatusCode,
		Code:   http.StatusText(resp.StatusCode),
		Detail: fmt.Sprintf("request failed with HTTP %d", resp.StatusCode),
	}
}

// newBusinessError 创建业务失败的 APIError（HTTP 2xx 但 code != 0）
func newBusinessError(httpStatus int, bizCode int, message string) *APIError {
	return &APIError{
		Status: httpStatus,
		Code:   strconv.Itoa(bizCode),
		Detail: message,
	}
}
