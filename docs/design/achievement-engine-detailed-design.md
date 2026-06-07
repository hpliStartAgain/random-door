# 成就引擎详细设计（achievement-engine-detailed-design.md）

> 对应概要设计 18 章。这是 internal/achievement 包的实现依据。约束见 go-backend-rules.md。
> 包含 2 个文件：rules.go（规则解析与单条判定）/ evaluator.go（整体评估）。
> 被 achievement_service.go 调用，在 POST /api/checkin 完成打卡后触发。

## 0. 设计目标
成就**配置化**：成就定义全部存 achievements 表(由 seed 导入)，引擎按 rule_type + rule_value 判定，新增成就**只改 seed 不改代码**(除非引入新 rule_type)。

## 1. 成就分类（对应概要 18.1）
| 类别 | 说明 | 解锁条件特征 |
|---|---|---|
| 通用成就 | 自由探索 + 游戏模式都可解锁 | 基于打卡/城市标签 |
| 游戏模式专属 | 仅游戏模式可解锁 | 基于 dice_rolls / game 访问 |

## 2. rule_type 规则字典（与 database-detailed-design.md 对齐）
| rule_type | rule_value 示例 | 判定数据源 | 含义 |
|---|---|---|---|
| first_checkin | "" | checkins | 完成任意首次打卡 |
| checkin_count | "5" | checkins | 累计打卡 ≥5 次 |
| visit_count | "5" | city_visits | 累计到过 ≥5 座城市(城市去重) |
| city_tag | "古都" | city_visits+city_tags | 到过任意一座带该标签城市 |
| tag_count | "江南:3" | city_visits+city_tags | 到过 ≥3 座带该标签城市(城市去重) |
| game_visit_count | "5" | city_visits(mode=game) | 游戏模式到达 ≥5 城 |
| dice_direction | "北:2" | dice_rolls | 连续 ≥2 次掷出该方向 |
| dice_distance | "1200" | dice_rolls | 单次掷出 ≥该距离 |

## 3. 示例成就映射（来自 achievement-design.md，供 seed）
通用：
| code | name | rule_type | rule_value |
|---|---|---|---|
| first_checkin | 初次打卡 | first_checkin | "" |
| ancient_capital_first | 古都初见 | city_tag | "古都" |
| ancient_capital_tour | 古都巡礼 | tag_count | "古都:3" |
| jiangnan_rain | 烟雨江南 | tag_count | "江南:3" |
| dongbei_first | 东北初见 | city_tag | "东北" |
| foodie_traveler | 美食旅人 | tag_count | "美食:3" |

游戏专属：
| code | name | rule_type | rule_value |
|---|---|---|---|
| game_first | 大富翁初体验 | game_visit_count | "1" |
| destiny_traveler | 命运旅人 | game_visit_count | "3" |
| roamer | 随机漫游家 | game_visit_count | "5" |
| go_north | 一路向北 | dice_direction | "北:2" |
| chosen_voyager | 天选远行者 | dice_distance | "1200" |

## 4. rules.go — 单条规则判定

### 4.1 导出
```text
// 用户行为快照，由 evaluator 一次性查好传入，避免每条规则重复查库
type UserStats struct {
  CheckinCount    int
  VisitedCityCount int              // 到过城市去重数
  VisitedCityTags map[string]int    // 标签 → 去重城市数(基于到过的城市)
  VisitedTagAny   map[string]bool   // 是否到过某标签城市
  GameVisitCount  int              // city_visits mode=game 去重城市数
  MaxDiceDistance int              // 历史最大单次距离
  MaxSameDirRun   map[string]int   // 各方向"最大连续次数"
}
func Match(ruleType, ruleValue string, s UserStats) bool
```

### 4.2 各 rule_type 判定逻辑
```text
first_checkin    : s.CheckinCount >= 1
checkin_count    : s.CheckinCount >= atoi(ruleValue)
visit_count      : s.VisitedCityCount >= atoi(ruleValue)
city_tag         : s.VisitedTagAny[ruleValue] == true
tag_count        : 解析"tag:n" → s.VisitedCityTags[tag] >= n
game_visit_count : s.GameVisitCount >= atoi(ruleValue)
dice_distance    : s.MaxDiceDistance >= atoi(ruleValue)
dice_direction   : 解析"方向:n" → s.MaxSameDirRun[方向] >= n
```
- rule_value 解析失败 → 该规则判 false(并记日志)，不 panic。

## 5. evaluator.go — 整体评估

### 5.1 导出
```text
func Evaluate(userID int64, repos Repos) (newlyUnlocked []Achievement, err error)
```

### 5.2 主流程（对应概要 18.4）
```text
1. 一次性构建 UserStats：
   - checkin_repo: CheckinCount
   - visit_repo:   到过城市去重数、按城市标签聚合(连 city_tags)、game 模式去重城市数
   - dice_repo:    历史最大距离、各方向最大连续次数(按时间排序扫描)
2. achievement_repo.ListAll() 取全部成就定义
3. achievement_repo.ListUserAchievements(userID) 取已解锁 code 集合
4. 遍历成就定义：
     已解锁 → 跳过
     rules.Match(rule_type, rule_value, stats)==true → 收集为新解锁
5. 事务内 CreateUserAchievement 逐条写入(UNIQUE 去重防并发重复)
6. 返回 newlyUnlocked 列表
```

### 5.3 "连续方向"算法（dice_direction 最复杂）
```text
取该用户 dice_rolls 按 created_at 升序
扫描，维护 currentDir / currentRun / maxRun[dir]
  本次 direction == 上一次 → currentRun++
  否则 currentRun = 1
  maxRun[direction] = max(maxRun[direction], currentRun)
```
- "一路向北 北:2" = maxRun["北"] >= 2。

## 6. 触发时机与幂等
- 触发：POST /api/visits/free、POST /api/game/roll 写入 city_visits 后自动调用 Evaluate；POST /api/checkin 写入 checkins 后仍会调用 Evaluate 以支持打卡类成就。
- 幂等：user_achievements UNIQUE(user_id, achievement_id) 保证不重复；Evaluate 已先过滤已解锁。
- 返回前端的 unlocked_achievements 只含"本次新解锁"。

## 7. 可测试性
- Match 为纯函数：构造 UserStats 覆盖每个 rule_type 边界(刚好达标/差一)。
- 连续方向算法：用方向序列(如 北,北,东,北)断言 maxRun。
- Evaluate：mock repos，验证已解锁不重复返回、新达标正确写入。
