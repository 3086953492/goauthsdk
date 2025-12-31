package goauthsdk

import (
	"net/http"
	"strings"
)

// normalizeConfig 标准化配置
// - 去掉 BaseURL 末尾的 /
// - 若未提供 HTTPClient，补充默认 http.DefaultClient
func normalizeConfig(cfg *Config) {
	cfg.FrontendBaseURL = strings.TrimSuffix(cfg.FrontendBaseURL, "/")
	cfg.BackendBaseURL = strings.TrimSuffix(cfg.BackendBaseURL, "/")

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}
}
