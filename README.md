# Echo API

AWS LambdaとAPI Gatewayを使用したHTTPエコーサービスです。受信したHTTPリクエストの詳細情報をJSON形式で返却します。

## 概要

このプロジェクトは、開発者がAPIリクエストをテストし、デバッグするためのエコーサービスを提供します。GETとPOSTリクエストをサポートし、リクエストのメソッド、ヘッダー、クエリパラメータ、ボディなどの詳細情報をレスポンスとして返却します。

## 特徴

- **サーバーレス**: AWS LambdaとAPI Gatewayによる完全サーバーレス構成
- **軽量コンテナ**: マルチステージビルドによる最適化されたコンテナイメージ（331MB）
- **Go言語**: 高パフォーマンスで静的リンクされたバイナリ
- **構造化ログ**: JSON形式での包括的なログ記録
- **CORS対応**: クロスオリジンリクエストをサポート
- **プロキシ統合**: API GatewayとLambdaの完全プロキシ統合
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
- Docker
- AWS CLI
- SAM CLI
- jq (テスト用、オプション)

## セットアップ

### 1. 依存関係のインストール

```bash
make install-tools
```



## ビルドとデプロイ

### Lambda関数のビルド

```bash
sam build
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

### API Gatewayのログ

```bash
make logs-api
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

## トラブルシューティング

### よくある問題

1. **ECRログインエラー**
   ```bash
   aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com
   ```

2. **IAM権限エラー**
   - Lambda実行権限
   - ECRアクセス権限
   - CloudFormation操作権限が必要

3. **SAM buildエラー**
   ```bash
   # キャッシュをクリア
   make clean
   make build
   ```

4. **API Gateway統合エラー**
   - プロキシ統合が正しく設定されていることを確認
   - OpenAPI定義でhttpMethod: POSTが指定されていることを確認

5. **コンテナイメージが大きい場合**
   - 軽量版Dockerfileを使用（マルチステージビルド）
   - ECRプッシュ時間を短縮

### ログレベルの変更

環境変数でログレベルを調整できます：

```bash
# template.yamlまたはsamconfig.tomlで設定
LOG_LEVEL=DEBUG  # DEBUG, INFO, WARN, ERROR
```

## 技術仕様

### アーキテクチャ
- **AWS Lambda**: Go 1.21、コンテナイメージ実行
- **API Gateway**: REST API、プロキシ統合
- **Amazon ECR**: コンテナイメージストレージ
- **CloudWatch**: ログとモニタリング

### パフォーマンス
- **コンテナサイズ**: 331MB（マルチステージビルド最適化）
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

### 利用可能なパス
- `GET /test` - テスト用エンドポイント
- `POST /api/echo` - エコー専用エンドポイント
- `ANY /{proxy+}` - すべてのパスに対応

## ライセンス

MIT License

## 貢献

プルリクエストやIssueの投稿をお待ちしています。

## サポート

問題が発生した場合は、以下を確認してください：

1. CloudWatch Logsでのエラーログ
2. API GatewayとLambda関数の設定
3. IAM権限の確認