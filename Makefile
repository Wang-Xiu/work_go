.PHONY: help test test-coverage run-example clean deps

help: ## 显示帮助信息
	@echo "可用命令："
	@echo "  make deps          - 安装项目依赖"
	@echo "  make test          - 运行所有单元测试"
	@echo "  make test-coverage - 运行测试并生成覆盖率报告"
	@echo "  make run-example   - 运行示例程序"
	@echo "  make clean         - 清理生成的文件"

deps: ## 安装项目依赖
	@echo "正在安装依赖..."
	go mod download
	go mod tidy

test: ## 运行所有单元测试
	@echo "运行单元测试..."
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率报告..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

run-example: ## 运行示例程序（需要先启动Redis）
	@echo "运行示例程序..."
	@echo "请确保Redis已启动（localhost:6379）"
	go run example/main.go

clean: ## 清理生成的文件
	@echo "清理临时文件..."
	rm -f coverage.out coverage.html
	rm -f *.test
	go clean

.DEFAULT_GOAL := help
