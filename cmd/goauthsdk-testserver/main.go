package main

import (
	"log"
)

// ============================================================================
// goauthsdk 手工测试服务 - 用于开发/测试阶段手动验证 OAuth 流程
// ============================================================================

func main() {
	log.Printf("启动 goauthsdk 测试服务于 %s", serverAddr)
	log.Printf("配置信息:")
	log.Printf("  - 前端地址: %s", testFrontendBaseURL)
	log.Printf("  - 后端地址: %s", testBackendBaseURL)
	log.Printf("  - 客户端ID: %s", testClientID)
	log.Printf("  - 回调地址: %s", testRedirectURI)
	log.Printf("\n访问 http://localhost%s/ 查看使用说明\n", serverAddr)

	if err := startServer(serverAddr); err != nil {
		log.Fatal(err)
	}
}
