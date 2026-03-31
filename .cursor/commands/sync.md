---
description: 同步远程主分支并更新当前分支
---

目标：把当前分支与远程主分支同步到最新，减少后续冲突。

请按以下顺序执行：

1) 前置检查
- `git status`
- 若有未提交改动，先提示我处理（提交、暂存或放弃），不要继续

2) 拉取远程信息
- `git fetch --all --prune`

3) 识别主分支
- 优先使用 `main`
- 若仓库使用 `master` 或其他默认分支，请自动识别并说明

4) 同步当前分支
- 默认采用 rebase：`git rebase origin/<base>`
- 如果 rebase 失败，先停止并说明冲突文件与建议处理方式

5) 结果验证
- `git status`
- `git log --oneline --decorate -5`
- 简要说明是否需要 `git push --force-with-lease`（仅建议，不自动执行）

约束：
- 不自动执行 `--force` 或 `--force-with-lease`
- 不修改与同步无关的文件
