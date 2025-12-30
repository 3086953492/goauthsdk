package main

import (
	"github.com/3086953492/goauthsdk"
)

// newTestClient 创建测试客户端
func newTestClient() (*goauthsdk.Client, error) {
	return goauthsdk.NewClient(goauthsdk.Config{
		FrontendBaseURL: testFrontendBaseURL,
		BackendBaseURL:  testBackendBaseURL,
		ClientID:        testClientID,
		ClientSecret:    testClientSecret,
		RedirectURI:     testRedirectURI,
	})
}

// newTestClientWithJWT 创建支持离线验签的测试客户端
func newTestClientWithJWT() (*goauthsdk.Client, error) {
	return goauthsdk.NewClient(goauthsdk.Config{
		FrontendBaseURL:    testFrontendBaseURL,
		BackendBaseURL:     testBackendBaseURL,
		ClientID:           testClientID,
		ClientSecret:       testClientSecret,
		RedirectURI:        testRedirectURI,
		AccessTokenSecret:  testAccessTokenSecret,
		RefreshTokenSecret: testRefreshTokenSecret,
	})
}
