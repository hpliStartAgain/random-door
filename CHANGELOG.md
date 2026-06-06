# Changelog

## [1.0.0] - 2026-06-06

### Added
- **WelcomeOverlay**：全屏品牌开场页，双入口（自由探索 / 随机漫游）+ 微动效
- **DiceConsole**：3D 骰子滚动动画 + 罗盘方向 + 距离展示 + 揭晓仪式（Phase 1 wow①）
- **MapCanvas**：掷骰后飞行轨迹折线 + 电影级 flyTo + 降落标记动效
- **Sidebar**：城市卡片渲染 `cover_image_url` 真实图片，渐变兜底；搜索过滤
- **RightDrawer**：打字机逐字效果 + 角色专属开场白 + 推荐问题 chips + 头像支持（Phase 2 wow②）
- **CheckinFlow**：三步打卡向导（选地标→上传自拍→AI生成进度→确认打卡）
- **CheckinPoster**：Canvas 合成品牌水印海报，一键下载；可传播（Phase 2 wow③）
- **AchievementUnlock**：全屏成就解锁庆祝，CSS 粒子动画，6s 自动关闭（Phase 2 wow④）
- **Toast + useToastStore**：全局 Toast 通知系统，替换所有 `alert()`
- **CityDetailPanel**：集成 CheckinFlow 滑入层 + 方言速记卡 + AchievementUnlock 联动
- **AdminPage（后台管理）**：Token 门控 + 城市/地标/人物/美食 批量上传文件或粘贴外链 URL
- **后端 admin API**：`POST /api/admin/*/image`（文件上传）+ `PATCH /api/admin/*/image`（URL绑定），ADMIN_TOKEN 保护
- **后端加固**：错误码常量化；`service.ErrConflict` / `ErrPermission`；chat 消息长度校验（最大 500 字符）
- **docker-compose**：app + caddy 双 healthcheck；caddy 等待 app healthy 后启动
- **docs/product/demo-script.md**：<10 分钟评委演示脚本

### Changed
- `useViewStore`：新增 `hasEntered`、`enter`、漫游动画状态
- `useGameStore`：暴露 `fromPoint`、`direction`、`distance_km` 供地图动画消费
- `.env.example`：新增 `ADMIN_TOKEN` 说明

### Fixed
- 修复 JSX 中 IIFE 过滤引起的语法错误（Sidebar 搜索逻辑提取为变量）
- 修复 `CityDetailPanel` 内联打卡逻辑替换为 CheckinFlow 解耦

## [0.1.0] - 2026-05-30
### Added
- **Frontend-Backend Integration**: fully connected the React frontend with the Go Gin backend.
- **Game Engine**: wired up the Dice Roller (`useGameStore.roll`) with backend API (`/api/game/roll`) to perform remote geospatial random rolling.
- **AI Chat Module**: connected the `RightDrawer` chat UI with the `/api/chat` endpoint.
- **Cyber Check-in**: added image upload and `api.generateImage` and `api.createCheckin` call to `CityDetailPanel.tsx` for Cyberpunk check-ins.
- **Achievement Wall**: connected `AchievementPage.tsx` with `/api/users/:id/achievements` and integrated it as an overlay via `Navbar`.
- **Backend Architecture**: implemented Go + Gin + GORM backend supporting MySQL and dynamic seeding.
- **Frontend Architecture**: built React + Vite + Zustand + TailwindCSS app with AMap 3D and Pannellum 360 viewer.

### Fixed
- Fixed TypeScript errors related to mismatched fields between Mock data (`city.figures`) and real API (`city.characters`).
- Fixed missing arguments in `roll()` calls within `Sidebar.tsx` and `DiceConsole.tsx`.
