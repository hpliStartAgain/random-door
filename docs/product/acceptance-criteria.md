# MVP 验收标准（acceptance-criteria.md）

> 对应概要设计 27 章。既是开发完成定义(DoD)，也是 PR 审查与现场演示的 checklist。

## 验收清单（20 条，逐条勾选）
| # | 验收条件 | 关联接口/页面 | 通过 |
|---|---|---|---|
| 1 | 用户可打开网页进入首页 | HomePage | ☐ |
| 2 | 可看到自由探索与游戏互动两个入口 | ModeSelectPage | ☐ |
| 3 | 可进入自由探索模式 | FreeExplorePage | ☐ |
| 4 | 自由探索下可在地图点击城市 | GET /cities、MapCanvas | ☐ |
| 5 | 点击城市后进入城市详情页 | POST /visits/free、GET /cities/{id} | ☐ |
| 6 | 可进入游戏互动模式 | GameModePage | ☐ |
| 7 | 游戏模式可获取当前位置或用默认位置 | POST /game/init | ☐ |
| 8 | 可点击掷骰按钮 | DicePanel | ☐ |
| 9 | 系统生成随机方向与随机距离 | geo 算法 | ☐ |
| 10 | 系统计算并返回目标城市 | POST /game/roll | ☐ |
| 11 | 地图展示移动到目标城市的过程 | MapCanvas 动画 | ☐ |
| 12 | 城市详情展示地标/美食/人物/方言 | GET /cities/{id} | ☐ |
| 13 | 可与至少一个城市人物 AI 对话 | POST /chat | ☐ |
| 14 | 可上传照片生成赛博游客照 | POST /checkin/generate-image | ☐ |
| 15 | 可完成城市打卡 | POST /checkin | ☐ |
| 16 | 打卡后可解锁至少一种通用成就 | 成就引擎 | ☐ |
| 17 | 游戏模式打卡后可解锁至少一种游戏专属成就 | 成就引擎 | ☐ |
| 18 | 系统记录用户通过何种模式进入城市 | city_visits.visit_mode | ☐ |
| 19 | 2C2G 单机可通过 Docker Compose 启动 | docker-compose.yml | ☐ |
| 20 | 演示过程不依赖手改数据库或临时脚本 | seed 自动导入 | ☐ |

## 验收方式
- 自动：docker compose up 后服务健康检查通过；seed 自动导入 12 城与成就。
- 手动：按上表 1→20 走通完整主链路（含一次游戏模式掷骰打卡解锁游戏专属成就）。

## 主链路定义（必须走通）
```text
首页 → 模式选择 → (自由探索点城市 / 游戏互动掷骰到城) → 城市详情 → AI 对话 → 赛博打卡 → 成就解锁
```