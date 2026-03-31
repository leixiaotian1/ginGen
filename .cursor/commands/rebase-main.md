---
description: 将当前分支 rebase 到最新 main
---

目标：把当前工作分支基于最新 `main` 重放提交历史，保持线性提交。

请按以下顺序执行：

1) 前置检查
- `git status`
- 若存在未提交改动，先提示我处理，不继续 rebase

2) 更新主分支引用
- `git fetch origin main`

3) 执行 rebase
- `git rebase origin/main`

4) 若发生冲突
- 列出冲突文件
- 提示我逐个解决后执行：`git add <file>` + `git rebase --continue`
- 若我明确要求终止，再执行 `git rebase --abort`

5) 完成后检查
- `git status`
- `git log --oneline --decorate -8`
- 说明是否需要 `git push --force-with-lease`（仅说明，不自动执行）

约束：
- 未经我明确确认，不进行任何强推
- 不做与 rebase 无关的代码修改
