# 前端详细设计

## 1. 应用形态

当前前端是单页全屏工作台，不是多路由页面栈。`App.tsx` 根据 URL 和 Zustand 状态组合以下区域：

```text
/?guess=<code> -> GuessChallengePage + Toast

默认入口
  Navbar
  Sidebar
  MapCanvas
  StreetViewCanvas (canvasMode=street 时覆盖地图)
  RandomCityModal (currentView=GAME_DICE)
  RightDrawer
  AchievementPage / AssetPage 覆盖层
  WelcomeOverlay
  AdminPage
  ProfilePanel
  Toast
```

## 2. pages

| 文件 | 职责 | 主要接口 |
|---|---|---|
| `AdminPage.tsx` | ADMIN_TOKEN 门控后台 CMS，维护城市/标签/地标/美食/人物/成就/媒体。 | `/admin/*` |
| `AchievementPage.tsx` | 成就墙，展示已解锁、未解锁、进度。 | `GET /users/{id}/achievements` |
| `AssetPage.tsx` | 用户资产页，展示足迹城市、海报、成就进度。 | `GET /users/{id}/assets` |
| `GuessChallengePage.tsx` | 分享挑战落地页，提交猜测并显示结果。 | `GET/POST /guess/challenges/*` |

## 3. components

| 文件 | 职责 |
|---|---|
| `MapCanvas.tsx` | 唯一封装高德 SDK；加载城市/地标 marker，处理点击、飞行和降落动画。 |
| `StreetViewCanvas.tsx` | 图片沉浸视图，提供声景、猜一猜、赛博打卡入口。 |
| `layout/Navbar.tsx` | 品牌、我的足迹、成就墙、后台入口。 |
| `layout/Sidebar.tsx` | 城市列表、搜索/区域/标签筛选、任意门入口、城市详情容器。 |
| `CityDetailPanel.tsx` | 城市详情、地标/美食/人物/方言、打卡和成就联动。 |
| `RightDrawer.tsx` | 人物 AI 对话抽屉。 |
| `CheckinFlow.tsx` | 选择地标、上传自拍、创建生图任务、轮询/重试、确认打卡。 |
| `CheckinPoster.tsx` | Canvas 合成并下载打卡海报。 |
| `CommentThread.tsx` | 评论列表与创建。 |
| `GuessChallengeModal.tsx` | 创建挑战链接，复制分享文案。 |
| `SoundscapeControl.tsx` | 用户点击后播放地标声景。 |
| `ProfilePanel.tsx` | 右侧个人足迹面板，资料编辑、城市/海报/成就标签页。 |
| `ProfileVisitedList.tsx` | 足迹城市列表复用组件。 |
| `ProfilePosterGrid.tsx` | 打卡海报网格复用组件。 |
| `Toast.tsx` | 全局 toast 展示。 |

### overlays

| 文件 | 职责 |
|---|---|
| `WelcomeOverlay.tsx` / `WelcomeCard.tsx` | 首屏欢迎、模式入口、注册/登录。 |
| `RandomCityModal.tsx` | 任意门随机漫游弹窗。 |
| `AchievementUnlock.tsx` | 成就解锁庆祝层。 |
| `CityDrawer.tsx` / `DiceConsole.tsx` | 旧组件仍在仓库中，后续清理见 `TODO.md`。 |

## 4. api

| 文件 | 职责 |
|---|---|
| `client.ts` | axios 实例，baseURL `/api`，注入 `X-User-Id`，统一错误结构。 |
| `index.ts` | 所有业务 API 函数。 |
| `types.ts` | 与 API 契约对齐的 TypeScript 类型。 |

组件不得直接请求业务 API。

## 5. store

| store | 主要状态 |
|---|---|
| `useUserStore` | userId、anonymousId、username、nickname、currentCityId、注册/登录/登出。 |
| `useCityStore` | cities、cityCache、searchQuery、activeRegion、activeTag、filteredCities。 |
| `useGameStore` | 起点、上次 roll、方向、距离、目标城市、rolling 状态。 |
| `useMapStore` | map 实例相关动作，例如 flyTo。 |
| `useViewStore` | currentView、activeCityId、canvasMode、drawer、rollPhase、profileOpen。 |
| `useToastStore` | toast 队列。 |

## 6. lib / assets

| 文件 | 职责 |
|---|---|
| `lib/cityFilters.ts` | 城市区域、标签与搜索筛选。 |
| `lib/shareImage.ts` | 分享/截图相关纯前端工具。 |
| `assets/foxImages.ts` | 狐狸视觉资源映射。 |

## 7. 安全与边界

- 前端只允许高德公开 Key；不得出现 LLM / IMAGE Key。
- 高德 SDK 只在 `MapCanvas.tsx` 内使用。
- 上传预校验只做体验优化，最终校验在后端。
- 与后端字段不一致时，先改 `api-contract.md`，再改 `types.ts` 和调用点。
