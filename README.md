# Echo API

AWS LambdaとAPI Gatewayを使用したHTTPエコーサービスです。受信したHTTPリクエストの詳細情報∂をJSON形式で返却します。

## 概要

このプロジェクトは、開発者がAPIリクエストをテストし、デバッグするためのエコーサービスを提供します。GETとPOSTリクエストをサポートし、リクエストのメソッド、ヘッダー、クエリパラメータ、ボディなどの詳細情報をレスポンスとして返却します。

## 特徴

- **サーバーレス**: AWS LambdaとAPI Gatewayによる完全サーバーレス構成
- **軽量コンテナ**: マルチステージビルドによる最適化されたコンテナイメージ（331MB）
- **Go言語**: 高パフォーマンスで静的リンクされたバイナリ
- **構造化ログ**: JSON形式での包括的なログ記録
- **CORS対応**: クロスオリジンリクエストをサポート
- **エラーハンドリング**: 適切なHTTPステータスコードとエラーメッセージ

## サポートするHTTPメソッド

- **GET**: クエリパラメータとヘッダー情報をエコー
- **POST**: リクエストボディ、ヘッダー情報をエコー
- **OPTIONS**: CORS preflight リクエストをサポート

## レスポンス形式

### 成功レスポンス (200 OK)

```json
{
  "request": {
    "method": "GET",
    "path": "/test",
    "headers": {
      "Content-Type": "application/json",
      "User-Agent": "curl/7.68.0"
    },
    "queryParams": {
      "param1": "value1",
      "param2": "value2"
    },
    "body": "",
    "timestamp": "2023-01-01T12:00:00Z"
  },
  "message": "Request successfully echoed",
  "processedAt": "2023-01-01T12:00:01Z"
}
```

### エラーレスポンス (405 Method Not Allowed)

```json
{
  "error": "Method Not Allowed",
  "message": "Only GET and POST methods are supported",
  "timestamp": "2023-01-01T12:00:00Z"
}
```

## 必要な前提条件

- Go 1.21以上
  - https://go.dev/ 
- Docker
  - https://www.docker.com/ja-jp/
- AWS CLI
  - https://docs.aws.amazon.com/ja_jp/cli/v1/userguide/cli-chap-install.html
- SAM CLI
  - https://docs.aws.amazon.com/ja_jp/serverless-application-model/latest/developerguide/install-sam-cli.html
- jq (テスト用、オプション 無くてもかまいません。)
  - https://jqlang.org/


## ビルドとデプロイ

### Lambda関数のビルド

```bash
sam build --use-container
```

### デプロイ

```bash
sam deploy

```



### デプロイ後のテスト

デプロイが完了すると、API Gateway URLが表示されます。

```bash
# 環境変数にAPI URLを設定
export API_URL=https://your-api-id.execute-api.ap-northeast-1.amazonaws.com/prod/

```

## ログとモニタリング

### Lambda関数のログ

```bash

# または直接aws logsコマンド
aws logs tail /aws/lambda/echo-api-prodxxxxxxx --follow
```

## プロジェクト構造

```
echo-api/
├── cmd/lambda/           # Lambda関数のエントリーポイント
│   └── main.go
├── internal/
│   ├── handler/          # Lambda ハンドラー
│   │   ├── lambda.go
│   │   └── lambda_test.go
│   └── models/           # データモデル
│       ├── echo.go
│       └── echo_test.go
├── pkg/logger/           # ログユーティリティ
│   ├── logger.go
│   └── logger_test.go
├── scripts/              # ビルド・デプロイスクリプト
│   ├── build.sh
│   ├── deploy.sh
│   └── test-api.sh
├── Dockerfile            # 軽量コンテナイメージ定義（マルチステージビルド）
├── Dockerfile.original   # 元のDockerfile（バックアップ）
├── template.yaml         # SAMテンプレート（OpenAPI統合付き）
├── samconfig.toml        # SAM設定（ECR自動作成）
├── Makefile              # ビルドコマンド（テスト除く）
├── go.mod                # Go modules
├── go.sum                # Go dependencies
├── QUICK_TESTS.md        # クイックテストガイド
└── README.md
```

## 技術仕様

### アーキテクチャ
- **AWS Lambda**: Go 1.21、コンテナイメージ実行
- **API Gateway**: REST API、プロキシ統合
- **Amazon ECR**: コンテナイメージストレージ
- **CloudWatch**: ログとモニタリング

### パフォーマンス
- **メモリ**: 128MB
- **タイムアウト**: 30秒
- **コールドスタート**: 通常1-2秒

### セキュリティ
- **CORS設定**: すべてのオリジンに対応
- **認証**: なし（パブリックAPI）
- **HTTPS**: API Gateway経由で自動対応

## 現在のエンドポイント

```
Base URL: https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/
```

∂
