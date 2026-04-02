MODULE_NAME := $(shell cat go.mod 2>/dev/null | head -1 | awk '{print $$2}' || echo "test")
PROJECT_NAME := $(notdir $(MODULE_NAME))

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# 构建标志
BUILDTAGS := jsoniter
LDFLAGS := -X '$(MODULE_NAME)/internal/version.Version=$(VERSION)' \
           -X '$(MODULE_NAME)/internal/version.GitCommit=$(GIT_COMMIT)' \
           -X '$(MODULE_NAME)/internal/version.BuildTime=$(BUILD_TIME)'

APPS=$(shell ls -1 ./cmd/ | grep -v -E "(internal|inject|.*\.md|.*\.go)")
APPS_COUNT=$(shell ls -1 ./cmd/ | grep -v -E "(internal|inject|.*\.md|.*\.go)" | wc -l)

.PHONY: all build test clean install lint fmt help vet check typos table svc config wire enum swag

# 默认目标
all: test build

# 构建项目
build:
	@echo "Building $(PROJECT_NAME) Project"
	@mkdir -p ./bin
	@for app in $$(ls -1 ./cmd/ | grep -v -E "(internal|inject|.*\.md|.*\.go)"); do \
		echo "Building $$app..."; \
		go build -tags $(BUILDTAGS) -ldflags "$(LDFLAGS)" -o ./bin/$$app $(MODULE_NAME)/cmd/$$app; \
	done
	@echo "Build complete: ./bin/"

# 运行测试
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 清理构建产物
clean:
	@echo "Cleaning..."
	rm -rf $(PROJECT_NAME)/bin
	rm -f coverage.out coverage.html
	go clean
	@echo "Clean complete"

# 安装到 GOPATH/bin
install:
	@echo "Installing $(PROJECT_NAME)..."
	go install -ldflags "$(LDFLAGS)" ./cmd/...
	@echo "Install complete"

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run ./...

typos:
	@echo "Check typos code..."
	typos

# 格式化代码
fmt:
	@echo "Formatting code..."
	goimports -l -d -w .
	@echo "Format complete"

# 静态分析代码
vet:
	@echo "Static Analysis code ..."
	go vet ./...
	@echo "Static Analysis complete"

# 漏洞检测分析
check:
	@echo "Vulnerability Detection code ..."
	govulncheck ./...
	@echo "Vulnerability Detection complete"

# 运行项目
run:
	@if [ ${APPS_COUNT} -gt 1 ]; \
	then sh -c 'read -p "please input run app name [ ${APPS} ] : " RUN_PROJECT;go run -tags $(BUILDTAGS) -ldflags "$(LDFLAGS)" ./cmd/$${RUN_PROJECT} $(ARGS);'; else \
	go run -tags $(BUILDTAGS) -ldflags "$(LDFLAGS)" ./cmd/$(APPS) $(ARGS); \
	fi;

# 生成数据库结构
table:
	@echo "Generating database structure..."
	gogenie db2struct

# 生成服务接口
svc:
	@echo "Generating service interface..."
	gogenie http router && gogenie http client && gogenie impl && gogenie mount && gogenie http swagger && gogenie autowire

# 生成配置注入
config:
	@echo "Generating config injection..."
	gogenie mount config

# 生成自动注入
wire:
	@echo "Generating auto-injection..."
	gogenie autowire

# 枚举生成
enum:
	@echo "Generating enum..."
	gogenie enum

# 生成 Swagger 文档（自动扫描含 @Router 的目录，无需随接口增多改命令）
swag:
	@echo "Generating swagger docs..."
	@SWAG_DIRS="./cmd/server,./service"; \
	for d in $$(grep -rl '@Router' --include='*.go' . 2>/dev/null | grep -v -E '(\.git|/vendor/)' | xargs -I {} dirname {} | sort -u); do \
		[ -n "$$d" ] && SWAG_DIRS="$$SWAG_DIRS,$$d"; \
	done; \
	swag init -d "$$SWAG_DIRS" --parseInternal=true
	@echo "Swagger docs generated in ./docs"

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  all            - Run tests and build (default)"
	@echo "  build          - Build the project"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install to GOPATH/bin"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Static analysis code"
	@echo "  check          - Vulnerability detection code"
	@echo "  typos          - Check typos in code"
	@echo "  run            - Run the project (use ARGS='...' for arguments)"
	@echo "  table          - Generate database structure"
	@echo "  svc            - Generate service interface"
	@echo "  config         - Generate config injection"
	@echo "  wire       - Generate auto-injection"
	@echo "  enum           - Generate enum"
	@echo "  swag           - Generate swagger docs (auto-scans @Router dirs)"
	@echo "  help           - Show this help message"
	@echo "Examples:"
	@echo "  make build"
	@echo "  make test"
	@echo "  make run ARGS='db2struct users'"
	@echo "  make install"
