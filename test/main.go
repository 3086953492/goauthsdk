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
				"GET /":           "本说明页",
				"GET /auth":       "发起 OAuth 授权（可选参数: ?scope=read&state=test）",
				"GET /callback":   "OAuth 回调地址（自动接收 code 并交换 token）",
				"GET /introspect": "内省访问令牌（必需参数: ?token=xxx）",
				"GET /refresh":    "刷新访问令牌（必需参数: ?refresh_token=xxx）",
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

	// 内省访问令牌
	r.GET("/introspect", handleIntrospect)

	// 刷新访问令牌
	r.GET("/refresh", handleRefresh)

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

	log.Printf("成功获取访问令牌: %s", token.AccessToken.AccessToken)

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功获取访问令牌",
		"state":   state,
		"token": gin.H{
			"access_token":             token.AccessToken.AccessToken,
			"access_token_expires_in":  token.AccessToken.ExpiresIn,
			"refresh_token":            token.RefreshToken.RefreshToken,
			"refresh_token_expires_in": token.RefreshToken.ExpiresIn,
			"token_type":               token.TokenType,
			"scope":                    token.Scope,
		},
	})
}

// handleIntrospect 处理内省请求
func handleIntrospect(c *gin.Context) {
	// 读取 token 参数
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少 token 参数",
			"detail": "请提供 token 参数，例如: /introspect?token=xxx",
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

	// 打印日志（只显示 token 前 16 个字符以保护敏感信息）
	tokenPreview := token
	if len(token) > 16 {
		tokenPreview = token[:16] + "..."
	}
	log.Printf("开始内省访问令牌: %s", tokenPreview)

	// 调用内省接口
	resp, err := client.IntrospectToken(context.Background(), token)
	if err != nil {
		log.Printf("内省令牌失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "内省访问令牌失败",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("内省成功: active=%v", resp.Active)

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "内省成功",
		"introspection": gin.H{
			"active":     resp.Active,
			"scope":      resp.Scope,
			"client_id":  resp.ClientID,
			"username":   resp.Username,
			"token_type": resp.TokenType,
			"exp":        resp.Exp,
			"sub":        resp.Sub,
		},
	})
}

// handleRefresh 处理刷新令牌请求
func handleRefresh(c *gin.Context) {
	// 读取 refresh_token 参数
	refreshToken := c.Query("refresh_token")

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "缺少 refresh_token 参数",
			"detail":  "请提供 refresh_token 参数，例如: /refresh?refresh_token=xxx",
		})
		return
	}

	// 创建客户端
	client, err := newTestClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建客户端失败",
			"detail":  err.Error(),
		})
		return
	}

	// 打印日志（只显示 refresh_token 前 16 个字符以保护敏感信息）
	tokenPreview := refreshToken
	if len(refreshToken) > 16 {
		tokenPreview = refreshToken[:16] + "..."
	}
	log.Printf("开始刷新访问令牌: %s", tokenPreview)

	// 调用刷新令牌接口
	token, err := client.RefreshToken(context.Background(), refreshToken)
	if err != nil {
		log.Printf("刷新令牌失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "刷新访问令牌失败",
			"detail":  err.Error(),
		})
		return
	}

	log.Printf("成功刷新访问令牌: %s", token.AccessToken.AccessToken)

	// 返回精简结果
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "刷新令牌成功",
		"token": gin.H{
			"access_token":             token.AccessToken.AccessToken,
			"access_token_expires_in":  token.AccessToken.ExpiresIn,
			"refresh_token":            token.RefreshToken.RefreshToken,
			"refresh_token_expires_in": token.RefreshToken.ExpiresIn,
			"token_type":               token.TokenType,
			"scope":                    token.Scope,
		},
	})
}
