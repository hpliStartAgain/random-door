# 前端详细设计（frontend-detailed-design.md）

> 对应概要设计 13 章。这是 agent 创建所有前端 `.ts/.tsx` 文件的直接依据。约束见 react-frontend-rules.md。
> 状态用 Zustand；请求统一走 api 层；地图只在 MapCanvas 封装；字段类型对齐 api-contract.md。

## 0. 完整前端文件树
```
frontend/src/
  pages/      HomePage ModeSelectPage FreeExplorePage GameModePage
              CityPage ChatPage CheckinPage AchievementPage
  components/ MapCanvas CityMarker DicePanel CityCard LandmarkCard
              FoodCard CharacterCard DialectCard ChatBox ImageUploader
              CheckinResult BadgeWall ModeEntryCard
  api/        client.ts city.ts game.ts visit.ts chat.ts checkin.ts achievement.ts
  store/      useUserStore.ts useGameStore.ts useCityStore.ts
  App.tsx main.tsx router.tsx
```

---

## 1. pages（8 个路由页）
| 页面 | 路由 | 职责 | 数据来源 | 调用接口 |
|---|---|---|---|---|
| HomePage | / | 产品介绍 + 进入按钮；首次进入触发匿名用户创建 | useUserStore | POST /users/anonymous |
| ModeSelectPage | /mode | 展示自由探索/游戏互动两入口卡片 | - | - |
| FreeExplorePage | /explore | 地图打点，点击城市→写 free visit→跳详情 | useCityStore | GET /cities、POST /visits/free |
| GameModePage | /game | 获取定位→init→掷骰→移动动画→跳详情 | useGameStore | POST /game/init、POST /game/roll |
| CityPage | /city/:id | 城市详情（两模式共用）：简介/地标/美食/人物/方言 + 对话/打卡/成就入口 | useCityStore | GET /cities/{id} |
| ChatPage | /city/:id/chat/:cid | 与人物对话 | - | POST /chat |
| CheckinPage | /city/:id/checkin | 上传照片→生图→确认打卡 | - | POST /checkin/generate-image、POST /checkin |
| AchievementPage | /achievements | 成就墙（已解锁/未解锁/进度） | useUserStore | GET /users/{id}/achievements |

---

## 2. components（13 个）
| 组件 | props（关键） | 职责/交互 |
|---|---|---|
| MapCanvas | cities, onCityClick, mode, movePath? | **唯一封装高德 SDK**；打 Marker；游戏模式播放移动动画 |
| CityMarker | city, onClick | 地图标记（也可由 MapCanvas 内部生成） |
| DicePanel | onRoll, rolling, result | 掷骰按钮 + 动画 + 显示方向/距离 |
| CityCard | city | 城市简介卡 |
| LandmarkCard | landmark, onCheckin? | 地标卡，可触发打卡 |
| FoodCard | food | 美食卡 |
| CharacterCard | character, onChat | 人物卡，点击进对话 |
| DialectCard | sample, explanation | 方言样例 + 解释 |
| ChatBox | messages, onSend, loading | 对话气泡 + 输入框 |
| ImageUploader | onSelect, maxSizeMB=5, accept | 选图 + 前端预校验（类型/大小） |
| CheckinResult | imageUrl, onConfirm, onRetry | 展示生成图 + 确认/重试 |
| BadgeWall | unlocked, locked, progress | 成就墙网格 |
| ModeEntryCard | mode, title, desc, onEnter | 模式入口卡 |

---

## 3. api（接口封装，组件不得直接 fetch）
### 3.1 client.ts
- 创建 axios 实例：baseURL=/api；请求拦截器注入 X-User-Id（从 useUserStore）；响应拦截器统一解析错误 `{error:{code,message}}` 并抛出。
### 3.2 各业务文件（函数 + TS 类型，类型对齐 api-contract.md）
| 文件 | 导出函数 |
|---|---|
| city.ts | listCities()、getCity(id) |
| game.ts | initGame(lat,lng)、rollDice(fromCityId,lat,lng) |
| visit.ts | createAnonymousUser(anonymousId)、createFreeVisit(cityId,source) |
| chat.ts | sendChat(cityId,characterId,message) |
| checkin.ts | generateImage(form)、createCheckin(payload) |
| achievement.ts | getAchievements(userId) |

---

## 4. store（Zustand，3 个）
| store | 状态 | 动作 |
|---|---|---|
| useUserStore | userId, anonymousId, currentCityId | initUser()（读 localStorage，无则建匿名用户并存） |
| useGameStore | fromCity, lastRoll, targetCity, rolling | setInit、roll、reset |
| useCityStore | cities[], cityCache{} | loadCities()、loadCity(id)（带缓存） |

---

## 5. 入口与路由
- main.tsx：挂载 App，初始化高德 SDK loader（仅 VITE_AMAP_KEY，前端公开 Key）。
- App.tsx：布局 + 路由出口；首屏调用 useUserStore.initUser()。
- router.tsx：定义上表 8 条路由。

## 6. 安全提醒
- 前端**不得出现 LLM/生图 Key**；只允许 VITE_AMAP_KEY。
- 所有请求经 api 层；高德调用只在 MapCanvas。