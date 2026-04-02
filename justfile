set shell := ["bash", "-c"]

# 项目变量
APP_NAME := "{{.ProjectName}}"
APP := ""
APPS := `find ./cmd -mindepth 1 -maxdepth 1 -type d -not -name internal -exec basename {} \; | tr '\n' ' '`
MODULE := `go list -m`
VERSION := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
BUILD_TIME := `date -u +"%Y-%m-%d_%H:%M:%S"`
GIT_COMMIT := `git rev-parse HEAD 2>/dev/null || echo "unknown"`

# Go构建标志
LDFLAGS := "-ldflags \"-X '" + MODULE + "/internal/version.Version=" + VERSION + "' -X '" + MODULE + "/internal/version.GitCommit=" + GIT_COMMIT + "' -X '" + MODULE + "/internal/version.BuildTime=" + BUILD_TIME + "'\""

# 默认任务
default: build

# 开发任务
dev: fmt lint test build

# 构建应用
build:
    @echo "构建应用..."
    @mkdir -p bin
    @for app in {{APPS}}; do \
        [ -d "./cmd/$app" ] || continue; \
        echo "  -> $app"; \
        go build {{LDFLAGS}} -o bin/$app ./cmd/$app; \
    done
    @echo "构建完成: ./bin"

# 构建所有平台
all: clean
    @echo "构建所有平台版本..."
    @mkdir -p bin
    @for app in {{APPS}}; do \
        GOOS=linux GOARCH=amd64 go build {{LDFLAGS}} -o bin/$app-linux-amd64 ./cmd/$app; \
        GOOS=linux GOARCH=arm64 go build {{LDFLAGS}} -o bin/$app-linux-arm64 ./cmd/$app; \
        GOOS=darwin GOARCH=amd64 go build {{LDFLAGS}} -o bin/$app-darwin-amd64 ./cmd/$app; \
        GOOS=darwin GOARCH=arm64 go build {{LDFLAGS}} -o bin/$app-darwin-arm64 ./cmd/$app; \
        GOOS=windows GOARCH=amd64 go build {{LDFLAGS}} -o bin/$app-windows-amd64.exe ./cmd/$app; \
    done
    @echo "所有平台构建完成"

# 运行应用
apps:
    @echo "可用服务:"
    @for app in {{APPS}}; do echo "  - $app"; done

run *ARGS:
    @set -- {{ARGS}}; \
    app=""; \
    if [ $# -ge 1 ]; then \
        for a in {{APPS}}; do \
            if [ "$1" = "$a" ]; then app="$1"; shift; break; fi; \
        done; \
    fi; \
    if [ -z "$app" ]; then \
        echo "请选择服务:"; \
        select app in {{APPS}}; do \
            if [ -n "$app" ]; then break; fi; \
        done; \
    fi; \
    echo "运行 $app..."; \
    go run {{LDFLAGS}} ./cmd/$app "$@"

# 运行开发模式
run-dev:
    @echo "启动开发模式..."
    ENV=development go run {{LDFLAGS}} ./cmd/{{APP}} --log-dir ./logs

# 代码格式化
fmt:
    @echo "格式化代码..."
    go fmt ./...
    goimports -w .

# 静态分析
lint:
    @echo "运行静态分析..."
    golangci-lint run --config .golangci.yml
    go vet ./...

# 安全扫描
security:
    @echo "运行安全扫描..."
    gosec ./...

# 运行测试
test:
    @echo "运行测试..."
    go test -v -race -coverprofile=coverage.out ./...

# 运行基准测试
bench:
    @echo "运行基准测试..."
    go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./...

# 测试覆盖率
coverage: test
    @echo "生成测试覆盖率报告..."
    go tool cover -html=coverage.out -o coverage.html
    @echo "覆盖率报告已生成: coverage.html"

# 性能分析
profile:
    @echo "启动性能分析..."
    go tool pprof -http=:8080 cpu.prof

# 生成代码
generate:
    @echo "生成代码..."
    go generate ./...

# 依赖管理
deps:
    @echo "更新依赖..."
    go mod tidy
    go mod verify

# 依赖检查
deps-check:
    @echo "检查依赖更新..."
    go list -u -m all

# 清理构建文件
clean:
    @echo "清理构建文件..."
    rm -rf bin/
    rm -f coverage.out coverage.html
    rm -f cpu.prof mem.prof
    go clean -cache
    go clean -testcache

# Docker构建
docker-build:
    @echo "构建Docker镜像..."
    docker build -t {{APP_NAME}}:{{VERSION}} -t {{APP_NAME}}:latest .

# Docker多架构构建
docker-build-multi:
    @echo "构建多架构Docker镜像..."
    docker buildx build --platform linux/amd64,linux/arm64 -t {{APP_NAME}}:{{VERSION}} -t {{APP_NAME}}:latest --push .

# Docker运行
docker-run:
    @echo "运行Docker容器..."
    docker run -d --name {{APP_NAME}} \
        -v {{APP_NAME}}-data:/root/.mcp-ssh \
        -v ~/.ssh:/root/.ssh:ro \
        {{APP_NAME}}:latest

# Docker停止
docker-stop:
    docker stop {{APP_NAME}} || true
    docker rm {{APP_NAME}} || true

# 安装到本地
install:
    @echo "安装应用..."
    go install {{LDFLAGS}} ./cmd/...
    @echo "安装完成"

# 卸载
uninstall:
    @echo "从本地卸载..."
    sudo rm -f /usr/local/bin/{{APP_NAME}}

# 创建发布包
release: clean all
    @echo "创建发布包..."
    mkdir -p dist
    cd bin && for file in {{APP_NAME}}-*; do \
        if [[ $file == *.exe ]]; then \
            zip -r ../dist/${file%.exe}.zip $file ../README.md ../LICENSE; \
        else \
            tar -czf ../dist/$file.tar.gz $file ../README.md ../LICENSE; \
        fi \
    done
    @echo "发布包已创建在 dist/ 目录"

# 开发环境初始化
init:
    @echo "初始化开发环境..."
    go mod download
    @echo "开发环境初始化完成"

# 检查代码质量
quality: fmt lint security test
    @echo "代码质量检查完成"

# CI/CD流水线
ci: deps quality all docker-build
    @echo "CI流水线执行完成"

# 静态分析代码
vet:
    @echo "静态分析代码 ..."
    go vet ./...
    @echo "静态分析代码完成"

typos:
    @echo "检查错别字代码..."
    typos
    @echo "检查错别字代码完成..."

# 挂载配置
config:
    @echo "挂载配置文件..."
    gogenie mount
    @echo "挂载配置文件完成..."

# 注入依赖
wire:
    @echo "自动注入依赖..."
    gogenie autowire
    @echo "自动注入依赖完成..."

# 构建数据模型
table:
    @echo "自动获取数据模型..."
    gogenie db2struct
    @echo "自动获取数据模型完成..."

# 构建数据模型
svc:
    @echo "自动实现服务中..."
    gogenie impl && gogenie http client && gogenie http router && gogenie http swagger && gogenie mount && gogenie autowire
    @echo "自动实现服务完成..."

impl:
    @echo "自动实现服务中..."
    gogenie impl
    @echo "自动实现服务完成..."

# 生成 Swagger 文档（自动扫描含 @Router 的目录）
swag:
    #!/usr/bin/env bash
    set -e
    SWAG_DIRS="./cmd/server,./service"
    for d in $(grep -rl '@Router' --include='*.go' . 2>/dev/null | grep -v '\.git' | grep -v '/vendor/' | xargs -I {} dirname {} | sort -u); do
        [ -n "$d" ] && SWAG_DIRS="$SWAG_DIRS,$d"
    done
    swag init -d "$SWAG_DIRS" --parseInternal --parseDependency
    echo "Swagger 文档已生成到 ./docs"

# 帮助信息
help:
    @echo "可用命令:"
    @just --list
