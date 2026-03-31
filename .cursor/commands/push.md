---
description: 检查分支状态并安全推送到远程
---

目标：将当前分支安全推送到远程仓库。

请按以下顺序执行：

1) 运行检查
- `git status`
- `git branch -vv`
- `git remote -v`

2) 判断是否可推送
- 若工作区有未提交改动，先提示我处理
- 若当前分支未设置 upstream，使用 `git push -u origin HEAD`
- 若已设置 upstream，使用 `git push`

3) 推送后验证
- 重新执行 `git status`
- 简要汇报推送结果（分支、远程、最新提交）

约束：
- 不使用 `--force`
- 如远程拒绝，先解释原因并给出修复建议，不做破坏性操作
