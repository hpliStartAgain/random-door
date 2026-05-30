# AI 编排详细设计（ai-orchestration-detailed-design.md）

> 对应概要设计 19、20 章。这是 internal/ai 包的实现依据。约束见 ai-integration-rules.md（含合规红线）。
> 包含 3 个文件：llm_client.go / image_client.go / prompt_builder.go。
> 服务于 POST /api/chat 与 POST /api/checkin/generate-image。

## 0. 铁律重申
- 所有外部 AI 调用**只在本包**；API Key 只读环境变量(LLM_API_KEY / IMAGE_API_KEY)；**前端禁止直连**。
- 每个调用必须：超时(AI_TIMEOUT_SECONDS) + 有限重试 + 失败兜底且不阻断主链路。

## 1. prompt_builder.go — Prompt 组装

### 1.1 对话 Prompt 上下文输入
```text
type ChatContext struct {
  CityName      string
  CharacterName string
  Persona       string     // 来自 characters.persona(后端私有)
  Landmarks     []string   // 城市地标名
  Foods         []string   // 城市美食名
  DialectStyle  string     // characters.dialect_style
}
func BuildChatPrompt(ctx ChatContext) (systemPrompt string)
```

### 1.2 对话 Prompt 模板（对应概要 19.4）
```text
你现在扮演一个城市文化导览角色。

当前城市：{{city_name}}
当前人物：{{character_name}}
人物设定：{{persona}}
城市地标：{{landmarks}}
城市美食：{{foods}}
方言特点：{{dialect_style}}

回答要求：
1. 保持人物风格，但不要声称自己是真实复活的人。
2. 可以介绍城市历史、地标、美食和地方文化。
3. 可以使用少量方言词汇，但必须给出普通话解释。
4. 不要编造无法确认的历史事实。
5. 回答控制在 150 字以内。
6. 尽量引导用户继续探索当前城市。
```
- {{landmarks}}/{{foods}} 用顿号拼接；空则填"暂无"。
- 用户消息作为 user role 传入，不拼进 system。

### 1.3 生图 Prompt 输入与模板（对应概要 20.4）
```text
func BuildImagePrompt(cityName, landmarkName string) string
```
```text
Create a realistic travel photo of the uploaded person visiting {{landmark_name}} in {{city_name}}.
Keep the person's identity consistent.
Use the landmark image as background reference.
Natural daylight, tourist photo style, high quality, realistic composition.
```

## 2. llm_client.go — 对话调用

### 2.1 导出
```text
type LLMClient struct { baseURL, apiKey, model string; timeout time.Duration; maxRetry int }
func NewLLMClient(cfg) *LLMClient
func (c *LLMClient) Chat(ctx context.Context, systemPrompt, userMessage string) (reply string, err error)
```

### 2.2 主逻辑
```text
1. messages = [ {role:system, content:systemPrompt}, {role:user, content:userMessage} ]
2. 带 ctx+timeout 发起 HTTP POST 到 LLM_API_BASE(OpenAI 兼容 /chat/completions)
3. Header: Authorization: Bearer {LLM_API_KEY}
4. 失败重试：网络错误/5xx 最多 maxRetry 次(指数退避)
5. 解析 choices[0].message.content
6. 错误映射：
   - 超时(ctx.DeadlineExceeded) → ErrAITimeout → 上层 504 AI_TIMEOUT
   - 上游非2xx/解析失败 → ErrAIUpstream → 上层 502 AI_UPSTREAM_ERROR
```

### 2.3 合规约束
- system prompt 已含"不声称真实复活/不编史/方言需解释/150字内"。
- 不做声音克隆相关调用。

## 3. image_client.go — 生图调用

### 3.1 导出
```text
type ImageClient struct { baseURL, apiKey string; timeout time.Duration; maxRetry int }
func NewImageClient(cfg) *ImageClient
func (c *ImageClient) Generate(ctx, selfiePath, refImagePath, prompt string) (imageURL string, err error)
```

### 3.2 主逻辑
```text
1. 读取 selfie(用户上传) + refImage(景点参考图 landmark.image_url 对应本地文件)
2. 按外部生图 API 规范组装请求(multipart 或 base64)：prompt + 人像图 + 背景参考图
3. 带 ctx+timeout 发起调用；失败重试 maxRetry 次
4. 生成图二进制 → 交 upload.storage 落 uploads/generated/{uuid}.png
5. 返回可访问 URL(/uploads/generated/xxx.png)
6. 错误映射同 llm：超时→504，上游失败→502
```

### 3.3 合规约束
- 只处理"用户主动上传的本人照片"(由 upload.validator 前置把关类型/大小)。
- 不接受在世名人/他人恶搞(MVP 靠用户协议+提示，不做人脸识别)。
- prompt 不含色情/暴力/侮辱性词汇。

## 4. 兜底与健壮性（对应概要 20.5 / 25.3）
| 风险 | 应对 |
|---|---|
| 生图耗时长 | 设 AI_TIMEOUT_SECONDS；前端 loading + 可取消 |
| 生图失败/超时 | 返回明确错误码(502/504)，前端提示重试，**不阻断浏览与对话** |
| 人物一致性差 | prompt 加 Keep the person identity consistent |
| 演示稳定性 | 演示前为重点城市**预生成样例图**，必要时降级展示 |
| 对话失败 | 返回友好兜底文案，城市详情/打卡不受影响 |

## 5. 配置项（来自 .env，经 config 注入）
| env | 用途 |
|---|---|
| LLM_API_BASE / LLM_API_KEY / LLM_MODEL | 对话 |
| IMAGE_API_BASE / IMAGE_API_KEY | 生图 |
| AI_TIMEOUT_SECONDS | 统一超时 |
| (代码内常量) maxRetry | 重试次数，建议 1~2 |

## 6. 与 service 的边界
- chat_service：取上下文数据、落库 chat_messages；**Prompt 组装与调用在本包**。
- checkin_service：落 selfie、取参考图、落 generated、(确认后)写 checkins；**生图调用在本包**。