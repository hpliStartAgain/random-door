// backend/go.mod — Go 依赖清单（版本可按需微调，提交时锁定 go.sum）
module github.com/your-org/city-roam/backend

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/google/uuid v1.6.0
	github.com/spf13/viper v1.18.2
	gorm.io/driver/mysql v1.5.6
	gorm.io/gorm v1.25.10
)

// 说明：
// - gin：HTTP 框架
// - gorm + driver/mysql：ORM 与 MySQL 驱动
// - viper：配置/环境变量加载
// - uuid：上传文件名/匿名 id 辅助
// - 日志使用标准库 log/slog（Go 1.21+ 内置，无需额外依赖）
// 禁止引入：redis / kafka / es / 向量库 等中间件客户端