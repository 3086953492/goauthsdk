package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
