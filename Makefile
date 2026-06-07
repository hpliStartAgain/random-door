# Makefile — 常用命令封装，降低团队上手成本
# 对应 README.md。用法：make up / make seed ...

.PHONY: up down build logs migrate seed seed-audit seed-sync seed-landmark-coords lint test fe-dev fe-build

# ---------- Docker ----------
up:            ## 启动全部容器
	docker-compose up -d

down:          ## 停止并移除容器
	docker-compose down

build:         ## 构建镜像
	docker-compose build

logs:          ## 查看后端日志
	docker-compose logs -f app

# ---------- 数据 ----------
# docker compose 使用外部 MySQL；应用启动只自动建表，不默认写入 seed 内容。
migrate:       ## 已废弃，应用启动自动执行
	@echo "Migration is auto-executed on backend startup."

seed-audit:    ## 只读盘点真实库与 seed 的差异
	cd backend && go run ./cmd/seedtool -mode audit

seed:          ## 显式补齐缺失 seed 行，不覆盖后台已有内容
	cd backend && go run ./cmd/seedtool -mode bootstrap

seed-sync:     ## 谨慎：按 seed 覆盖自然键匹配行
	cd backend && go run ./cmd/seedtool -mode sync -confirm-overwrite

seed-landmark-coords: ## 仅回填缺失的地标经纬度，不覆盖图片/描述等内容
	cd backend && go run ./cmd/seedtool -mode backfill-landmark-coordinates

# ---------- 质量 ----------
lint:          ## 后端 + 前端静态检查
	cd backend && gofmt -l . && go vet ./...
	cd frontend && npm run lint

test:          ## 后端单元测试（geo / achievement 为重点）
	cd backend && go test ./...

# ---------- 前端本地 ----------
fe-dev:        ## 前端开发服务器
	cd frontend && npm run dev

fe-build:      ## 前端构建产物
	cd frontend && npm run build
