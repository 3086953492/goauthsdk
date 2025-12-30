package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

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

// handleClientCredentials 处理客户端凭证模式请求
func handleClientCredentials(c *gin.Context) {
	// 从 query 读取可选参数
	scope := c.DefaultQuery("scope", "")

	// 创建客户端
	client, err := newTestClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "创建客户端失败",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("开始客户端凭证模式获取令牌: scope=%s", scope)

	// 调用客户端凭证模式接口
	token, err := client.ClientCredentialsToken(context.Background(), scope)
	if err != nil {
		log.Printf("客户端凭证模式获取令牌失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "客户端凭证模式获取令牌失败",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("成功获取客户端凭证令牌: %s", token.AccessToken[:16]+"...")

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "客户端凭证模式获取令牌成功",
		"token": gin.H{
			"access_token": token.AccessToken,
			"expires_in":   token.ExpiresIn,
			"token_type":   token.TokenType,
			"scope":        token.Scope,
		},
		"note": "该 token 的 JWT sub 为 client:<client_id>，不适用于 /userinfo 接口",
	})
}

// handleIntrospect 处理内省请求（RFC 7662）
func handleIntrospect(c *gin.Context) {
	// 读取 token 参数
	token := c.Query("token")
	tokenTypeHint := c.Query("token_type_hint")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少 token 参数",
			"detail": "请提供 token 参数，例如: /introspect?token=xxx 或 /introspect?token=xxx&token_type_hint=refresh_token",
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
	log.Printf("开始内省令牌: %s (hint: %s)", tokenPreview, tokenTypeHint)

	// 调用内省接口（支持 token_type_hint）
	resp, err := client.IntrospectTokenWithHint(context.Background(), token, tokenTypeHint)
	if err != nil {
		log.Printf("内省令牌失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "内省令牌失败",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("内省成功: active=%v", resp.Active)

	// 返回 RFC 7662 格式的内省结果
	c.JSON(http.StatusOK, gin.H{
		"active":     resp.Active,
		"scope":      resp.Scope,
		"client_id":  resp.ClientID,
		"username":   resp.Username,
		"token_type": resp.TokenType,
		"exp":        resp.Exp,
		"sub":        resp.Sub,
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

// handleRevoke 处理撤销令牌请求（RFC 7009）
func handleRevoke(c *gin.Context) {
	// 读取参数
	token := c.Query("token")
	tokenTypeHint := c.Query("token_type_hint")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少 token 参数",
			"detail": "请提供 token 参数，例如: /revoke?token=xxx 或 /revoke?token=xxx&token_type_hint=refresh_token",
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
	log.Printf("开始撤销令牌: %s (hint: %s)", tokenPreview, tokenTypeHint)

	// 调用撤销接口（支持 token_type_hint）
	err = client.RevokeTokenWithHint(context.Background(), token, tokenTypeHint)
	if err != nil {
		log.Printf("撤销令牌失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "撤销令牌失败",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("撤销令牌成功")

	// 按 RFC 7009，成功返回 HTTP 200
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "令牌撤销成功",
	})
}

// handleUserInfo 处理获取用户信息请求
func handleUserInfo(c *gin.Context) {
	// 读取 token 参数
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少 token 参数",
			"detail": "请提供 token 参数，例如: /userinfo?token=xxx",
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
	log.Printf("开始获取用户信息: %s", tokenPreview)

	// 调用用户信息接口
	info, err := client.UserInfo(context.Background(), token)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)

		// 尝试断言为 ProblemDetails 以获取详细错误
		if pd, ok := err.(*goauthsdk.ProblemDetails); ok {
			c.JSON(pd.Status, gin.H{
				"error":  pd.Code,
				"detail": pd.Detail,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "获取用户信息失败",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("获取用户信息成功: sub=%s, nickname=%s", info.Sub, info.Nickname)

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取用户信息成功",
		"data": gin.H{
			"sub":        info.Sub,
			"nickname":   info.Nickname,
			"picture":    info.Picture,
			"updated_at": info.UpdatedAt,
		},
	})
}

// handleGetUser 处理获取用户详情请求
func handleGetUser(c *gin.Context) {
	// 读取参数
	token := c.Query("token")
	userIDStr := c.Query("user_id")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少 token 参数",
			"detail": "请提供 token 参数，例如: /user?token=xxx&user_id=123",
		})
		return
	}

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少 user_id 参数",
			"detail": "请提供 user_id 参数，例如: /user?token=xxx&user_id=123",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "user_id 参数无效",
			"detail": "user_id 必须为正整数",
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
	log.Printf("开始获取用户详情: token=%s, user_id=%d", tokenPreview, userID)

	// 调用获取用户详情接口
	user, err := client.GetUser(context.Background(), token, userID)
	if err != nil {
		log.Printf("获取用户详情失败: %v", err)

		// 尝试断言为 ProblemDetails 以获取详细错误
		if pd, ok := err.(*goauthsdk.ProblemDetails); ok {
			errorCode := pd.Code
			if errorCode == "" {
				errorCode = pd.Title
			}
			c.JSON(pd.Status, gin.H{
				"error":  errorCode,
				"detail": pd.Detail,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "获取用户详情失败",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("获取用户详情成功: id=%d, username=%s, nickname=%s", user.ID, user.Username, user.Nickname)

	// 返回用户详情
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取用户详情成功",
		"data": gin.H{
			"id":         user.ID,
			"subject":    user.Subject,
			"username":   user.Username,
			"nickname":   user.Nickname,
			"avatar":     user.Avatar,
			"status":     user.Status,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// handleParse 处理离线解析令牌请求
func handleParse(c *gin.Context) {
	// 读取参数
	token := c.Query("token")
	tokenType := c.DefaultQuery("type", "access")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少 token 参数",
			"detail": "请提供 token 参数，例如: /parse?token=xxx 或 /parse?token=xxx&type=refresh",
		})
		return
	}

	// 创建支持 JWT 的客户端
	client, err := newTestClientWithJWT()
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
	log.Printf("开始离线解析令牌: %s (type: %s)", tokenPreview, tokenType)

	// 根据类型调用不同的解析方法
	if tokenType == "refresh" {
		claims, err := client.ParseRefreshToken(token)
		if err != nil {
			log.Printf("离线解析刷新令牌失败: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "解析刷新令牌失败",
				"detail": err.Error(),
			})
			return
		}

		log.Printf("离线解析刷新令牌成功: subject=%s", claims.Subject)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "离线解析刷新令牌成功",
			"claims": gin.H{
				"token_type": claims.TokenType,
				"extra":      claims.Extra,
				"issuer":     claims.Issuer,
				"subject":    claims.Subject,
				"expires_at": claims.ExpiresAt,
				"issued_at":  claims.IssuedAt,
			},
		})
		return
	}

	// 默认解析访问令牌
	claims, err := client.ParseAccessToken(token)
	if err != nil {
		log.Printf("离线解析访问令牌失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "解析访问令牌失败",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("离线解析访问令牌成功: subject=%s", claims.Subject)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "离线解析访问令牌成功",
		"claims": gin.H{
			"token_type": claims.TokenType,
			"extra":      claims.Extra,
			"issuer":     claims.Issuer,
			"subject":    claims.Subject,
			"expires_at": claims.ExpiresAt,
			"issued_at":  claims.IssuedAt,
		},
	})
}

// handleValidate 处理离线验证令牌请求
func handleValidate(c *gin.Context) {
	// 读取参数
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少 token 参数",
			"detail": "请提供 token 参数，例如: /validate?token=xxx",
		})
		return
	}

	// 创建支持 JWT 的客户端
	client, err := newTestClientWithJWT()
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
	log.Printf("开始离线验证令牌: %s", tokenPreview)

	// 验证令牌
	err = client.ValidateToken(token)
	if err != nil {
		log.Printf("离线验证令牌失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "令牌验证失败",
			"detail":  err.Error(),
		})
		return
	}

	log.Printf("离线验证令牌成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "令牌验证成功，令牌有效",
	})
}
