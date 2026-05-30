---
status: paused_for_backend
branch: main
owner: antigravity
updated: 2026-05-30 12:10
tier: COMPLEX
---

## 1. 需求理解

核心需求是为“AI城市漫游”产品搭建完整的前端应用。
前端需承担极具表现力的交互与视觉体验：用户将在地图上自由探索或游戏化漫游，与 AI 城市人物进行沉浸式对话，并生成赛博打卡照片，最后收集成就。
应用需基于 React + Vite + Tailwind CSS + shadcn/ui，地图依赖高德 API。

## 2. 设计方案（前端架构与全局色彩风格）

**全局色彩与视觉风格：新中式雅白极简 (Neo-Chinese Light Minimal)**
为了在浅色下依然达到“WOW”的惊艳效果，同时保持文化沉浸感：
*   **布局与质感**：采用浅色高德地图底板（标准或水墨风）。界面大面积留白（如宣纸白 `#F8F9FA`），通过柔和的阴影（Soft Drop Shadow）和极简边框让卡片悬浮于地图之上。
*   **基础色调 (Backgrounds)**：主背景采用「宣纸白」与纯白（White）交织，文字使用黛蓝/墨黑（如 slate-800 或 slate-900）以保证极佳阅读性。
*   **强调色 (Accents)**：
    *   主强调色：**千里青绿 / 天青 (`#0D9488` / teal-600)** —— 点缀交互控件、选中态。
    *   次强调色：**朱砂红 (`#E11D48` / rose-600)** —— 用于关键交互（掷骰子、打卡动作）。
*   **排版与动画**：留白充分的现代排版，配合微秒级的丝滑过渡动效，让页面显得轻盈且通透。

**前端目录与技术栈架构：**
*   **脚手架**：Vite + React + TS。
*   **UI 库**：Tailwind CSS + shadcn/ui 提供极简可定制组件。
*   **状态与路由**：Zustand 管理全局状态 (User/City/Game)；React Router V6 控制单页跳转。
*   **外部集成**：Axios (API 请求), `@amap/amap-jsapi-loader` (高德地图)。

## 3. 阶段划分

- [x] Phase 1: 基础设施搭建 (Vite/Tailwind/路由/状态)。
- [x] Phase 2: 全局色彩体系与核心悬浮布局实现 (Glassmorphism)。
- [x] Phase 3: 沉浸式单页架构重构 (左侧 Sidebar + 右侧地图布局)。
- [x] Phase 4: API 客户端层 (Axios) 与 Zustand Store 前端联调 (使用高质量 Mock 数据闭环)。
- [x] Phase 5: 接入真实高德地图 3D 引擎，实现飞越、摇骰子漫游动画。
- [x] Phase 6: 研发城市详情下钻面板，以及集成 Pannellum 引擎实现真正的 360° 高清全景沉浸。

## 4. 文件级任务

| 文件 | 动作 | 说明 |
|------|------|------|
| frontend/package.json | NEW | 初始化依赖 |
| frontend/tailwind.config.js | NEW | 注入青绿/朱砂红等全局设计令牌 |
| frontend/src/index.css | NEW | 定义基础样式与 Glassmorphism 变量 |
| frontend/src/router.tsx | NEW | 定义 8 大路由结构 |
| frontend/src/store/* | NEW | user, city, game 状态库 |
| frontend/src/api/* | NEW | axios 实例与接口对接 |
| frontend/src/pages/* | NEW | 各业务页面入口 |

## 5. 外部变更（必须显式列出）

- [ ] 新增依赖：`react-router-dom`, `zustand`, `axios`, `lucide-react`, `@amap/amap-jsapi-loader`, `clsx`, `tailwind-merge`。
- [ ] 环境变量：前端根目录新增 `.env`（用于存放 `VITE_AMAP_KEY`）。

## 6. 待确认问题（@老板）

- Q1: “新中式暗夜赛博 (Neo-Chinese Cyber Dark)”这套以深色地图为底、毛玻璃卡片悬浮、青绿和朱砂红点缀的视觉方案是否符合您心目中惊艳的预期？
- Q2: 初始化前端工程将直接在 `frontend` 目录下进行（会清理空目录并重新生成）。确认可直接执行构建命令？

## 7. 备选方案（如有）

- 方案 A（当前推荐）：全局暗黑赛博风，强调氛围沉浸感和高级感。
- 方案 B：高亮极简风 (Light & Clean)，类似 Apple Maps / Airbnb 风格，以大面积白色留白和卡片投影为主。但这可能在“赛博打卡”这种概念下显得不够炫酷。
