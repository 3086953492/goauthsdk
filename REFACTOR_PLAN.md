# goauthsdk 分布式改造计划（不考虑向后兼容）

本文档基于仓库内现行规则：
- `.cursor/rules/codestyle.mdc`：代码风格与工程约束

目标：把当前“根目录平铺 + 逻辑分散”的形态，改造成**层次清晰、错误一致、HTTP 一致、可演进**的 SDK 工程结构。

> 说明：本计划不要求一次性完成。建议按阶段逐步推进，每次提交都可回滚、可定位。

## 目标目录结构（建议）

根包 `goauthsdk` 仅保留对外 API（稳定入口），实现细节下沉 `internal/`。

建议最终形态：

- `client.go`：`Client` 对外入口与依赖注入
- `config.go`：`Config` 与校验/标准化（可拆到 `config_validate.go`）
- `types/`（可选）或保留 `types.go` 但按域分组：对外 DTO（与后端接口对齐）
- `internal/httpx/`：HTTP 通用能力（构建、发送、读 body、通用 header）
- `internal/apierr/`：统一错误模型与解码（ProblemDetails / code-message）
- `internal/validate/`（可选）：参数校验工具（避免散落）
- `cmd/goauthsdk-testserver/`：手工验证服务（可应用化，但只能依赖公开 API）

## 执行原则

- **不考虑向后兼容**：允许改导出 API、错误类型、文件/包结构。
- **一致性优先**：同类接口必须同样的请求/错误/解析策略。
- **先收敛通用能力再改业务**：先统一错误与 HTTP，再统一各业务接口文件。
- **安全优先**：错误/日志中禁止输出密钥与 token 明文。

## 分阶段计划

### 阶段 0：定边界与骨架（低风险）

目标：先把“未来要往哪里放”确定下来，让后续改造有稳定落点。

- **动作**
  - 创建 `internal/httpx`、`internal/apierr`（允许先是最小实现）
  - 约定各接口文件的三段式结构：`buildXRequest` / `doXRequest` / `parseXResponse`
  - 约定 DTO 的归属（对外 vs 内部）
- **验收**
  - 仓库结构已出现 `internal/` 目录，且导入关系清晰
  - `cmd/...` 不依赖 `internal/...`
- **风险**
  - 仅搬目录通常不会引入逻辑变更；注意保持 package 名一致

### 阶段 1：错误体系收敛（收益最大）

目标：让 SDK 的错误可预测、可识别、可 `errors.Is/As`，并且不泄漏敏感信息。

- **动作**
  - 全仓库错误字符串统一为 **小写开头、无结尾标点**
  - 统一包装错误：`fmt.Errorf("<verb phrase>: %w", err)`
  - 新增统一错误类型（建议）：`APIError`（包含 `Status/Code/Detail` 等）
  - 新增统一入口：`decodeAPIError(resp, body)`，用于解析：
    - RFC 7807 `ProblemDetails`
    - 通用 `{code,message,data}` 结构（如适用）
    - 兜底：只暴露必要信息（不输出敏感内容；body 截断仅用于调试）
- **验收**
  - 调用方可以 `errors.As(err, *APIError)` 获取结构化信息
  - 所有 API 在非 2xx 情况下走统一的错误解析路径
- **风险**
  - 错误信息变化会影响调用方日志/告警关键词（但本计划不考虑兼容）

### 阶段 2：HTTP 通用逻辑下沉（build/do/parse 标准化）

目标：每个 API 文件只保留“本接口特有”的内容，HTTP 发送/读 body/通用 header 等全部复用。

- **动作**
  - 在 `internal/httpx` 收敛：
    - `Do(req)`：发送请求并读取 body（统一 Close）
    - 可选：默认 User-Agent、通用 Header（如 Trace-ID）
  - 在每个接口的 `parseXResponse`：
    - 先调用 `decodeAPIError`
    - 成功时再 decode 业务响应
  - 删除重复的 `truncateBody` 等散落工具函数，改为集中实现
- **验收**
  - 所有网络调用都接收 `context.Context`，且库内部不使用 `context.Background()`
  - `doXRequest` 逻辑几乎一致或直接复用 `internal/httpx`
- **风险**
  - 统一读 body/错误路径后，某些接口原先的“特殊处理”可能丢失，需要逐个确认

### 阶段 3：DTO 与类型治理（对齐后端、可读、可维护）

目标：DTO 清晰可读，避免“过度复用导致语义混乱”。

- **动作**
  - DTO 按域拆分（推荐二选一）：
    - A：保留 `types.go`，但按域分组并加清晰注释
    - B：拆成 `token_types.go` / `users_types.go` / `oauth_types.go` 等
  - 统一时间字段策略：
    - Unix 时间戳：`int64` 并注明单位
    - ISO 8601：`string`（SDK 不强行解析）
  - 内部专用 DTO 不导出，必要时下沉 `internal/`
- **验收**
  - JSON tag 全部 `snake_case`
  - 导出类型/字段注释符合规则（中文注释、说明含义与来源）
- **风险**
  - 类型拆分可能影响 import 与循环依赖，需要小步推进

### 阶段 4：Client / Config 重塑（入口变干净）

目标：`Client` 只承载“配置 + 依赖 + 方法”，实现细节不堆在入口文件里。

- **动作**
  - `Config` 校验/标准化逻辑拆分清晰（校验错误字符串按规则）
  - 明确 `HTTPClient` 注入策略（必须可注入）
  - JWT 离线验签能力明确边界：可独立为 verifier/manager（视现状决定）
- **验收**
  - `NewClient` 的副作用最小、可预测；敏感字段不出现在错误里
  - `Client` 代码体量下降，实现细节迁移到 `internal/`
- **风险**
  - 入口重塑通常会带来导出 API 变化（本计划允许）

### 阶段 5：testserver 收尾（只依赖公开 API）

目标：`cmd/goauthsdk-testserver` 继续作为手工验证入口，但与 SDK 私有实现解耦。

- **动作**
  - 清理对 `internal/` 的引用（如有）
  - 随 SDK 破坏性变更同步更新 testserver
- **验收**
  - testserver 仅依赖 `goauthsdk` 公共 API
- **风险**
  - 若 SDK 改动较大，需要同步更新路由/handler 参数与返回

## 推荐切入点（按收益/复杂度排序）

1. **token 相关**：通常最核心、重复逻辑多，适合先做错误与 HTTP 收敛
2. **userinfo/users**：DTO 清晰度提升明显
3. **introspect/revoke/authorize**：补齐一致性与边界