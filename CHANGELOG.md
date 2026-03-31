# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- `internal/feature`：统一 CLI / Web 的功能注册与生成管线
- `.gingen.yaml` / `.gingen.yaml.example`：`defaultModule`、`templateRoot`、`presets`
- `new --features`、`new --preset`、`new`/`add --template-root`；`add --force` 严格/宽松策略
- `internal/version` 与 Makefile / GoReleaser 的 `-ldflags` 注入
- CI（GitHub Actions）、`CONTRIBUTING.md`、`docs/BASELINE_REGRESSION.md`
- 生成器与功能的单元测试

### Changed

- `add` 与 Web 生成路径对齐：MySQL/Redis/Kafka 等自动合并 `config.yaml` 与 `config.go`
- `swagger` 功能同时生成 `docs/swagger.go` 与 `docs/swagger_annotations.txt`
