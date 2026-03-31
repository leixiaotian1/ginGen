---
description: 基于当前分支创建规范 Pull Request
---

目标：为当前分支创建一个清晰、可审查的 Pull Request。

请按以下顺序执行：

1) 收集上下文
- `git status`
- `git branch -vv`
- `git log --oneline --decorate -10`
- 识别基线分支（默认 `main`，如仓库约定不同请自动识别）
- 运行 `git diff <base>...HEAD` 了解完整变更

2) 生成 PR 内容草稿
- 给出精炼标题
- 产出 PR 描述，至少包含：
  - 变更动机
  - 主要改动点（1-3 条）
  - 风险与影响范围
  - 测试说明

3) 如当前分支尚未推送，先推送
- 使用 `git push -u origin HEAD`

4) 创建 PR
- 使用 `gh pr create` 创建
- 创建后返回 PR 链接

约束：
- 未经我允许，不修改代码
- 若 `gh` 未登录或权限不足，先提示我完成认证
