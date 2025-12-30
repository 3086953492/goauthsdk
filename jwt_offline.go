package goauthsdk

import (
	"errors"
	"fmt"

	"github.com/3086953492/gokit/jwt"
)

// ErrJWTNotConfigured 表示未配置 JWT 密钥，无法进行离线验签
var ErrJWTNotConfigured = errors.New("jwt manager not configured: access_token_secret or refresh_token_secret is required")

// ParseAccessToken 离线解析并验证访问令牌
// 返回令牌中的 Claims 信息，包括 Subject（用户标识）、令牌类型、自定义扩展字段等
//
// 参数:
//   - token: 需要解析的访问令牌字符串
//
// 返回值:
//   - *jwt.Claims: 解析出的令牌声明
//   - error: 解析失败时返回错误（令牌无效、过期、签名错误等）
//
// 注意:
//   - 使用此方法前，必须在初始化 Client 时配置 AccessTokenSecret
//   - 若未配置，将返回 ErrJWTNotConfigured 错误
//
// 示例用法:
//
//	claims, err := client.ParseAccessToken(accessToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Subject: %s\n", claims.Subject)
//	fmt.Printf("TokenType: %s\n", claims.TokenType)
func (c *Client) ParseAccessToken(token string) (*jwt.Claims, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}
	if c.jwtManager == nil {
		return nil, ErrJWTNotConfigured
	}
	return c.jwtManager.ParseAccessToken(token)
}

// ParseRefreshToken 离线解析并验证刷新令牌
// 返回令牌中的 Claims 信息，包括 Subject（用户标识）、令牌类型等
//
// 参数:
//   - token: 需要解析的刷新令牌字符串
//
// 返回值:
//   - *jwt.Claims: 解析出的令牌声明
//   - error: 解析失败时返回错误（令牌无效、过期、签名错误等）
//
// 注意:
//   - 使用此方法前，必须在初始化 Client 时配置 RefreshTokenSecret
//   - 若未配置，将返回 ErrJWTNotConfigured 错误
//
// 示例用法:
//
//	claims, err := client.ParseRefreshToken(refreshToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Subject: %s\n", claims.Subject)
func (c *Client) ParseRefreshToken(token string) (*jwt.Claims, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}
	if c.jwtManager == nil {
		return nil, ErrJWTNotConfigured
	}
	return c.jwtManager.ParseRefreshToken(token)
}

// ValidateToken 离线验证令牌的有效性（不返回 Claims）
// 仅检查令牌签名和过期时间，不区分 access/refresh 类型
//
// 参数:
//   - token: 需要验证的令牌字符串
//
// 返回值:
//   - error: 若令牌有效返回 nil，否则返回具体错误
//
// 注意:
//   - 使用此方法前，必须在初始化 Client 时配置 AccessTokenSecret 或 RefreshTokenSecret
//   - 若未配置，将返回 ErrJWTNotConfigured 错误
func (c *Client) ValidateToken(token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}
	if c.jwtManager == nil {
		return ErrJWTNotConfigured
	}
	return c.jwtManager.ValidateToken(token)
}
