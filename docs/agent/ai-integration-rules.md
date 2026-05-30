# AI 集成约束（含合规红线）

> 涉及任何 AI 代码前必读。父约束：`CLAUDE.md`。对应概要 19 / 20 / 23 章。

## 1. 架构约束
- 所有外部 AI 调用**封装在 backend/internal/ai/**（llm_client.go / image_client.go / prompt_builder.go）。
- **API Key 只读环境变量**（LLM_API_KEY / IMAGE_API_KEY），禁止硬编码、禁止下发前端。
- **前端禁止直连任何 AI API**，一律经后端中转。

## 2. 健壮性
- 每个外部调用必须设置 **超时**（AI_TIMEOUT_SECONDS）与 **重试**（有限次）。
- 失败有**兜底**且**不阻断主链路**：
  - 对话失败 → 返回友好提示，城市浏览不受影响。
  - 生图失败/超时 → 明确错误 + 可重试；演示前为重点城市准备**预生成样例图**。

## 3. AI 对话合规红线
- 不得声称自己是「真实复活的人物」。
- 不编造无法确认的确定性史实。
- 可用少量方言，但**必须给出普通话解释**。
- 避免地域刻板印象；回答控制在约 150 字内（见 prompt 模板）。

## 4. AI 生图合规红线
- 只接受**用户本人主动上传**的照片。
- 禁止：在世公众人物肖像、色情、暴力、侮辱性内容、他人照片恶搞。
- 不做真人声音克隆。

## 5. Prompt 来源
- 对话 Prompt 模板见 docs/design/ai-orchestration-detailed-design.md（概要 19.4）。
- 生图 Prompt 模板见同文档（概要 20.4）。
- 人物 persona 见 docs/product/ai-character-design.md。

## 6. 提交前自检
- [ ] AI 调用都在 internal/ai 内
- [ ] 无前端 Key、无硬编码 Key
- [ ] 有超时/重试/兜底
- [ ] 合规红线在 prompt 与校验中体现