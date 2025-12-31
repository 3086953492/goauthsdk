package configx

import "fmt"

// Validate 校验配置的必填字段
// 错误字符串遵循规则：小写开头、无结尾标点，且不泄漏敏感信息
func Validate(cfg *Config) error {
	if cfg.FrontendBaseURL == "" {
		return fmt.Errorf("frontend_base_url is required")
	}
	if cfg.BackendBaseURL == "" {
		return fmt.Errorf("backend_base_url is required")
	}
	if cfg.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if cfg.ClientSecret == "" {
		return fmt.Errorf("client_secret is required")
	}
	if cfg.RedirectURI == "" {
		return fmt.Errorf("redirect_uri is required")
	}
	return nil
}
