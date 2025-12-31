package main

import (
	"github.com/3086953492/goauthsdk"
)

// newTestClient 创建测试客户端
func newTestClient() (*goauthsdk.Client, error) {
	return goauthsdk.NewClient(
		testFrontendBaseURL,
		testBackendBaseURL,
		testClientID,
		testClientSecret,
		testRedirectURI,
	)
}

// newTestClientWithJWT 创建支持离线验签的测试客户端
func newTestClientWithJWT() (*goauthsdk.Client, error) {
	return goauthsdk.NewClient(
		testFrontendBaseURL,
		testBackendBaseURL,
		testClientID,
		testClientSecret,
		testRedirectURI,
		goauthsdk.WithJWTSecrets(testAccessTokenSecret, testRefreshTokenSecret),
	)
}
