# AI 集成约束

> 涉及任何 AI 代码前必读。父约束：`AGENTS.md` / `CLAUDE.md`。

## 1. 架构约束

- 外部 AI 调用封装在 `backend/internal/ai/`。
- API Key 只读后端环境变量：`LLM_API_KEY` / `IMAGE_API_KEY`。
- 前端禁止直连任何 LLM 或图像生成 API。

## 2. 健壮性

- 每个外部调用必须设置 `AI_TIMEOUT_SECONDS` 超时与有限重试。
- LLM 未配置或失败时返回标准错误，由前端友好提示，城市浏览不受影响。
- 生图走 `ai_tasks` 异步任务；失败后可重试，达到最大尝试次数后标记 failed。
- `IMAGE_API_BASE=mock` 或未配置 Key 时允许使用本地占位图，保证非 AI 链路可用。

## 3. AI 对话合规红线

- 不得声称自己是“真实复活的人物”。
- 不编造无法确认的确定性史实。
- 可少量使用方言风格，但必须给普通话解释。
- 避免地域刻板印象；回答控制在约 150 字内。

## 4. AI 生图合规红线

- 只接受用户主动上传的本人照片。
- 禁止在世公众人物肖像、色情、暴力、侮辱性内容、他人照片恶搞。
- 不做真人声音克隆。

## 5. Prompt 来源

- 对话 Prompt 与生图 Prompt 模板见 `docs/design/ai-orchestration-detailed-design.md`。
- 人物 persona / prompt 来自数据库 `characters` 表；seed 仅提供初始化模板。

## 6. 提交前自检

- [ ] AI 调用都在后端。
- [ ] 无前端 Key、无硬编码 Key。
- [ ] 有超时、重试、兜底或异步任务状态。
- [ ] 合规红线在 prompt 与校验中体现。
