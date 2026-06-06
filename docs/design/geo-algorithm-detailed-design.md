# 随机漫游算法详细设计（geo-algorithm-detailed-design.md）

> 对应概要设计 17 章。这是 internal/geo 包的实现依据。约束见 go-backend-rules.md。
> 包含 4 个文件：distance.go / bearing.go / target_point.go / city_matcher.go。
> 被 game_service.go 调用，服务于 POST /api/game/init 与 POST /api/game/roll。

## 0. 算法总览
游戏互动模式（大富翁式漫游）的核心：从当前位置出发，**随机方向 + 随机距离** → 算出一个目标经纬度 → 在 seed 城市库中匹配**距离目标点最近、且非当前城、优先未访问**的城市。

```text
当前位置(lat,lng)
  ↓ RandomDirection()      → 8 方向之一 + 方位角
  ↓ RandomDistance()       → 6 档距离之一
  ↓ TargetPoint()          → 目标经纬度(球面投影)
  ↓ MatchNearestCity()     → 遍历城市库算距离 + 过滤 + 兜底
  ↓
目标城市
```

## 1. 输入参数
| 参数 | 类型 | 来源 | 说明 |
|---|---|---|---|
| lat | float64 | 请求体/init 起点 | 当前纬度 |
| lng | float64 | 请求体/init 起点 | 当前经度 |
| from_city_id | int64 | 请求体 | 当前城市，用于过滤 |
| direction | string | RandomDirection 产出 | 随机方向 |
| distance_km | int | RandomDistance 产出 | 随机距离 |

## 2. bearing.go — 方向

### 2.1 八方向角度映射表（正北为 0°，顺时针）
| 方向 | 方位角 deg |
|---|---|
| 北 | 0 |
| 东北 | 45 |
| 东 | 90 |
| 东南 | 135 |
| 南 | 180 |
| 西南 | 225 |
| 西 | 270 |
| 西北 | 315 |

### 2.2 导出（伪签名）
```text
type Direction struct { Name string; Bearing float64 }
var Directions [8]Direction   // 上表
func RandomDirection() Direction   // 均匀随机取一个
func RandomDirectionWithRand(rng IntnSource) Direction   // 测试注入固定随机源
```
- 实现：rand.Intn(8) 取下标。
- 测试点：返回 Name 必属上表 8 值；Bearing 与 Name 对应。

## 3. target_point.go — 距离与目标点

### 3.1 距离档位
```text
var DistanceLevels = []int{100, 200, 300, 500, 800, 1200}   // km
func RandomDistance() int   // rand 取一档
func RandomDistanceWithRand(rng IntnSource) int   // 测试注入固定随机源
```

### 3.2 目标点计算（球面正算 / Destination point，地球半径 R=6371km）
已知起点(lat1,lng1)、方位角 θ、距离 d，求目标点(lat2,lng2)：
```text
δ = d / R                       // 角距离(弧度)
θ = bearing 转弧度
φ1 = lat1 转弧度, λ1 = lng1 转弧度

φ2 = asin( sin φ1 · cos δ + cos φ1 · sin δ · cos θ )
λ2 = λ1 + atan2( sin θ · sin δ · cos φ1 , cos δ − sin φ1 · sin φ2 )

lat2 = φ2 转角度
lng2 = ((λ2 转角度) + 540) mod 360 − 180   // 归一化到 [-180,180)
```

### 3.3 导出
```text
func TargetPoint(lat, lng, bearingDeg float64, distKm int) (lat2, lng2 float64)
```
- 边界：跨经线 ±180 已由公式归一化；任意圈数的输入经度也会稳定归一化；中国境内不跨极点，无需特判。

## 4. distance.go — Haversine 两点距离
```text
func Haversine(lat1, lng1, lat2, lng2 float64) float64   // 返回 km
```
公式：
```text
dφ = (lat2−lat1) 弧度
dλ = (lng2−lng1) 弧度
a = sin²(dφ/2) + cos φ1 · cos φ2 · sin²(dλ/2)
c = 2 · atan2(√a, √(1−a))
d = R · c       // R=6371
```
- 用途：MatchNearestCity 算每城到目标点距离；也可算实际到达距离写日志。

## 5. city_matcher.go — 最近城市匹配（核心，含兜底）

### 5.1 导出
```text
type MatchOptions struct {
    ExcludeCityID  int64     // 当前城市，必排除
    VisitedCityIDs []int64   // 该用户已访问城市集合
}
func MatchNearestCity(cities []City, targetLat, targetLng float64, opt MatchOptions) (City, error)
```

### 5.2 主流程
```text
1. 遍历 cities，对每城算 Haversine(目标点, 城市)
2. 过滤掉 ExcludeCityID（当前城市）
3. 候选拆两组：未访问组 / 已访问组
4. 未访问组非空 → 返回该组中距目标点最近城市   ← 正常路径(优先未访问)
5. 否则进入兜底(见 5.3)
```

### 5.3 兜底策略（对应概要 17.5）
```text
全部非当前城市都已访问过 → 在"全部非当前城"中选距目标点最近者，允许重复访问
```
- 随机方向与随机距离已经编码为目标点；兜底仍按目标点选最近城市，避免再次随机或重复施加方向偏差。
- 同距离时按 city_id 较小者稳定决胜，保证输入顺序不影响结果。
- 保证：城市库 ≥2 座且坐标合法时必返回非当前城市；error 仅在 cities 为空、仅有当前城或坐标非法时返回。

### 5.4 边界条件
| 场景 | 处理 |
|---|---|
| cities 为空 | 返回 error（理论不会，seed 至少有 12 城） |
| 仅当前城 1 座 | 返回 error 或允许停留(service 决定提示) |
| 目标点落海/境外 | 不影响——只看哪座城离目标点最近 |
| 已访问全部城市 | 在全部非当前城中选距目标点最近者，允许重复 |
| 目标点或城市坐标非法 | 返回 error，交由 service 转成友好错误 |

## 6. game_service 调用示例（逻辑，非代码）
```text
dir  = geo.RandomDirection()
dist = geo.RandomDistance()
tLat,tLng = geo.TargetPoint(lat, lng, dir.Bearing, dist)
city,_ = geo.MatchNearestCity(allCities, tLat, tLng, MatchOptions{
    ExcludeCityID: fromCityID, VisitedCityIDs: visitedIDs })
// 写 dice_rolls(direction=dir.Name, distance_km=dist, target_lat/lng, to_city_id=city.ID)
// 写 city_visits(visit_mode=game, source=dice_roll, from_city_id, dice_roll_id)
```

## 7. 可测试性
- RandomDirection/RandomDistance 可注入 rand 源，便于固定种子单测。
- TargetPoint/Haversine 为纯函数，用已知经纬度对(如北京→西安约 910km)断言。
- MatchNearestCity 覆盖：正常/排除当前城/优先未访问/全访问允许重复/稳定决胜/空城市/仅当前城/非法坐标。
