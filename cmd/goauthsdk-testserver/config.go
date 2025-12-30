package main

const (
	// testFrontendBaseURL OAuth 前端站点地址
	testFrontendBaseURL = "http://localhost:5173"

	// testBackendBaseURL OAuth 后端服务地址
	testBackendBaseURL = "http://localhost:9000"

	// testClientID OAuth 客户端 ID
	testClientID = "1"

	// testClientSecret OAuth 客户端密钥
	testClientSecret = "mC9dvSBXPIIDLWP2MSauuxybZmICfNpq"

	// testRedirectURI OAuth 回调地址，需与客户端注册的回调地址一致
	testRedirectURI = "http://localhost:7000/callback"

	// testAccessTokenSecret 访问令牌签名密钥，用于离线验证访问令牌（需与 goauth 服务端配置一致）
	testAccessTokenSecret = "GO4ymlqBMkucpQ60roh17ZADPcY8outx"

	// testRefreshTokenSecret 刷新令牌签名密钥，用于离线验证刷新令牌（需与 goauth 服务端配置一致）
	testRefreshTokenSecret = "tnwBPejxaajp3m1AzLMAs9viS4GLGoLj"

	// serverAddr 测试服务监听地址
	serverAddr = ":7000"
)
