package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
