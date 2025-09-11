# Echo API Makefile

.PHONY: help build test test-unit test-local deploy clean deps

# デフォルト環境
ENV ?= prod

# AWS設定
AWS_REGION ?= ap-northeast-1

help: ## ヘルプを表示
	@echo "Echo API - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

deps: ## Go モジュールの依存関係をダウンロード
	go mod download
	go mod tidy

build-local: ## ローカルでGoバイナリをビルド
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main ./cmd/lambda

build: ## SAMでLambda関数をビルド
	@echo "Building Lambda function with SAM..."
	sam build --config-env $(ENV)

deploy: ## SAMでprod環境にデプロイ
	@echo "Deploying to $(ENV) environment with SAM..."
	sam deploy --config-env $(ENV)

deploy-prod: ## prod環境にデプロイ
	$(MAKE) deploy ENV=prod

build-and-deploy: ## ビルドとデプロイを連続実行
	@echo "Building and deploying..."
	$(MAKE) build ENV=$(ENV)
	$(MAKE) deploy ENV=$(ENV)

clean: ## ビルドアーティファクトをクリーンアップ
	rm -rf .aws-sam/
	rm -f bin/main
	docker image prune -f

logs: ## Lambda関数のログを表示
	sam logs --name echo-api-prod --stack-name echo-api-prod --tail

logs-api: ## API Gatewayのログを表示
	aws logs tail /aws/apigateway/echo-api-gateway-prod --follow --region $(AWS_REGION)

sam-build: ## SAM buildのみ実行
	sam build --config-env $(ENV)

sam-deploy: ## SAM deployのみ実行
	sam deploy --config-env $(ENV)

sam-validate: ## SAMテンプレートの検証
	sam validate

sam-local-api: ## SAM local start-apiを起動
	sam local start-api --port 3000

sam-local-invoke: ## SAM local invokeでテスト
	echo '{"httpMethod":"GET","path":"/test","headers":{},"queryStringParameters":{"test":"value"}}' | sam local invoke EchoFunction

sam-sync: ## SAM syncで高速デプロイ（開発用）
	sam sync --stack-name $(shell sam list stack-outputs --config-env $(ENV) --output json | jq -r '.[0].StackName') --watch

format: ## コードをフォーマット
	go fmt ./...

lint: ## コードをリント
	golangci-lint run

install-tools: ## 開発ツールをインストール
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# SAM用ビルドターゲット
build-EchoFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/lambda

# Docker関連のターゲット
docker-build: ## Dockerイメージをビルド
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o bootstrap ./cmd/lambda
	docker build -t echo-api-lambda:latest .

docker-run: ## ローカルでDockerコンテナをテスト実行
	docker run -p 9000:8080 echo-api-lambda:latest

# サンプルリクエスト
sample-get: ## デプロイされたAPIにGETリクエストを送信
	@if [ -z "$(API_URL)" ]; then \
		echo "Error: API_URL environment variable is required"; \
		echo "Example: make sample-get API_URL=https://your-api-id.execute-api.us-east-1.amazonaws.com/dev/"; \
		exit 1; \
	fi
	curl -X GET "$(API_URL)test?param1=value1&param2=value2" -H "Content-Type: application/json" | jq .

sample-post: ## デプロイされたAPIにPOSTリクエストを送信
	@if [ -z "$(API_URL)" ]; then \
		echo "Error: API_URL environment variable is required"; \
		echo "Example: make sample-post API_URL=https://your-api-id.execute-api.us-east-1.amazonaws.com/dev/"; \
		exit 1; \
	fi
	curl -X POST "$(API_URL)api/echo" \
		-H "Content-Type: application/json" \
		-d '{"message": "Hello, World!", "timestamp": "2023-01-01T00:00:00Z"}' | jq .

# クイックテスト（直接URL使用）
quick-test: ## 本番APIに対してクイックテストを実行
	@echo "=== Quick API Test ==="
	@echo "GET Test:"
	@curl -s -X GET "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/test?param=value" | jq .
	@echo -e "\nPOST Test:"
	@curl -s -X POST "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/api/echo" \
		-H "Content-Type: application/json" \
		-d '{"message": "Quick test"}' | jq .
	@echo -e "\nError Test:"
	@curl -s -X DELETE "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/test" | jq .