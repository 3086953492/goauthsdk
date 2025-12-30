package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
