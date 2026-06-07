# yourwonder.us 对标改造清单

> 目的：基于 OpenCLI 对 `https://www.yourwonder.us/` 的实测探索，提炼本项目可借鉴的产品与工程改造项，并落到文件增删改级别。
> 范围：黑客松 MVP v0.1，不变更技术栈，不引入 Redis / Kafka / Elasticsearch / 向量库 / 微服务网关，AI Key 仍只在后端。

## 1. OpenCLI 探索结论

本次使用 `opencli browser` 观察了首页、地标详情、人物对话入口、街景区域、猜位置分享弹层、个人面板与资源加载。

可借鉴点：

1. 首页即是产品体验，不是营销页：左侧发现列表，右侧地图，顶部只保留品牌与当前上下文。
2. 发现页用强信息密度建立探索感：`50 landmarks · 149 historical figures · 50 nearby` 这类统计让内容规模可感知。
3. 列表卡片比文字更先建立兴趣：大图、地名、国家/城市、Nearby 标签、图标化 marker 同时出现。
4. 地标详情把沉浸画面和人物入口放在同一状态：详情左侧是地标介绍与历史人物，右侧进入 Street View。
5. 历史人物卡有一句可引用的话、身份、年份和明确 CTA，比单纯“人物名 + 对话按钮”更有叙事感。
6. “Where am I?” 分享弹层很完整：截图预览、Challenge Friends、Copy Challenge、Save Photo、Copy Image、Post to X 在同一链路里。
7. 个人面板低门槛：不要求注册，也能看到 visited places，并可补姓名、年龄、国家。
8. 资源按场景加载：首页预取地标图，进入详情后再加载人物图、soundscape 音频和街景资源。

不建议照搬：

1. 不照搬世界地标范围。本项目核心是中国城市文化，内容仍控制在 35 个精选城市。
2. 不照搬 Snap Camera Kit 这类重前端 AR 依赖。MVP 目标服务器 2C2G，且已有 AI 打卡链路。
3. 不把任何 AI Key 放前端。yourwonder 的 Google Maps embed key 可见，但本项目只允许前端出现公开地图 Key，不允许 LLM / 生图 Key。
4. 不让 AI 人物声称真实复活；人物对话继续遵循现有 prompt 合规约束。

## 2. 本项目已有基础

这些能力已经存在，应增量增强：

1. `frontend/src/App.tsx` 已是常驻地图 + 侧栏 + 街景浮层 + 右抽屉结构。
2. `frontend/src/components/MapCanvas.tsx` 已按 `useCityStore.filteredCities()` 重绘地图点，搜索和地图天然可共享同一状态源。
3. `frontend/src/components/StreetViewCanvas.tsx` 已支持“截图生成文案”和“用当前视角作为赛博打卡参考图”。
4. `docs/design/api-contract.md` 已有 `POST /api/guess/caption`，可支撑分享文案生成。
5. `frontend/src/components/CheckinFlow.tsx` 已支持 `sceneFile`，与 yourwonder 的截图打卡方向一致。
6. `frontend/src/pages/AssetPage.tsx` 已展示 visited cities、posters、achievement progress，可改造成更轻量的 profile 面板。
7. `frontend/src/components/RightDrawer.tsx` 已有 AI 对话、建议问题和评论区。

## 3. 改造包 A：发现页信息密度与筛选

目标：借鉴 yourwonder 的发现列表信息密度，但保持本项目的“中国城市 + 地图探索 + 任意门”定位。

### A1. 前端筛选与统计增强，不改接口

- [x] 新增 `frontend/src/lib/cityFilters.ts`
  - 定义区域筛选映射，例如 `north_china`、`jiangnan`、`southwest`、`coastal`、`food_city` 等现有 tag 分组。
  - 提供 `getRegionOptions(cities)`、`filterCities(cities, query, region, tag)`。
- [x] 修改 `frontend/src/store/useCityStore.ts`
  - 新增 `activeRegion`、`activeTag`、`setActiveRegion`、`setActiveTag`、`resetFilters`。
  - `filteredCities()` 同时考虑搜索词、区域、标签。
- [x] 修改 `frontend/src/components/layout/Sidebar.tsx`
  - 在搜索框下增加 tag/区域 chips。
  - 标题统计从单一城市数改为：总城市数、当前筛选命中数、已访问数。
  - 空状态保留“清空搜索”，同时清空 region/tag。
- [x] 修改 `frontend/src/components/MapCanvas.tsx`
  - 保持使用 `filteredCities()`，确保筛选 chips 与地图点同步。
  - marker 增加“已访问 / 当前目标 / 普通城市”三种视觉状态。
- [x] 修改 `frontend/src/index.css`
  - 降低圆角堆叠，chips 与 marker hover 状态统一到现有色板。

验收：

- 搜索、区域、标签任一变化后，左侧卡片与地图 marker 数量一致。
- 空状态可一键清空所有筛选。
- 不新增接口，不改数据库。

### A2. 城市列表统计字段增强，需要接口改动

如果要做到 yourwonder 式的 `landmarks · figures · nearby` 统计，走这一项。

- [x] 修改 `docs/design/api-contract.md`
  - `GET /api/cities` 的城市项增加 `landmark_count`、`food_count`、`character_count`。
- [x] 修改 `docs/design/frontend-detailed-design.md`
  - 更新 `CityListItem` 与发现页职责。
- [x] 修改 `docs/arch/backend-detailed-design.md`
  - 更新 `city_repo.go` 与 `city_service.go` 的聚合职责。
- [x] 修改 `backend/internal/model/city.go`
  - 仅在响应 DTO 中增加统计字段，不建议把统计字段落库。
- [x] 修改 `backend/internal/repository/city_repo.go`
  - 增加按城市聚合 landmarks / foods / characters 数量的查询。
- [x] 修改 `backend/internal/service/city_service.go`
  - List 聚合统计字段。
- [x] 修改 `backend/internal/api/city_handler.go`
  - 返回新字段。
- [x] 修改 `frontend/src/api/types.ts`
  - `CityListItem` 增加三个 count 字段。
- [x] 修改 `frontend/src/components/layout/Sidebar.tsx`
  - 展示当前筛选集合内的总地标、总美食、总人物。
- [x] 修改 `backend/internal/service/city_service_test.go`
  - 覆盖统计字段。

验收：

- `GET /api/cities` 对旧字段兼容，新增字段有默认 0。
- 前端缺少 count 字段时不崩溃，方便旧后端联调。

## 4. 改造包 B：城市详情与人物叙事增强

目标：让城市详情从“资料卡”更接近“可探索的文化现场”。

### B1. 不改表的前端排版优化

- [x] 修改 `frontend/src/components/CityDetailPanel.tsx`
  - 地标卡 CTA 从“查看图片”改为“进入风光视角”或“进入地标视角”。
  - 历史人物卡改成 quote-like 卡片：头像、姓名、类型、方言风格、对话 CTA。
  - 沉浸体验区减少说明文案，突出“风光视角”“赛博打卡”两个动作。
  - 修复当前文件里的异常字符显示：`<span className="text-3xl ...">�️</span>`。
- [x] 修改 `frontend/src/components/RightDrawer.tsx`
  - 人物对话开场文案不再泛化为“汝好”，根据 `character_type` 和 `dialect_style` 生成更稳妥的开场。
  - 建议问题按城市/人物上下文展示，保留通用兜底。
- [x] 修改 `frontend/src/index.css`
  - `panel` 圆角从 `rounded-2xl` 收敛到更紧凑的半径，避免卡片堆叠过重。

验收：

- 城市详情首屏能看到城市名、简介、一个核心行动。
- 人物卡不声称真实复活，不展示后端 prompt/persona。

### B2. 增加人物 quote / 年份 / 身份字段，需要表结构改动

- [x] 修改 `docs/design/api-contract.md`
  - `GET /api/cities/{id}` 的 `characters[]` 增加 `role_title`、`life_span`、`intro_quote`。
- [x] 修改 `docs/design/database-detailed-design.md`
  - `characters` 表增加 `role_title VARCHAR(128)`、`life_span VARCHAR(64)`、`intro_quote VARCHAR(255)`。
- [x] 修改 `backend/migrations/schema.sql`
  - 同步新增字段。注意已有数据库升级需另配迁移说明，MVP seed 初始化可直接建表。
- [x] 修改 `backend/internal/model/character.go`
  - 增加字段与 JSON tag，仍保持 `Persona` / `Prompt` 不下发。
- [x] 修改 `backend/internal/repository/city_repo.go`
  - 查询并返回新增字段。
- [x] 修改 `backend/internal/seed/seed.go`
  - seed 校验允许并校验新增字段长度。
- [x] 修改 `backend/data/seed/cities.json`
  - 每个城市至少 1 个人物补 `role_title`、`life_span`、`intro_quote`。
- [x] 修改 `frontend/src/api/types.ts`
  - `Character` 增加新增字段。
- [x] 修改 `frontend/src/components/CityDetailPanel.tsx`
  - 按 yourwonder 的人物卡结构展示身份、年份、引用语。
- [x] 修改 `backend/internal/model/schema_contract_test.go`
  - 覆盖新增字段。
- [x] 修改 `backend/internal/seed/seed_test.go`
  - 覆盖新增字段校验。

验收：

- 人物卡能显示“身份 + 年份 + 引用语 + 对话按钮”。
- 引用语必须是角色口吻的引导文案，不写成确定性史实断言。

## 5. 改造包 C：街景猜一猜与分享链路

目标：把现有 `StreetViewCanvas` 的“截图生成文案”升级为完整的猜位置分享体验。

### C1. 前端本地分享弹层，不新增后端

- [x] 新增 `frontend/src/components/GuessChallengeModal.tsx`
  - 展示截图预览、标题“猜猜我在哪”、生成文案、复制文案、保存图片、复制图片。
  - 使用现有 `api.generateGuessCaption()`。
- [x] 新增 `frontend/src/lib/shareImage.ts`
  - 封装 canvas/blob/dataURL 转换。
  - 提供 `downloadImage()`、`copyImageToClipboard()`。
- [x] 修改 `frontend/src/components/StreetViewCanvas.tsx`
  - 移出当前内联猜一猜面板的重复逻辑。
  - 点击“猜一猜”后打开 `GuessChallengeModal`。
  - 保留“生成赛博打卡”，继续传 `sceneFile` 给 `CheckinFlow`。
- [x] 修改 `frontend/src/api/types.ts`
  - 如需要，补充 `GuessCaptionResponse` 的可选 `share_title` 字段，但优先不改接口。
- [x] 修改 `frontend/src/index.css`
  - 增加分享弹层的稳定尺寸、移动端最大宽度和截图区域 aspect-ratio。

验收：

- 进入街景后可生成截图文案。
- 可复制文案，可保存图片。
- 浏览器不支持复制图片时降级到保存图片，不报错阻断。

### C2. 真实好友挑战链接，需要新增接口和表

只有在需要“朋友打开链接并提交猜测”时做这一项。

- [x] 修改 `docs/design/api-contract.md`
  - 新增 `POST /api/guess/challenges`。
  - 新增 `GET /api/guess/challenges/{code}`。
  - 新增 `POST /api/guess/challenges/{code}/answers`。
- [x] 修改 `docs/design/database-detailed-design.md`
  - 新增 `guess_challenges` 表：`code`、`user_id`、`city_id`、`target_name`、`image_url`、`caption`、`created_at`、`expires_at`。
  - 新增 `guess_answers` 表：`challenge_code`、`answer_text`、`is_correct`、`created_at`。
- [x] 修改 `backend/migrations/schema.sql`
  - 新增两张表和索引。
- [x] 新增 `backend/internal/model/guess_challenge.go`
- [x] 新增 `backend/internal/model/guess_answer.go`
- [x] 新增 `backend/internal/repository/guess_challenge_repo.go`
- [x] 新增 `backend/internal/service/guess_challenge_service.go`
- [x] 新增 `backend/internal/api/guess_challenge_handler.go`
- [x] 修改 `backend/internal/api/router.go`
  - 注册挑战相关路由。
- [x] 确认复用 `backend/internal/upload/storage.go`
  - 如分享截图需要落本地，增加 `uploads/guess` 子目录。
- [x] 确认 `.gitignore`
  - 确认 `uploads/` 仍不提交。
- [x] 修改 `frontend/src/api/types.ts`
  - 增加 challenge 请求/响应类型。
- [x] 修改 `frontend/src/api/index.ts`
  - 增加 create/get/answer challenge API。
- [x] 新增 `frontend/src/pages/GuessChallengePage.tsx`
  - 展示截图，输入猜测，显示正确/接近/错误反馈。
- [x] 修改 `frontend/src/App.tsx`
  - 如果现有项目继续无路由，可用 `useViewStore` 管理；如引入 route，必须与前端设计文档同步。
- [x] 修改 `frontend/src/components/GuessChallengeModal.tsx`
  - `Challenge Friends` 生成链接并复制。
- [x] 新增 `backend/internal/service/guess_challenge_service_test.go`
- [x] 新增 `backend/internal/api/guess_challenge_handler_test.go`

验收：

- 未登录用户也能创建挑战。
- challenge code 不暴露 user_id。
- 过期挑战返回 404 或明确过期错误。

## 6. 改造包 D：声景与沉浸资产

目标：借鉴 yourwonder 的 soundscape，让城市/地标不只有图片，也有“环境感”。

### D1. 城市或地标声景字段

- [x] 修改 `docs/design/api-contract.md`
  - `GET /api/cities/{id}` 的 city 或 `landmarks[]` 增加 `soundscape_url`。
  - 建议优先放在 landmark 上，城市可做兜底。
- [x] 修改 `docs/design/database-detailed-design.md`
  - `landmarks` 增加 `soundscape_url VARCHAR(512)`。
- [x] 修改 `backend/migrations/schema.sql`
  - 同步新增字段。
- [x] 修改 `backend/internal/model/landmark.go`
  - 增加 `SoundscapeURL *string`。
- [x] 修改 `backend/internal/seed/seed.go`
  - 校验声景 URL 只能是 `/static/soundscapes/...` 或空。
- [x] 修改 `backend/data/seed/cities.json`
  - 先给 5 个演示城市的核心地标补声景路径，其余可为空。
- [x] 新增 `backend/static/soundscapes/README.md`
  - 说明音频来源、授权、压缩目标和命名规则。
- [x] 新增 `frontend/src/components/SoundscapeControl.tsx`
  - 播放/暂停、音量、当前地标名，默认静音不自动播放。
- [x] 修改 `frontend/src/components/StreetViewCanvas.tsx`
  - 当前 street target 有 `soundscape_url` 时显示声景控制。
- [x] 修改 `frontend/src/components/CityDetailPanel.tsx`
  - 点击地标进入街景时把 `soundscape_url` 传入 `streetTarget`。
- [x] 修改 `frontend/src/store/useViewStore.ts`
  - `StreetTarget` 增加 `soundscape_url`。
- [x] 修改 `frontend/src/api/types.ts`
  - `Landmark` 增加 `soundscape_url`。

验收：

- 声景不自动播放，必须用户点击。
- 音频文件体积受控，避免拖慢 2C2G 演示。

## 7. 改造包 E：个人足迹面板

目标：将现有“我的资产”从独立页面增强为 yourwonder 式轻量个人面板，降低用户查看足迹的成本。

### E1. 前端先重组，不新增用户字段

- [x] 新增 `frontend/src/components/ProfilePanel.tsx`
  - 展示匿名头像、当前昵称或“未设置昵称”、走过城市、打卡海报、成就进度摘要。
  - 数据复用 `api.getUserAssets()`。
- [x] 修改 `frontend/src/components/layout/Navbar.tsx`
  - 右上角改为图标按钮打开 ProfilePanel。
  - 保留“成就墙”入口，减少顶部文字导航拥挤。
- [x] 修改 `frontend/src/App.tsx`
  - 挂载 ProfilePanel。
- [x] 修改 `frontend/src/store/useViewStore.ts`
  - 增加 `profileOpen` 或将 profile 纳入 `currentView`。
- [x] 修改 `frontend/src/pages/AssetPage.tsx`
  - 保留完整资产页，但复用 ProfilePanel 的列表组件，避免两套 UI。
- [x] 新增 `frontend/src/components/ProfileVisitedList.tsx`
  - 资产页和 ProfilePanel 共同复用。
- [x] 新增 `frontend/src/components/ProfilePosterGrid.tsx`
  - 资产页和 ProfilePanel 共同复用。

验收：

- 顶部个人入口一键打开，不离开当前地图/街景上下文。
- 资产页仍可作为完整视图存在。

### E2. 可编辑姓名、年龄、地区，需要接口改动

- [x] 修改 `docs/design/api-contract.md`
  - 新增 `GET /api/users/{id}/profile`。
  - 新增 `PATCH /api/users/{id}/profile`，允许 `nickname`、`age`、`home_region`。
  - `GET /api/users/{id}/assets` 可内嵌 profile 摘要，或前端并发读取 profile。
- [x] 修改 `docs/design/database-detailed-design.md`
  - `users` 增加 `age INT`、`home_region VARCHAR(64)`。
- [x] 修改 `backend/migrations/schema.sql`
  - 同步新增字段。
- [x] 修改 `backend/internal/model/user.go`
  - 增加字段。
- [x] 修改 `backend/internal/repository/user_repo.go`
  - 增加 `FindByID`、`UpdateProfile`。
- [x] 新增 `backend/internal/service/user_service.go`
  - 处理 profile 读写校验。
- [x] 新增 `backend/internal/api/user_handler.go`
  - 处理 profile 接口。
- [x] 修改 `backend/internal/api/router.go`
  - 注册 profile 接口。
- [x] 修改 `frontend/src/api/types.ts`
  - 增加 `UserProfileResponse`。
- [x] 修改 `frontend/src/api/index.ts`
  - 增加 `getUserProfile`、`updateUserProfile`。
- [x] 修改 `frontend/src/components/ProfilePanel.tsx`
  - 增加内联编辑。
- [x] 新增 `backend/internal/service/user_service_test.go`
- [x] 新增 `backend/internal/api/user_handler_test.go`

验收：

- 昵称 1 到 64 字，年龄合理范围校验。
- 匿名用户可编辑，不需要注册登录。

## 8. 不删除项

本轮不建议删除任何现有业务文件。理由：

1. `StreetViewCanvas.tsx`、`CheckinFlow.tsx`、`AssetPage.tsx` 已承载主链路能力，适合拆分组件和增量增强。
2. `RightDrawer.tsx` 已有对话和评论能力，不应改成独立页面导致主体验割裂。
3. `MapCanvas.tsx` 已是高德 SDK 唯一封装点，必须保留这个边界。

允许的“删除”只限于文件内部重复 JSX 或异常字符，不删除文件：

- [x] 修改 `frontend/src/components/StreetViewCanvas.tsx`：将内联分享面板 JSX 抽到 `GuessChallengeModal.tsx`。
- [x] 修改 `frontend/src/components/CityDetailPanel.tsx`：删除异常字符显示和过度说明文案。

## 9. 推荐实施顺序

1. 先做 A1、B1、C1、E1：全是前端增量，风险低，不破坏接口。
2. 再做 D1：需要 schema 和 seed，但可只给少量演示城市补声景。
3. 再做 A2、B2、E2：涉及接口和表，必须先改设计文档和契约。
4. 最后做 C2：真实挑战链接是增长玩法，范围最大，不应阻塞 MVP 演示。

## 10. 验证清单

每个改造包完成后至少执行：

- [x] `npm run build`，工作目录 `frontend/`。
- [x] `go test ./...`，工作目录 `backend/`。
- [x] 如改 seed：`python scripts/build_seed.py validate`。
- [ ] 手动走通：进入首页、搜索/筛选城市、点击城市、进入地标视角、生成猜一猜文案、生成赛博打卡、查看个人足迹。
- [x] 如改接口：先确认 `docs/design/api-contract.md` 已更新，再确认前后端字段一致。

