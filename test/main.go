package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/3086953492/goauthsdk"
)

// ============================================================================
// goauthsdk 手工测试服务 - 用于开发/测试阶段手动验证 OAuth 流程
// ============================================================================

const (
	// testFrontendBaseURL OAuth 前端站点地址
	testFrontendBaseURL = "http://localhost:5173"

	// testBackendBaseURL OAuth 后端服务地址
	testBackendBaseURL = "http://localhost:9000"

	// testClientID OAuth 客户端 ID
	testClientID = "1"

	// testClientSecret OAuth 客户端密钥
	testClientSecret = "kwNC6fyfds303wOkqBtPNyMY03xSswbY"

	// testRedirectURI OAuth 回调地址，需与客户端注册的回调地址一致
	testRedirectURI = "http://localhost:7000/callback"

	// serverAddr 测试服务监听地址
	serverAddr = ":7000"
)

func main() {
	log.Printf("启动 goauthsdk 测试服务于 %s", serverAddr)
	log.Printf("配置信息:")
	log.Printf("  - 前端地址: %s", testFrontendBaseURL)
	log.Printf("  - 后端地址: %s", testBackendBaseURL)
	log.Printf("  - 客户端ID: %s", testClientID)
	log.Printf("  - 回调地址: %s", testRedirectURI)
	log.Printf("\n访问 http://localhost%s/ 查看使用说明\n", serverAddr)

	if err := startTestServer(serverAddr); err != nil {
		log.Fatal(err)
	}
}

// startTestServer 启动测试 HTTP 服务
func startTestServer(addr string) error {
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
				"GET /":         "本说明页",
				"GET /auth":     "发起 OAuth 授权（可选参数: ?scope=read&state=test）",
				"GET /callback": "OAuth 回调地址（自动接收 code 并交换 token）",
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

	return r.Run(addr)
}

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

// handleAuth 处理授权发起请求
func handleAuth(c *gin.Context) {
	// 创建客户端
	client, err := newTestClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "创建客户端失败",
			"detail": err.Error(),
		})
		return
	}

	// 从 query 读取可选参数
	scope := c.DefaultQuery("scope", "profile")
	state := c.DefaultQuery("state", "manual-test")

	// 构建授权 URL
	authURL, err := client.BuildAuthorizationURL(state, scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "构建授权 URL 失败",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("发起授权请求: scope=%s, state=%s", scope, state)
	log.Printf("重定向到: %s", authURL)

	// 重定向到授权页面
	c.Redirect(http.StatusFound, authURL)
}

// handleCallback 处理授权回调请求
func handleCallback(c *gin.Context) {
	// 读取授权码
	code := c.Query("code")
	state := c.Query("state")

	log.Printf("收到回调请求: code=%s, state=%s", code, state)

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少授权码",
			"detail": "未收到 code 参数，授权可能失败",
			"state":  state,
		})
		return
	}

	// 创建客户端
	client, err := newTestClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "创建客户端失败",
			"detail": err.Error(),
		})
		return
	}

	// 交换访问令牌
	log.Printf("开始交换访问令牌...")
	token, err := client.ExchangeToken(context.Background(), code)
	if err != nil {
		log.Printf("交换令牌失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "交换访问令牌失败",
			"detail": err.Error(),
			"code":   code,
			"state":  state,
		})
		return
	}

	log.Printf("成功获取访问令牌: %s", token.AccessToken)

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功获取访问令牌",
		"state":   state,
		"token": gin.H{
			"access_token":  token.AccessToken,
			"token_type":    token.TokenType,
			"expires_in":    token.ExpiresIn,
			"refresh_token": token.RefreshToken,
			"scope":         token.Scope,
		},
	})
}

