package goauthsdk

import "net/http"

// HTTPDoer 是发送 HTTP 请求的最小接口
// *http.Client 自动满足此接口
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
