# AI 编排详细设计

## 1. 总原则

- LLM 与图像生成 API 只由后端调用。
- API Key 只从环境变量读取，不硬编码、不下发前端。
- 所有 AI 调用有超时、有限重试、错误映射和成本限制。
- AI 能力失败不应阻断城市浏览、访问记录、评论、资产等非 AI 主链路。

## 2. 对话链路

```text
POST /api/chat
  -> chat_service 校验 user/city/character/message
  -> AI 用量限制 IncrementIfBelow(chat)
  -> 读取历史消息
  -> prompt_builder.BuildChatPrompt
  -> llm_client.ChatWithHistory
  -> 落库 user 与 assistant 消息
  -> 返回 reply
```

### Prompt 约束

```text
你在和用户进行角色扮演的游戏，你扮演的人物是{{character_name}}。当前城市是{{city_name}}。
不要声称自己是真实复活的人，不编造无法确认的确定性史实。回答控制在150字以内。
{{dialect_rule}}
```

`dialect_rule` 仅在人物有 `dialect_style` 时追加，并要求附普通话解释。

### LLM client

- base URL：`LLM_API_BASE`，按 OpenAI 兼容 `/chat/completions` 调用。
- Header：`Authorization: Bearer <LLM_API_KEY>`。
- 未配置 base/key 或 base 为 `mock` 时返回 `AI_UPSTREAM_ERROR`。
- 超时映射 `AI_TIMEOUT`，上游错误映射 `AI_UPSTREAM_ERROR`。

## 3. 生图链路

```text
POST /api/checkin/generate-image
  -> checkin_handler 校验 multipart 文件
  -> storage 保存自拍与可选 scene_file
  -> checkin_service 校验 user/city/landmark
  -> AI 用量限制 IncrementIfBelow(image)
  -> 写 ai_tasks(status=queued,type=checkin_image)
  -> 返回 task_id

worker
  -> ClaimNext(checkin_image)
  -> image_client.Generate(selfie, refImage, prompt)
  -> storage.SaveBytes uploads/generated
  -> MarkSucceeded(result_url) 或 MarkRetryable/MarkFailed

frontend
  -> GET /api/checkin/image-tasks/{task_id}
  -> 失败可 POST /retry
  -> 成功后 POST /api/checkin 确认打卡
```

### 生图 Prompt

```text
请生成一张真实旅行打卡照：保留用户上传自拍中的人物身份、脸部特征和自然姿态，
将人物自然合成到{{city_name}}的{{landmark_name}}场景中。可参考地标图片作为背景风格和构图依据。
自然日光、游客照片风格、高质量、真实构图。禁止生成在世公众人物、色情、暴力或侮辱内容。
```

### Image client

- `IMAGE_API_BASE=mock` 或 Key 为空时返回本地占位 PNG。
- DashScope 模式调用 multimodal generation。
- 支持从 URL、data URL 或 base64 响应解析图像。
- 生成图由 upload storage 保存为 `/uploads/generated/*.png`。

## 4. 用量与任务配置

| env | 用途 |
|---|---|
| `AI_TIMEOUT_SECONDS` | LLM / 生图统一超时。 |
| `AI_CHAT_DAILY_LIMIT` | 单用户每日对话次数。 |
| `AI_IMAGE_DAILY_LIMIT` | 单用户每日生图任务次数。 |
| `AI_WORKER_INTERVAL_SECONDS` | 生图 worker 轮询间隔。 |
| `AI_WORKER_CONCURRENCY` | worker 并发数。 |
| `AI_MAX_TASK_ATTEMPTS` | 生图任务最大尝试次数。 |
| `LLM_API_BASE` / `LLM_API_KEY` / `LLM_MODEL` | 对话模型配置。 |
| `IMAGE_API_BASE` / `IMAGE_API_KEY` / `IMAGE_MODEL` | 图像模型配置。 |

## 5. 合规边界

- 不生成在世公众人物肖像。
- 不处理色情、暴力、侮辱性内容。
- 不做声音克隆。
- 不让 AI 自称真实复活人物。
- 人物 persona / prompt 不下发前端。

## 6. Service 边界

| service | 职责 |
|---|---|
| `chat_service` | 校验、用量限制、历史消息、调用 LLM、保存消息。 |
| `checkin_service` | 文件任务入队、worker 处理、任务状态、确认打卡。 |
| `guess_service` | 猜城市文案 prompt 与 LLM 兜底。 |
