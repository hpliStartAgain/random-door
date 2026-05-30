"# React 前端编码约束

> 写任何前端文件前必读。父约束：`CLAUDE.md`。

## 1. 组件与状态
- 一律函数组件 + Hooks，不用 class 组件。
- 全局状态用 **Zustand**（useUserStore / useGameStore / useCityStore），不滥用 Context。
- 组件内只放 UI 与局部状态，跨页共享状态放 store。

## 2. 数据请求
- **组件不得直接 fetch/axios**；统一走 src/api/ 下封装函数，底层用 api/client.ts 的 axios 实例（含拦截器：注入 user_id、统一错误处理）。
- 接口的 TS 类型来源于 docs/design/api-contract.md，与后端字段名一致。

## 3. 样式
- Tailwind CSS + shadcn/ui 为主，**不手写大量自定义 CSS**。
- 主题色 / 城市色在 tailwind.config.js 扩展。

## 4. 地图
- 高德地图 SDK **只在 MapCanvas.tsx 内封装**，其它组件通过 props/回调与地图交互，禁止在别处直接 new AMap.*。

## 5. 图片上传
- 走 ImageUploader.tsx，前端做基本类型/大小预校验（jpg/jpeg/png/webp，≤5MB），真正校验在后端。

## 6. 安全
- 前端**绝不出现任何 API Key**（LLM/生图 Key 全在后端）。高德 JS Key 属前端公开 Key，放 VITE_ 环境变量。

## 7. 目录约定
```
src/
  pages/        路由页（8 个）
  components/   复用组件（13 个）
  api/          接口封装 + 类型
  store/        Zustand store（3 个）
```
每个文件职责见 docs/design/frontend-detailed-design.md。

## 8. 提交前自检
- [ ] 无组件直接 fetch；请求都走 api 层
- [ ] 无前端硬编码 AI Key
- [ ] 地图调用只在 MapCanvas
- [ ] 类型与 api-contract.md 对齐"
