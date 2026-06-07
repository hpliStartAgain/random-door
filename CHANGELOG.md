# Changelog

## Unreleased

### Changed

- 统一项目名称为“任意门”，文档不再以旧验证阶段表述当前状态。
- 以当前代码为事实标准重写 README、产品文档、架构文档、前后端详设和 agent 约束。
- 新增 `TODO.md` 项目级待办，替代一次性修复任务记录。

### Removed

- 删除历史对标改造清单、目标任务记录和已入库 Python 缓存文件。

## [1.0.0] - 2026-06-07

### Added

- 后台 CMS 扩展：城市、标签、地标、美食、人物、成就、勋章和媒体资产维护。
- 用户资产与个人资料：足迹城市、打卡海报、成就进度、匿名资料编辑、注册/登录保留足迹。
- 猜城市挑战：截图/图片文案、匿名挑战链接、好友答题。
- 城市发现增强：区域/标签/搜索筛选，地图 marker 与列表同步。
- 地标坐标与声景：地标 marker、`soundscape_url`、坐标回填命令。
- 角色元数据：`role_title`、`life_span`、`intro_quote`，前端人物卡片可展示更完整叙事。
- 狐狸视觉与任意门品牌动效，`RandomCityModal` 替代旧 DiceConsole 命名。

### Changed

- seed 导入改为显式 `SEED_MODE=off|bootstrap|sync`，默认不写库；数据库成为内容事实源。
- seed 扩展到 35 座精选城市，并将城市标签规范化为中文。
- seed 重跑时保留后台上传媒体，避免覆盖运营内容。
- Docker Compose 改为外部 MySQL，Caddy 等待 app healthcheck 后启动。
- AI 生图改为任务化配置，增加 worker、每日限额和重试参数。

### Fixed

- 修复成就墙空数据白屏。
- 修复后台假鉴权，必须校验 `ADMIN_TOKEN` 后进入。
- 将不可用的全景依赖降级为图片风光浏览，保持打卡/文案链路可用。
- 修复搜索与地图不同步、成就进度溢出、城市数量文案不一致等体验问题。
- 修复打卡原子性、错误分类、图片重试逻辑和若干死页面问题。

## [0.2.0] - 2026-06-06

### Added

- 管理端媒体上传接口与 ADMIN_TOKEN 保护。
- 35 城 seed、静态图片资产、AI worker 配置和容器健康检查。
- Docker 构建中自动创建 static/uploads 目录。

### Changed

- 从嵌入式 MySQL 切换为外部 MySQL 实例。
- 使用 `docker-compose` 命令兼容更多部署环境。
- 移除已提交的 `.env` 敏感文件，更新环境变量模板。

## [0.1.0] - 2026-05-30

### Added

- React + Vite + Zustand + Tailwind 前端基础架构。
- Go + Gin + GORM 后端基础架构。
- 匿名用户、城市列表/详情、自由访问、随机漫游、AI 对话、打卡、成就墙主链路。
- 地理算法：Haversine、8 方向、6 距离档、目标点计算、最近城市匹配。
- seed 校验与幂等导入。
- HTTP 兼容 UUID 生成器，支持非安全上下文回退。

### Changed

- 升级 Go 基础镜像到 1.22。
- 配置 Go / npm 国内镜像与前端依赖安装兼容参数。
