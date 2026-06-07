# 用户流程文档

## 1. 首次进入

```text
打开网页
  -> useUserStore 从本地持久化读取 anonymousId
  -> 无 anonymousId 则生成 UUID
  -> POST /api/users/anonymous 创建或恢复用户
  -> 展示 WelcomeOverlay
  -> 用户选择自由探索 / 任意门 / 注册登录
```

相关接口：`POST /api/users/anonymous`、`POST /api/auth/register`、`POST /api/auth/login`。

## 2. 自由探索

```text
进入主工作台
  -> GET /api/cities 加载城市列表
  -> Sidebar 搜索 / 区域筛选 / 标签筛选
  -> 地图 marker 与列表同步过滤
  -> 点击城市
  -> POST /api/visits/free 写访问记录
  -> GET /api/cities/{city_id} 加载详情
  -> 展示 CityDetailPanel
```

相关接口：`GET /api/cities`、`POST /api/visits/free`、`GET /api/cities/{city_id}`。

## 3. 任意门随机漫游

```text
点击任意门入口
  -> POST /api/game/init 获取最近起点城市
  -> 用户开启随机
  -> POST /api/game/roll
       随机方向 + 随机距离
       计算目标经纬度
       匹配目标城市
       写 dice_rolls + city_visits
       触发成就评估
  -> 前端展示揭晓与地图飞行动画
  -> 进入目标城市详情
```

相关接口：`POST /api/game/init`、`POST /api/game/roll`。

## 4. 城市探索

```text
城市详情
  -> 查看简介、标签、地标、美食、人物、方言
  -> 点击地标或城市风光
  -> StreetViewCanvas 展示图片沉浸视图与声景控制
  -> 可进入猜一猜或赛博打卡
```

相关接口：`GET /api/comments`、`POST /api/comments`。

## 5. AI 人物对话

```text
选择城市人物
  -> RightDrawer 打开对话
  -> 用户发送问题
  -> POST /api/chat
       后端读取城市、人物、历史消息
       组装合规 prompt
       调外部 LLM
       保存 user / assistant 消息
  -> 前端展示回复
```

相关接口：`POST /api/chat`。

## 6. 赛博打卡

```text
选择地标或风光视图中的打卡入口
  -> 上传自拍
  -> 可附带当前风光截图作为 scene_file
  -> POST /api/checkin/generate-image 创建 ai_tasks
  -> 前端轮询 GET /api/checkin/image-tasks/{task_id}
  -> 成功后展示生成图，可失败重试
  -> POST /api/checkin 确认打卡
  -> 生成海报并触发成就解锁
```

相关接口：`POST /api/checkin/generate-image`、`GET /api/checkin/image-tasks/{task_id}`、`POST /api/checkin/image-tasks/{task_id}/retry`、`POST /api/checkin`。

## 7. 猜城市挑战

```text
风光浏览中点击截图生成文案
  -> POST /api/guess/caption
  -> POST /api/guess/challenges 创建分享 code
  -> 好友打开 /?guess=code
  -> GET /api/guess/challenges/{code}
  -> POST /api/guess/challenges/{code}/answers
  -> 返回是否猜中
```

## 8. 资产与成就

```text
点击我的足迹 / 我的资产
  -> GET /api/users/{user_id}/profile
  -> GET /api/users/{user_id}/assets
  -> 可编辑昵称、年龄、地区

点击成就墙
  -> GET /api/users/{user_id}/achievements
  -> 展示已解锁、未解锁与进度
```

相关接口：`GET/PATCH /api/users/{user_id}/profile`、`GET /api/users/{user_id}/assets`、`GET /api/users/{user_id}/achievements`。

## 9. 后台内容维护

```text
点击后台入口
  -> 输入 ADMIN_TOKEN
  -> GET /api/admin/catalog/coverage 验证 token 并获取覆盖率
  -> 维护城市、标签、地标、美食、人物、成就和媒体
```

相关接口：`/api/admin/*`，详见 `../design/api-contract.md`。
