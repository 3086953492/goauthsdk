package goauthsdk

// UserInfo 用户信息结构体
// 用于 GET /api/v1/oauth/userinfo 接口返回的用户信息
type UserInfo struct {
	Sub       string `json:"sub"`        // 用户唯一标识（用户ID）
	Nickname  string `json:"nickname"`   // 用户昵称
	Picture   string `json:"picture"`    // 用户头像URL
	UpdatedAt int64  `json:"updated_at"` // 用户信息更新时间（Unix 时间戳，秒）
}

// UserDetail 用户详情结构体
// 用于 GET /api/v1/users/{id} 接口返回的用户详情
type UserDetail struct {
	ID        uint64 `json:"id"`         // 用户主键 ID
	Subject   string `json:"subject"`    // 用户唯一标识（对外使用，推荐作为用户标识符）
	Username  string `json:"username"`   // 用户名
	Nickname  string `json:"nickname"`   // 昵称
	Avatar    string `json:"avatar"`     // 头像 URL
	Status    int    `json:"status"`     // 状态：1=正常，0=禁用
	Role      string `json:"role"`       // 角色：user / admin
	CreatedAt string `json:"created_at"` // 创建时间（ISO 8601 格式）
	UpdatedAt string `json:"updated_at"` // 更新时间（ISO 8601 格式）
}
