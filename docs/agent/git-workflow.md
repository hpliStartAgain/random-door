# Git 工作流约束

> 提交前必读。仓库可能存在他人未提交改动，任何 agent 都不得擅自回滚不属于自己的内容。

## 1. 基本原则

- 提交前先看 `git status --short`。
- 只暂存本次任务相关文件。
- 不执行 `git reset --hard`、`git checkout -- <path>` 等回滚命令，除非用户明确要求。
- 发现无关脏文件时保持原样，不纳入提交、不删除。

## 2. 提交内容

- 文档整理、代码功能、配置变更分清边界。
- 接口或表结构变更必须带对应 docs 更新。
- 不提交 `.env`、密钥、`uploads/` 运行时产物、构建产物、缓存文件。

## 3. 推荐命令

```bash
git status --short
git diff -- <paths>
git add <paths>
git commit -m "docs: align project documentation"
```

## 4. 提交前自检

- [ ] `git diff --cached` 只包含本次任务文件。
- [ ] README / TODO / CHANGELOG 与实际变更一致。
- [ ] 测试或文档审计命令已记录在最终说明中。
