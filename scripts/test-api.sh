#!/bin/bash

# Echo API テストスクリプト
set -e

# カラー設定
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# API設定
API_URL="https://o5sqxqj3e2.execute-api.ap-northeast-1.amazonaws.com/prod"

echo -e "${GREEN}=== Echo API テストスクリプト ===${NC}"
echo "API URL: ${API_URL}"
echo ""

# テスト関数
run_test() {
    local test_name="$1"
    local curl_command="$2"
    local expected_status="$3"
    
    echo -e "${BLUE}=== ${test_name} ===${NC}"
    echo "Command: ${curl_command}"
    echo ""
    
    # レスポンスとステータスコードを取得
    response=$(eval "${curl_command} -w '\n%{http_code}'" -s)
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo "Status Code: ${status_code}"
    echo "Response:"
    
    # JSONかどうかをチェックして整形
    if echo "$body" | jq . >/dev/null 2>&1; then
        echo "$body" | jq .
    else
        echo "$body"
    fi
    
    # ステータスコードの検証
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ テスト成功${NC}"
    else
        echo -e "${RED}✗ テスト失敗 (Expected: $expected_status, Got: $status_code)${NC}"
    fi
    
    echo ""
    echo "---"
    echo ""
}

# 引数チェック
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "使用方法: $0 [オプション]"
    echo ""
    echo "オプション:"
    echo "  --all, -a       全テストを実行"
    echo "  --basic, -b     基本テストのみ実行"
    echo "  --errors, -e    エラーテストのみ実行"
    echo "  --performance   パフォーマンステストを実行"
    echo "  --help, -h      このヘルプを表示"
    echo ""
    echo "引数なしの場合は基本テストを実行します。"
    exit 0
fi

# テストの実行
case "${1:-basic}" in
    "all"|"-a"|"--all")
        echo -e "${YELLOW}全テストを実行します...${NC}"
        echo ""
        
        # 基本テスト
        run_test "ルートパスGET" \
            "curl -X GET '${API_URL}/'" \
            "200"
        
        run_test "クエリパラメータ付きGET" \
            "curl -X GET '${API_URL}/test?param1=value1&param2=value2' -H 'Content-Type: application/json'" \
            "200"
        
        run_test "複数ヘッダー付きGET" \
            "curl -X GET '${API_URL}/api/users' -H 'Content-Type: application/json' -H 'Authorization: Bearer test-token' -H 'X-Custom-Header: custom-value'" \
            "200"
        
        run_test "基本的なPOST" \
            "curl -X POST '${API_URL}/api/echo' -H 'Content-Type: application/json' -d '{\"message\": \"Hello, World!\", \"timestamp\": \"2023-01-01T00:00:00Z\"}'" \
            "200"
        
        run_test "空のPOSTボディ" \
            "curl -X POST '${API_URL}/empty' -H 'Content-Type: application/json'" \
            "200"
        
        run_test "DELETEメソッド（405エラー）" \
            "curl -X DELETE '${API_URL}/test'" \
            "405"
        
        run_test "PUTメソッド（405エラー）" \
            "curl -X PUT '${API_URL}/test' -H 'Content-Type: application/json' -d '{\"data\": \"test\"}'" \
            "405"
        
        run_test "OPTIONSリクエスト" \
            "curl -X OPTIONS '${API_URL}/test' -H 'Origin: https://example.com'" \
            "200"
        ;;
        
    "basic"|"-b"|"--basic")
        echo -e "${YELLOW}基本テストを実行します...${NC}"
        echo ""
        
        run_test "ルートパスGET" \
            "curl -X GET '${API_URL}/'" \
            "200"
        
        run_test "クエリパラメータ付きGET" \
            "curl -X GET '${API_URL}/test?param1=value1&param2=value2' -H 'Content-Type: application/json'" \
            "200"
        
        run_test "基本的なPOST" \
            "curl -X POST '${API_URL}/api/echo' -H 'Content-Type: application/json' -d '{\"message\": \"Hello, World!\"}'" \
            "200"
        ;;
        
    "errors"|"-e"|"--errors")
        echo -e "${YELLOW}エラーテストを実行します...${NC}"
        echo ""
        
        run_test "DELETEメソッド（405エラー）" \
            "curl -X DELETE '${API_URL}/test'" \
            "405"
        
        run_test "PUTメソッド（405エラー）" \
            "curl -X PUT '${API_URL}/test' -H 'Content-Type: application/json' -d '{\"data\": \"test\"}'" \
            "405"
        
        run_test "PATCHメソッド（405エラー）" \
            "curl -X PATCH '${API_URL}/test'" \
            "405"
        ;;
        
    "performance"|"--performance")
        echo -e "${YELLOW}パフォーマンステストを実行します...${NC}"
        echo ""
        
        echo -e "${BLUE}=== レスポンス時間測定 ===${NC}"
        curl -X GET "${API_URL}/performance" \
            -w "\nResponse Time: %{time_total}s\nHTTP Code: %{response_code}\nSize: %{size_download} bytes\n" \
            -s -o /dev/null
        echo ""
        
        run_test "大きなJSONペイロード" \
            "curl -X POST '${API_URL}/large' -H 'Content-Type: application/json' -d '{\"data\": \"$(printf 'A%.0s' {1..100})\", \"description\": \"Large payload test\"}'" \
            "200"
        ;;
        
    *)
        echo -e "${RED}不明なオプション: $1${NC}"
        echo "使用方法については $0 --help を参照してください。"
        exit 1
        ;;
esac

echo -e "${GREEN}=== テスト完了 ===${NC}"