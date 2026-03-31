# 贡献指南

## 开发环境

- Go 版本见 [`go.mod`](go.mod) 的 `go` 指令。

## 常用命令

```bash
make fmt    # gofmt
make vet
make test
make build  # 输出到 bin/ginGen，并注入 version/commit/date
```

CI 会执行 `gofmt -l`、`go vet ./...`、`go test ./...`。

## 架构说明

- **CLI**：[`cmd/`](cmd/)（Cobra）
- **功能管线**：[`internal/feature/`](internal/feature/) — CLI 与 Web 共用，避免重复逻辑
- **工程生成**：[`internal/generator/`](internal/generator/) — 模板与 `new` 工程结构
- **用户配置**：[`internal/preset/`](internal/preset/) — 合并 `~/.gingen.yaml` 与当前目录 `.gingen.yaml`
- **Web**：[`internal/web/`](internal/web/)

新增功能时请在 `internal/feature/registry.go` 注册，并实现对应 `apply*.go` 中的函数。

## 回归与文档

- 基线清单：[`docs/BASELINE_REGRESSION.md`](docs/BASELINE_REGRESSION.md)

## 发布（可选）

使用 [GoReleaser](https://goreleaser.com/) 时参考根目录 [`.goreleaser.yaml`](.goreleaser.yaml)；版本信息通过 `-ldflags` 写入 [`internal/version`](internal/version/version.go)。
