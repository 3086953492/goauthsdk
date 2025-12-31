package goauthsdk

import (
	"fmt"

	"github.com/3086953492/gokit/jwt"
)

// JWTVerifier 提供 JWT 离线验签能力
// 可独立使用，也可由 Client 持有
type JWTVerifier struct {
	manager *jwt.Manager
}

// NewJWTVerifier 创建一个新的 JWTVerifier
// accessTokenSecret 和 refreshTokenSecret 至少提供一个
//
// 参数:
//   - accessTokenSecret: 访问令牌签名密钥
//   - refreshTokenSecret: 刷新令牌签名密钥
//
// 返回值:
//   - *JWTVerifier: 可用于离线验签的 verifier 实例
//   - error: 两个 secret 均为空、或底层 manager 创建失败时返回错误
func NewJWTVerifier(accessTokenSecret, refreshTokenSecret string) (*JWTVerifier, error) {
	if accessTokenSecret == "" && refreshTokenSecret == "" {
		return nil, fmt.Errorf("at least one of access_token_secret or refresh_token_secret is required")
	}
	mgr, err := jwt.NewManager(
		jwt.WithAccessSecret(accessTokenSecret),
		jwt.WithRefreshSecret(refreshTokenSecret),
	)
	if err != nil {
		return nil, fmt.Errorf("create jwt manager: %w", err)
	}
	return &JWTVerifier{manager: mgr}, nil
}

// ParseAccessToken 离线解析并验证访问令牌
// 返回令牌中的 Claims 信息
//
// 参数:
//   - token: 需要解析的访问令牌字符串
//
// 返回值:
//   - *jwt.Claims: 解析出的令牌声明
//   - error: 令牌无效、过期、签名错误等情况返回错误
func (v *JWTVerifier) ParseAccessToken(token string) (*jwt.Claims, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}
	return v.manager.ParseAccessToken(token)
}

// ParseRefreshToken 离线解析并验证刷新令牌
// 返回令牌中的 Claims 信息
//
// 参数:
//   - token: 需要解析的刷新令牌字符串
//
// 返回值:
//   - *jwt.Claims: 解析出的令牌声明
//   - error: 令牌无效、过期、签名错误等情况返回错误
func (v *JWTVerifier) ParseRefreshToken(token string) (*jwt.Claims, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}
	return v.manager.ParseRefreshToken(token)
}

// ValidateToken 离线验证令牌的有效性（不返回 Claims）
// 仅检查令牌签名和过期时间，不区分 access/refresh 类型
//
// 参数:
//   - token: 需要验证的令牌字符串
//
// 返回值:
//   - error: 若令牌有效返回 nil，否则返回具体错误
func (v *JWTVerifier) ValidateToken(token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}
	return v.manager.ValidateToken(token)
}
