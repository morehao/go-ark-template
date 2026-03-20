# Go应用程序Makefile

# 构建相关变量
APP =
BINARY = $(APP)
MAIN_DIR = ./apps/$(APP)/cmd
BUILD_DIR = ./output/build
VERSION = $(shell date +%Y%m%d%H%M%S)-$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

APP_CONFIG_PATH = /app/config.yaml

# go命令的环境变量
GO_ENV = CGO_ENABLED=0 GOPROXY=https://goproxy.cn,direct

# Docker 相关变量
# 获取 git tag（如果存在）
GIT_TAG = $(shell git describe --tags --exact-match 2>/dev/null)
# 获取 commit hash 短格式
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
# 构建镜像 tag：如果有 git tag 则使用 tag-commit，否则使用 commit
DOCKER_TAG = $(if $(GIT_TAG),$(GIT_TAG)-$(GIT_COMMIT),$(GIT_COMMIT))
# 完整的镜像名称：appname:tag
DOCKER_IMAGE = $(APP):$(DOCKER_TAG)
# 镜像名称（不含tag，用于查询）
DOCKER_IMAGE_NAME = $(APP)

# 伪目标
.PHONY: all build clean run lint test swag docker-build docker-run help list-apps deps tidy

# 通用入口：清理、依赖、构建并运行
all: clean deps build run

# 定义函数：验证 APP 参数是否有效
define validate_app
	@if [ -z "$(APP)" ]; then \
		echo "❌ 请使用 APP=<名称> 指定要操作的应用程序，例如：make build APP=demoapp"; \
		exit 1; \
	fi
	@if [ ! -d "./apps/$(APP)" ]; then \
		echo "❌ 应用程序 '$(APP)' 不存在于 ./apps 目录下"; \
		exit 1; \
	fi
endef

# 构建应用程序
build:
	$(call validate_app)
	@echo "正在构建应用程序 $(APP)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-X 'main.BuildVersion=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY) $(MAIN_DIR)
	@echo "✅ 构建完成: $(BUILD_DIR)/$(BINARY)"

# 为特定环境构建（例如 Linux）
build-env:
	$(call validate_app)
	@echo "正在为 $(GO_ENV) 构建 $(APP)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO_ENV) go build -ldflags="-X 'main.BuildVersion=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY) $(MAIN_DIR)
	@echo "✅ 构建完成: $(BUILD_DIR)/$(BINARY)"

# 清理构建产物
clean:
	@echo "🧹 正在清理构建目录..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ 已清理构建目录"

# 运行应用程序
run:
	$(call validate_app)
	@echo "🚀 正在运行应用程序 $(APP)..."
	@go run $(MAIN_DIR)

# 运行测试
test:
	$(call validate_app)
	@echo "🧪 正在运行测试..."
	@go test ./apps/$(APP)/internal/... -v

# 下载依赖项
deps:
	@echo "📦 正在下载依赖项..."
	@$(GO_ENV) go mod download
	@$(GO_ENV) go mod tidy
	@echo "✅ 依赖项已更新"

# 生成 Swagger 文档
swag:
	$(call validate_app)
	@echo "📚 正在生成 Swagger 文档..."
	@which swag > /dev/null || (echo "⚠️  swag 未安装，正在安装..." && go install github.com/swaggo/swag/cmd/swag@latest)
	@swag init \
		--parseDependency \
		--parseInternal \
		-g app.go \
		--dir apps/${APP} \
		--output apps/$(APP)/docs \
		--outputTypes go \
		--instanceName $(APP)
	@echo "✅ Swagger 文档已生成：apps/$(APP)/docs"


CLI_VERSION := v1.0.20
CLI_PKG     := github.com/morehao/gocli

codegen:
	$(call validate_app)
	$(if $(MODE),, $(error ❌ 请使用 MODE 参数指定生成模式，例如 MODE=api,module,model))

	@set -e; \
	if ! command -v gocli >/dev/null 2>&1; then \
		echo "⚠️ 未检测到 gocli，正在安装 $(CLI_VERSION)..."; \
		go install $(CLI_PKG)@$(CLI_VERSION); \
	else \
		INSTALLED_VER=$$(go version -m $$(which gocli) 2>/dev/null | grep -E "^\s+mod\s+$(CLI_PKG)" | awk '{print $$3}' || echo ""); \
		echo "🔍 已安装的 gocli 版本: $$INSTALLED_VER"; \
		echo "🎯 目标版本: $(CLI_VERSION)"; \
		if [ "$$INSTALLED_VER" != "$(CLI_VERSION)" ]; then \
			echo "⚠️ gocli 版本不匹配，重新安装 $(CLI_VERSION)..."; \
			go install $(CLI_PKG)@$(CLI_VERSION); \
		else \
			echo "✅ gocli 版本已是最新"; \
		fi; \
	fi

	@echo "🔧 开始生成代码：APP=$(APP)，MODE=$(MODE)"
	@gocli generate --mode=$(MODE) --app=$(APP)


# 构建 Docker 镜像
docker-build:
	$(call validate_app)
	@echo "🐳 正在构建 $(APP) 的 Docker 镜像..."
	@echo "   镜像名称: $(DOCKER_IMAGE)"
	@echo "   Git Tag: $(if $(GIT_TAG),$(GIT_TAG),无)"
	@echo "   Git Commit: $(GIT_COMMIT)"
	docker buildx build -t $(DOCKER_IMAGE) -f ./apps/$(APP)/scripts/Dockerfile .
	@echo "✅ Docker 镜像 $(DOCKER_IMAGE) 已构建完成"

# 运行 Docker 容器
docker-run: check-image
	@echo "🚀 正在运行 $(APP) 容器..."
	@echo "   使用镜像: $(DOCKER_IMAGE)"
	-@docker rm -f $(APP) 2>/dev/null || true
	@docker run -d \
		--name $(APP) \
		-e APP_CONFIG_PATH=$(APP_CONFIG_PATH) \
		-p 8099:8099 \
		$(DOCKER_IMAGE)
	@echo "✅ 容器 $(APP) 已启动，服务地址：http://localhost:8099"

# 检查镜像是否存在，没有就构建
check-image:
	@if [ -n "$$(docker images -q $(DOCKER_IMAGE))" ]; then \
		echo "⚠️ 镜像 $(DOCKER_IMAGE) 已存在，准备删除重建..."; \
		docker rmi -f $(DOCKER_IMAGE); \
	fi
	$(MAKE) docker-build

# 列出所有可用的应用程序
list-apps:
	@echo "📂 可用的应用程序:"
	@ls -1 ./apps

# 运行代码检查工具
lint:
	@echo "🔍 正在运行代码检查工具..."
	@golangci-lint run ./...

# 显示帮助信息
help:
	@echo "🆘 可用命令:"
	@echo "  make                    - 清理、下载依赖并构建应用程序"
	@echo "  make build APP=<名称>    - 构建指定的应用程序"
	@echo "  make build-env APP=<名称> - 为特定环境构建"
	@echo "  make clean              - 清理构建产物"
	@echo "  make deps               - 下载依赖项"
	@echo "  make run APP=<名称>     - 运行指定的应用程序"
	@echo "  make test APP=<名称>    - 运行测试"
	@echo "  make swag APP=<名称>    - 生成 Swagger 文档"
	@echo "  make codegen APP=<名称> MODE=<模式> - 生成代码（MODE=api,module,model）"
	@echo "  make docker-build APP=<名称>  - 构建 Docker 镜像（标签基于 git tag 和 commit hash）"
	@echo "  make docker-run APP=<名称> - 运行 Docker 容器"
	@echo "  make list-apps          - 列出所有可用的应用程序"
	@echo "  make lint               - 运行代码检查工具"
	@echo ""
	@echo "📝 Docker 镜像标签说明："
	@echo "  - 如果当前 commit 有 git tag：镜像标签为 <tag>-<commit>"
	@echo "  - 如果当前 commit 无 git tag：镜像标签为 <commit>"
	@echo "  - 示例：demoapp:v1.0.0-abc1234 或 demoapp:abc1234"
