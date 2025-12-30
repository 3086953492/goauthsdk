package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// startServer 启动测试 HTTP 服务
func startServer(addr string) error {
	r := gin.Default()

	// 首页 - 说明文档
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "goauthsdk 手工测试服务",
			"config": gin.H{
				"frontend_base_url": testFrontendBaseURL,
				"backend_base_url":  testBackendBaseURL,
				"client_id":         testClientID,
				"redirect_uri":      testRedirectURI,
			},
			"routes": gin.H{
				"GET /":                   "本说明页",
				"GET /auth":               "发起 OAuth 授权（可选参数: ?scope=read&state=test）",
				"GET /callback":           "OAuth 回调地址（自动接收 code 并交换 token）",
				"GET /client_credentials": "客户端凭证模式获取令牌（可选参数: ?scope=api）",
				"GET /introspect":         "内省令牌（必需: ?token=xxx，可选: &token_type_hint=access_token|refresh_token）",
				"GET /refresh":            "刷新访问令牌（必需参数: ?refresh_token=xxx）",
				"GET /revoke":             "撤销令牌（必需: ?token=xxx，可选: &token_type_hint=access_token|refresh_token）",
				"GET /userinfo":           "获取用户信息（必需: ?token=xxx）",
				"GET /user":               "获取用户详情（必需: ?token=xxx&user_id=123）",
				"GET /parse":              "离线解析令牌（必需: ?token=xxx，可选: &type=access|refresh）",
				"GET /validate":           "离线验证令牌有效性（必需: ?token=xxx）",
			},
			"usage": []string{
				"1. 访问 /auth 发起授权",
				"2. 在 OAuth 授权页面确认授权",
				"3. 自动跳转回 /callback 并显示访问令牌",
			},
		})
	})

	// 发起授权
	r.GET("/auth", handleAuth)

	// 授权回调
	r.GET("/callback", handleCallback)

	// 客户端凭证模式
	r.GET("/client_credentials", handleClientCredentials)

	// 内省访问令牌
	r.GET("/introspect", handleIntrospect)

	// 刷新访问令牌
	r.GET("/refresh", handleRefresh)

	// 撤销令牌
	r.GET("/revoke", handleRevoke)

	// 获取用户信息
	r.GET("/userinfo", handleUserInfo)

	// 获取用户详情
	r.GET("/user", handleGetUser)

	// 离线解析令牌
	r.GET("/parse", handleParse)

	// 离线验证令牌
	r.GET("/validate", handleValidate)

	return r.Run(addr)
}
