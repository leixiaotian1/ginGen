# 基线回归清单（ginGen）

本文档描述当前 CLI / Web 生成与 `add` 功能的**可验证基线**，供重构与 CI 对照。

## 命令基线

| 场景 | 命令 / 请求 | 期望 |
|------|-------------|------|
| 新建工程 | `ginGen new demo -m example.com/demo` | 生成 `go.mod`、`cmd/server/main.go`、`internal/router/router.go`、`internal/config/config.go`、`configs/config.yaml` 等；目录 `internal/service`、`internal/model`、`internal/middleware`、`pkg` 存在 |
| 追加功能 | `ginGen add mysql .`（在已有模块内） | 拉取 GORM/MySQL 相关依赖（或 `--force` 下允许失败后继续）、生成 `internal/config/db_config.go`、`internal/clients/gorm.go`；并自动合并 `configs/config.yaml` 与 `internal/config/config.go`（与 Web 生成路径一致） |
| 功能别名 | `gorm`≈`mysql`，`postgresql`≈`postgres`，`logging`≈`logger`，`hot-reload`≈`hotreload` | 与 canonical id 行为一致 |
| Web 列表 | `GET /api/features` | 返回与内置注册表一致的 `id` / 文案 / 分类 |
| Web 生成 | `POST /api/generate` JSON：`projectName`、`modulePath`、`features[]` | ZIP 内工程结构与 CLI `new` + 多次 `add` 一致（含配置补丁） |

## 功能矩阵（canonical id）

`mysql`, `postgres`, `redis`, `kafka`, `jwt`, `logger`, `swagger`, `middleware`, `health`, `prometheus`, `hotreload`, `cron`, `upload`

- **含 `go get` 的 feature**：mysql、gorm 路径、postgres、redis、kafka、jwt、logger、swagger、middleware、prometheus、hotreload、cron  
- **仅模板的 feature**：health、upload  

## 回归测试建议（自动化）

1. **无网络单测**：`new` 结构生成 + `add` 仅校验落盘文件（`SkipGoGet`）。  
2. **可选集成**（CI 可关或单独 job）：临时目录 `go mod tidy` / `go build ./...`。  

## 已知手工步骤（文档层）

部分功能仍可能需在业务代码中注册路由或初始化客户端；以生成后 `README` 与模板注释为准。
