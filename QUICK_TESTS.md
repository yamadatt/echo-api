# 🚀 Echo API クイックテストガイド

## 基本情報
- **API URL**: `https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/`
- **サポートメソッド**: GET, POST, OPTIONS
- **エラーメソッド**: DELETE, PUT, PATCH (405エラー)

## 🎯 最も使用頻度の高いテスト

### 1. 超簡単テスト
```bash
# 最もシンプルなGET
curl https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/

# JSON整形付き
curl -s https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/ | jq .
```

### 2. パラメータ付きテスト
```bash
# クエリパラメータ
curl -s "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/test?name=taro&age=25" | jq .
```

### 3. JSONデータ送信
```bash
# シンプルなJSON
curl -s -X POST "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/api/data" \
  -H "Content-Type: application/json" \
  -d '{"name": "太郎", "message": "こんにちは"}' | jq .
```

### 4. エラーテスト
```bash
# 405エラーの確認
curl -s -X DELETE "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/test" | jq .
```


## 💡 便利なワンライナー

### レスポンス時間測定
```bash
curl -w "Time: %{time_total}s\n" -s -o /dev/null \
  "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/test"
```

### HTTPヘッダー表示
```bash
curl -I "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/"
```

### 繰り返しテスト（5回）
```bash
for i in {1..5}; do
  echo "Test $i:"
  curl -s "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/test" | jq .message
done
```

## 🔧 トラブルシューティング

### jqがない場合
```bash
# jqなしでも見やすく
curl -s "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/" | python -m json.tool
```

### HTTPステータスコード確認
```bash
curl -w "HTTP Code: %{response_code}\n" -s -o /dev/null \
  "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/"
```

### 詳細なリクエスト情報
```bash
curl -v "https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod/debug" 2>&1 | head -20
```

## 📊 期待されるレスポンス

### GET成功 (200)
```json
{
  "request": {
    "method": "GET",
    "path": "/test",
    "headers": {...},
    "queryParams": {...},
    "timestamp": "2025-09-10T23:35:32Z"
  },
  "message": "Request successfully echoed",
  "processedAt": "2025-09-10T23:35:32Z"
}
```

### POST成功 (200)
```json
{
  "request": {
    "method": "POST",
    "path": "/api/echo",
    "headers": {...},
    "queryParams": {},
    "body": "{\"message\": \"data\"}",
    "timestamp": "2025-09-10T23:35:32Z"
  },
  "message": "Request successfully echoed",
  "processedAt": "2025-09-10T23:35:32Z"
}
```

### エラー (405)
```json
{
  "error": "Method Not Allowed",
  "message": "Only GET and POST methods are supported",
  "timestamp": "2025-09-10T23:35:32Z"
}
```

