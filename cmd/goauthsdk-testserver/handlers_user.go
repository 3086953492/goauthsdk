package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/3086953492/goauthsdk"
	"github.com/gin-gonic/gin"
)

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

		// 尝试断言为 APIError 以获取详细错误
		var apiErr *goauthsdk.APIError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.Status, gin.H{
				"error":  apiErr.Code,
				"detail": apiErr.Detail,
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
	sub := c.Query("sub")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少 token 参数",
			"detail": "请提供 token 参数，例如: /user?token=xxx&sub=user-sub-xxx",
		})
		return
	}

	if sub == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "缺少 sub 参数",
			"detail": "请提供 sub 参数，例如: /user?token=xxx&sub=user-sub-xxx",
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
	log.Printf("开始获取用户详情: token=%s, sub=%s", tokenPreview, sub)

	// 调用获取用户详情接口
	user, err := client.GetUser(context.Background(), token, sub)
	if err != nil {
		log.Printf("获取用户详情失败: %v", err)

		// 尝试断言为 APIError 以获取详细错误
		var apiErr *goauthsdk.APIError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.Status, gin.H{
				"error":  apiErr.Code,
				"detail": apiErr.Detail,
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
