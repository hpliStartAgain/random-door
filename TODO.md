---
status: done
branch: master
owner: devin
updated: 2026-06-07 00:00
tier: COMPLEX
---

## 1. 需求理解

修复线上演示阻断问题，分 P0/P1/P2 三级：
- P0（演示阻断）：成就墙白屏、后台假鉴权、全景加载失败
- P1（优先优化）：城市数量文案不一致、成就进度溢出、搜索不同步地图、语义化/可访问性
- P2（暂不实现）：图片性能、赛博打卡前置引导、高德 Key 域名限制（需运营控制台操作）

## 2. 设计方案

### P0-1 成就墙白屏
- 根因：Go `var unlocked []T` 在无数据时为 `nil`，JSON 序列化为 `null`；前端 `data?.unlocked.map(...)` 在 `null` 上崩溃。
- 后端修复：初始化为 `make([]T, 0)`。
- 前端修复：`(data?.unlocked ?? []).map(...)` 兜底；无成就时展示占位文案。

### P0-2 后台假鉴权
- 根因：`handleAuth` 直接 `setAuthed(true)`，不等 API 返回。
- 修复：改为 async，调用 `api.adminCoverage(token)` 成功才 `setAuthed(true)`，同时复用返回的 coverage；失败弹错误 toast 停留在 token 输入层。

### P0-3 全景不可用
- 根因：Pannellum 需要等距柱状全景图（equirectangular），而资源只有普通封面 JPEG；CORS 问题会导致 "The file could not be accessed"。
- 修复：移除 Pannellum，StreetViewCanvas 用 `<img>` 全屏展示封面图（沉浸体验不丢失，打卡/文案功能保留）；CityDetailPanel 中"全景"按钮改文案为"查看图片"、"进入 3D 街景"改为"查看城市风光"。

### P1-4 城市数量文案不一致
- WelcomeOverlay 硬编码 "35座城市" → 改为从 useCityStore 动态获取；Sidebar "处名胜" → "座城市"。

### P1-6 成就进度溢出
- 进度条高度 `clamp` 到 100%，progress current `clamp` 到 target，已完成时显示"已完成"。

### P1-7 搜索不同步地图
- `useCityStore` 增加 `searchQuery / setSearchQuery`，计算并导出 `filteredCities`。
- Sidebar 使用 store 的状态（去掉本地 state）。
- MapCanvas 改用 `filteredCities` 重绘 markers；无结果时清空 markers 并保留用户位置。

### P1-8 语义化
- CityDetailPanel 历史人物 / 美食卡片 `div[onClick]` → `button`。

## 3. 阶段划分

- [x] Phase 1: P0 后端修复（成就 null 问题）
- [x] Phase 2: P0 前端修复（成就墙兜底 + 假鉴权 + 全景降级）
- [x] Phase 3: P1 前端修复（文案、进度 clamp、搜索同步地图、语义化）

## 4. 文件级任务

| 文件 | 动作 | 说明 |
|------|------|------|
| backend/internal/service/achievement_service.go | MODIFY | 初始化 unlocked/locked/progress 为 make([]T,0) |
| frontend/src/pages/AchievementPage.tsx | MODIFY | ?? [] 兜底；进度 clamp；无成就占位 |
| frontend/src/pages/AdminPage.tsx | MODIFY | handleAuth 改 async，先验证 token |
| frontend/src/components/StreetViewCanvas.tsx | MODIFY | 移除 Pannellum，用 img 全屏背景 |
| frontend/src/components/CityDetailPanel.tsx | MODIFY | 全景按钮文案；人物/食物 div→button |
| frontend/src/store/useCityStore.ts | MODIFY | 增加 searchQuery/setSearchQuery/filteredCities |
| frontend/src/components/layout/Sidebar.tsx | MODIFY | 使用 store searchQuery，文案更新 |
| frontend/src/components/MapCanvas.tsx | MODIFY | 使用 filteredCities |
| frontend/src/components/overlays/WelcomeOverlay.tsx | MODIFY | "35座城市" 改动态 |

## 5. 外部变更

无表结构变更，无新依赖，无环境变量变更。

## 6. 待确认问题

- Q1: P2（图片 WebP、懒加载、高德域名限制）是否纳入本次？→ 暂缓，本次只做 P0+P1。
- Q2: 全景功能是否有计划补全（添加真正等距柱状全景图字段）？→ 暂按降级处理，保留入口但改文案。

## 7. 备选方案

- 搜索同步地图：方案 A（当前选取）：store 全局 searchQuery；方案 B：将 filteredCities 作为 prop 传给 MapCanvas（需改 App.tsx，更复杂）。
