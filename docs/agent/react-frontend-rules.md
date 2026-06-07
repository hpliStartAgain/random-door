# React 前端编码约束

> 写任何前端文件前必读。父约束：`AGENTS.md` / `CLAUDE.md`。

## 1. 组件与状态

- 一律函数组件 + Hooks，不使用 class 组件。
- 全局状态用 Zustand：`useUserStore`、`useGameStore`、`useCityStore`、`useViewStore`、`useMapStore`、`useToastStore`。
- 组件内只放 UI、交互编排和局部状态；跨视图共享状态放 store。

## 2. 数据请求

- 组件不得直接调用业务 API 的 `fetch` / `axios`。
- 所有业务请求走 `src/api/index.ts`，底层统一使用 `src/api/client.ts`。
- API 类型集中在 `src/api/types.ts`，字段名与 `docs/design/api-contract.md` 对齐。
- 读取本地静态图片转 File 等浏览器工具也应封装在 api/lib 层，避免散落到组件。

## 3. 样式

- Tailwind CSS 为主，保持现有 shadcn/ui 风格约定。
- 避免无约束的大段自定义 CSS；全局动效和主题变量集中在 `src/index.css`。
- 图标按钮优先使用现有 icon 库；已有内联 SVG 可维护，但新增按钮优先用库图标。

## 4. 地图

- 高德地图 SDK 只在 `MapCanvas.tsx` 内封装。
- 其它组件通过 store、props 或回调表达地图意图，禁止直接 `new AMap.*`。

## 5. 图片上传

- 前端做基础类型/大小预校验：jpg / jpeg / png / webp，默认不超过 5MB。
- 后端 `internal/upload/validator.go` 是最终校验点。

## 6. 安全

- 前端不得出现 LLM / IMAGE API Key。
- 仅允许 `VITE_AMAP_KEY`、`VITE_AMAP_SECURITY_CODE` 这类前端公开地图配置。

## 7. 目录约定

```text
src/
  pages/        AdminPage / AssetPage / AchievementPage / GuessChallengePage
  components/   地图、侧栏、抽屉、打卡、评论、个人面板、覆盖层
  api/          axios client、业务 API、TS 类型
  store/        Zustand store
  lib/          纯前端工具
  assets/       前端静态资源引用
```

每个文件职责见 `docs/design/frontend-detailed-design.md`。

## 8. 提交前自检

- [ ] 无组件直接请求业务 API。
- [ ] 无前端硬编码 AI Key。
- [ ] 高德 SDK 调用只在 `MapCanvas.tsx`。
- [ ] 类型与 `api-contract.md` 对齐。
