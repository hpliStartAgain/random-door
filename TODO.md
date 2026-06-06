---
status: active
branch: master
owner: devin
updated: 2026-06-06 09:30
tier: COMPLEX
---

> v1.0（演示前夕的产品打磨）已交付，参见 CHANGELOG `[1.0.0]`。
> 本轮目标：v1.1 — 把 seed 从精选 12 城扩至 35 城，并建立**可重复的批量补全工作流**（数据 + 图片）。

## 1. 需求理解

把"能跑通的 12 城 demo"扩展成"覆盖 35 城的丰满产品"。三个具体输入/输出：

- **输入**：DB 现状 35 城（其中 12 城 seed 完整 + 23 城仅 name/province/lat/lng + 8 张基础描述,无图无 tag）。
- **输出**：
  1. `backend/data/seed/cities.json` 扩到 35 城,**每城都满足 seed.ValidateCatalog 全部约束**（1-2 landmarks / 1-2 foods / 正好 1 character / 合规 prompt / 真实方言）。
  2. `backend/static/{landmarks,foods,characters,cities}/` 下补齐对应图片素材（约 35 cover + 35-70 landmarks + 35-70 foods + 35 avatars ≈ 200 张）。
  3. 一套**可复用的工作流脚本** `scripts/build_seed.py`,后续要扩 50/100 城时直接复用。
- **边界**：
  - 12 城原有数据**不动**（已经精修过）,只增量扩 23 城。
  - 不改 DB schema、不改 API 契约、不改前端。
  - 不引入新中间件（Redis / 向量库 / 队列均不要）。

## 2. 设计方案

**核心思路**：双流水线 + 单一脚本入口 + 字段元数据即 prompt 模板变量。

```
┌──────────────────────────────────────────────────────────────────┐
│ 用户准备:                                                          │
│   .env 填好 LLM_API_KEY / IMAGE_API_KEY / UNSPLASH_ACCESS_KEY      │
│   编辑 scripts/seed_inputs.csv (扩城清单: name,province,lat,lng)   │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─── 流水线 A: 数据生成（一次性） ────────────────────────────────────┐
│ python scripts/build_seed.py gen-data                              │
│   for each row in seed_inputs.csv:                                 │
│     LLM ( Prompt-CITY ) → 单城 JSON 片段                            │
│     合并进 cities.json（追加,不覆盖已有 12 城）                       │
│ python scripts/build_seed.py validate → 跑 seed.ValidateCatalog     │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─── 流水线 B: 图片获取（幂等增量） ──────────────────────────────────┐
│ python scripts/build_seed.py gen-images                            │
│   scan cities.json,for each image_url:                             │
│     if static/<path> exists → skip                                 │
│     else:                                                          │
│       kind=characters → IMAGE_API (插画风)                          │
│       kind in [landmarks,foods,cities] →                           │
│           Wikimedia Commons API → Unsplash API → IMAGE_API（兜底）  │
│       下载,做尺寸标准化（Pillow）,落到 static/<path>                  │
└──────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─── 收尾 ─────────────────────────────────────────────────────────┐
│ 改 ExpectedCityCount = 35                                          │
│ 重启服务 → seed.Load 自动 upsert → DB 与 JSON 对齐                   │
└──────────────────────────────────────────────────────────────────┘
```

**关键决策**：
- **铁律 4 升级**：从「12 城」改为「精选 N 个核心城市,N 由 seed JSON 决定」,并把 `ExpectedCityCount` 从硬编码改为 **`MinCityCount` + `MaxCityCount` 区间校验**（12-100）。
- **图片三级回落**：Wikimedia → Unsplash → AI 生图。每级失败有清晰日志,允许人工指定某条走特定源。
- **agent 调用接口**：通过 `.devin/skills/seed-builder/SKILL.md` 暴露,后续任何 coding agent（Devin / Claude Code / Codex）都能识别并按工作流执行。
- **幂等性**：脚本必须支持"中途中断重跑"、"已有文件不重生成"、"指定单城重跑"三种姿势。

## 3. 阶段划分

- [ ] **Phase 1**: 基础设施改造（铁律 / ExpectedCityCount / .env.example）— 解锁后续扩城
- [ ] **Phase 2**: scripts/build_seed.py 骨架 + skill 文档 + seed_inputs.csv 模板
- [ ] **Phase 3**: 实现 `gen-data` 子命令（LLM 调用 + JSON 合并 + 重试）
- [ ] **Phase 4**: 实现 `gen-images` 子命令（Wikimedia / Unsplash / AI 三级回落 + 标准化）
- [ ] **Phase 5**: 实现 `validate` 子命令（包装 `go run -tags seed_check`）
- [ ] **Phase 6**: 端到端跑通 1 城（端到端验证）→ 扩至 23 城 → 验证 DB 对齐
- [ ] **Phase 7**: 文档收尾（CHANGELOG [1.1.0] + README 城市数表述）

## 4. 文件级任务

### Phase 1：基础设施

| 文件 | 动作 | 说明 |
|------|------|------|
| AGENTS.md | MODIFY | 第 4 条铁律改为「精选核心城市,以 cities.json 为准,数量区间 12-100」 |
| CLAUDE.md | MODIFY | 同步修改第 4 条铁律 |
| backend/internal/seed/seed.go | MODIFY | `ExpectedCityCount = 12` 改为 `MinCityCount=12, MaxCityCount=100` 区间校验 |
| backend/internal/seed/seed_test.go | MODIFY | 调整测试断言 |
| .env.example | MODIFY | 新增 `UNSPLASH_ACCESS_KEY`、`WIKIMEDIA_USER_AGENT` 注释 |

### Phase 2：脚手架

| 文件 | 动作 | 说明 |
|------|------|------|
| scripts/build_seed.py | NEW | 主入口,argparse 三个子命令骨架 |
| scripts/seed_builder/__init__.py | NEW | Python 包初始化 |
| scripts/seed_builder/llm_client.py | NEW | 调 LLM_API 的小客户端（与 backend image_client 行为对齐） |
| scripts/seed_builder/image_sources.py | NEW | Wikimedia / Unsplash / AI 三个 source 实现 |
| scripts/seed_builder/prompts.py | NEW | Prompt-CITY / Prompt-IMAGE 模板常量 |
| scripts/seed_builder/pinyin_slug.py | NEW | 中文→拼音 slug 工具（pypinyin） |
| scripts/seed_inputs.csv | NEW | 扩城清单模板（含已知 23 城坐标） |
| scripts/requirements.txt | NEW | requests, pypinyin, Pillow（python-dotenv 已在使用？需确认） |
| .devin/skills/seed-builder/SKILL.md | NEW | 让 coding agent 能识别并自动调用本工作流 |

### Phase 3-5：脚本逻辑实现（覆盖 Phase 2 骨架）

scripts/seed_builder/{generator,image_pipeline,validator}.py 各模块的具体实现细节,在 Phase 2 设计完骨架后再具体定。

### Phase 6-7：产物 + 收尾

| 文件 | 动作 | 说明 |
|------|------|------|
| backend/data/seed/cities.json | MODIFY | 脚本产物,扩至 35 城 |
| backend/static/{cities,landmarks,foods,characters}/ | NEW（增量） | 约 200 张图,均来自 Wikimedia/Unsplash/IMAGE_API |
| CHANGELOG.md | MODIFY | 新增 [1.1.0] 段落 |
| README.md | MODIFY | 城市数表述（如果有"12 城"字眼） |

## 5. 外部变更（必须显式列出）

- [ ] **铁律修改**：AGENTS.md / CLAUDE.md 第 4 条
- [ ] **新增 Python 依赖**：`requests`、`pypinyin`、`Pillow`、`python-dotenv`（pymysql 已装,仅 audit 用,不入 requirements）
- [ ] **新增环境变量**：
  - `UNSPLASH_ACCESS_KEY`（必填,自助申请 https://unsplash.com/developers）
  - `WIKIMEDIA_USER_AGENT`（Wikimedia API 礼仪要求,如 `city-roam-seed/1.0 (contact@example.com)`）
  - 复用 `LLM_API_KEY`、`IMAGE_API_KEY`、`LLM_MODEL`（已存在）
- [ ] **配置文件**：`.env.example` 同步;`.env` 由老板自己填
- [ ] **代码常量变更**：`ExpectedCityCount` 单值→范围
- [ ] **不动**：DB schema、API 契约、前端、技术栈

## 6. 待确认问题（@老板）

- **Q1**：`ExpectedCityCount` 改为 `MinCityCount=12, MaxCityCount=100` 区间校验?还是简单改成 `ExpectedCityCount=35`?
  - **倾向 A（区间）**：避免下次扩城再改代码,但需调整一处测试。
- **Q2**：cities.json 里的 `image_url` 字段允许填 **https://** 远程 URL 吗?
  - 现有 seed.go `requireStaticURL` 强制 `/static/` 开头。如果允许远程,需放宽校验,但好处是省下载步骤。
  - **倾向 B（保持 /static/）**：本地兜底,演示离线也能跑;一致性强。
- **Q3**：文件命名规则用 `<拼音>_<序号>.jpg`（如 `tianjin_1.jpg`）还是 `<城市拼音>_<地标拼音>.jpg`（如 `tianjin_wudadao.jpg`）?
  - **倾向 B（语义命名）**：可读性强,出问题好定位。
- **Q4**：现有 DB 里 23 个 extra 城市已有 `intro / dialect_sample` 等字段值,LLM 重新生成时是否参考?
  - **倾向 A（不参考,LLM 重新生成）**：DB 现状字段质量未必高;LLM 严格按 prompt 约束生成更可控。upsert 会覆盖。
  - **副作用**：DB 现存的 9 条 city_visits + 9 条 dice_rolls 关联的城市 ID 不变,继续可用。
- **Q5**：图片生成允许 LLM/IMAGE API key 缺失时**只跑数据流水线**吗?
  - **倾向 A（允许）**：解耦,先把 JSON 数据补齐,图片可以后续异步补。

## 7. 备选方案（如有）

- **方案 A（推荐,本 TODO 已采纳）**：双流水线 + Python 脚本 + skill 文档,工程化交付。
- **方案 B**：纯 LLM 对话手工产出,每城在 ChatGPT/Claude 网页里手动跑 prompt + 人工合并 JSON + 人工下载图。
  - 优点：零代码;缺点：35 城不可重复,扩到 50/100 城时回到原点。
- **方案 C**：放弃扩城,**改为在 roll 端过滤无图城市**（之前给老板的兜底方案）。
  - 优点：30 分钟改完;缺点：不是真正的产品演进,跟"丰满闭环"目标背道而驰。

---

## 8. 待补充：风险登记

- **R1**：LLM 编造方言/史实。**缓解**：prompt 强制保守措辞 + 每城人工校验。
- **R2**：Wikimedia/Unsplash 对部分小城市覆盖度低（如三亚的地标"南山寺"可能找不到合适图）。**缓解**：三级回落到 AI 生图;允许人工指定单条 URL 强制覆盖。
- **R3**：API 成本。Unsplash 免费;Wikimedia 免费;IMAGE_API 按 35 城 × 1 头像 = 35 张 × $0.04 ≈ $1.5（其他类型都先走免费源,仅在失败时回落到 AI）。LLM 35 次调用按 GPT-4o-mini < $0.5。
- **R4**：脚本运行时间。35 城 × ~5 张图 × 2-5 秒下载 ≈ 5-15 分钟。可接受。
- **R5**：人物头像合规。**缓解**：prompt 严格要求插画风,不绘制特定真人面孔,不指代在世名人。
