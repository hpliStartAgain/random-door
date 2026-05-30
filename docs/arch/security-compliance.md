# 安全与合规设计（security-compliance.md）

> 对应概要设计 23 章。这是黑客松演示时体现"风险可控"的关键文档。红线在代码中由 upload/ai 包与中间件落实。

## 1. 上传图片安全（对应 23.1）
| 项 | 规则 |
|---|---|
| 文件类型 | 仅 jpg / jpeg / png / webp（upload.validator 校验，违反 → 415） |
| 文件大小 | ≤ 5MB（UPLOAD_MAX_SIZE_MB，违反 → 413） |
| 文件名 | 一律 UUID 重命名，**不信任用户原始文件名** |
| 路径安全 | 防路径穿越（拒绝 ../、绝对路径） |
| 目录隔离 | 上传图存 uploads/selfies，生成图存 uploads/generated |

## 2. AI 对话风控（对应 23.2）
1. 优先历史人物 / 低风险文化角色（见 ai-character-design.md）。
2. **不得让 AI 声称自己是真实复活的人物**。
3. 不鼓励编造确定性史实。
4. 方言以文化风格化为主，避免地域刻板印象。
5. 异常输出有兜底友好提示。
（以上约束已写入对话 system prompt，见 ai-orchestration-detailed-design.md。）

## 3. 生图风控（对应 23.3）
1. 仅支持用户**主动上传本人照片**。
2. 不支持上传他人照片进行恶搞。
3. 不生成色情、暴力、侮辱性内容。
4. 不使用真实在世公众人物肖像作为生成目标。
5. 不做真人声音克隆。

## 4. API 与密钥安全（对应 23.4）
| 项 | 规则 |
|---|---|
| 密钥存放 | 外部 API Key 只放后端环境变量（LLM_API_KEY/IMAGE_API_KEY） |
| 前端 | 前端**不得直连**任何模型 API；仅可持有 VITE_AMAP_KEY（高德公开 Key） |
| 上传接口 | 限大小 + 限频（rate_limit 中间件） |
| AI 接口 | 设超时（AI_TIMEOUT_SECONDS），失败有兜底 |
| CORS | 白名单（CORS_ALLOW_ORIGINS），cors 中间件控制 |
| 错误处理 | 后端统一错误响应，不向前端泄露内部细节/堆栈 |

## 5. 合规红线汇总（一页速查 / 演示问答用）
- ❌ AI 自称真实复活人物　❌ 编造确定性史实　❌ 在世名人肖像生图
- ❌ 他人照片恶搞　❌ 色情/暴力/侮辱内容　❌ 声音克隆
- ❌ 前端出现/调用 AI Key　❌ 接受非白名单类型/超大文件
- ✅ 历史人物优先　✅ 用户本人照片　✅ 方言必附普通话解释　✅ Key 仅后端

## 6. 责任落点（代码层）
| 红线 | 落实位置 |
|---|---|
| 上传校验 | internal/upload/validator.go |
| 对话约束 | internal/ai/prompt_builder.go（system prompt） |
| 密钥隔离 | internal/config + .env（前端无 Key） |
| CORS/限流 | internal/middleware/cors.go、rate_limit.go |
| 统一错误 | internal/middleware/recover.go + handler |