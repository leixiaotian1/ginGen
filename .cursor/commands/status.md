---
description: 快速汇总当前仓库与分支状态
---

目标：在最短时间内给出当前 Git 状态全景，方便我决定下一步动作。

请按以下顺序执行：

1) 基础状态
- `git status`
- `git branch -vv`
- `git remote -v`

2) 变更概览
- `git diff --stat`
- `git diff --staged --stat`

3) 提交上下文
- `git log --oneline --decorate -10`

4) 输出结论
- 用 3-6 条要点总结：当前分支、是否落后/领先、是否有未提交改动、建议下一步（commit/push/pr）

约束：
- 只读操作，不修改代码和 Git 历史
