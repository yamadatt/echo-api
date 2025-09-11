#!/bin/bash

# ローカルテストスクリプト for Echo API
set -e

# カラー設定
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Echo API Local Test Script ===${NC}"

# ポート設定
PORT=${PORT:-3000}
API_URL="http://localhost:${PORT}"

echo "Local API URL: ${API_URL}"

# SAM local start-apiをバックグラウンドで起動
echo -e "${YELLOW}SAM local start-api を起動中...${NC}"
sam local start-api --port ${PORT} &
SAM_PID=$!

# SAMの起動を待つ
echo -e "${YELLOW}API の起動を待機中...${NC}"
sleep 10

# サーバーが起動しているかチェック
check_server() {
    curl -s "${API_URL}" > /dev/null 2>&1
    return $?
}

# 最大30秒待機
for i in {1..30}; do
    if check_server; then
        echo -e "${GREEN}API サーバーが起動しました${NC}"
        break
    fi
    echo "待機中... ($i/30)"
    sleep 1
done

if ! check_server; then
    echo -e "${RED}API サーバーの起動に失敗しました${NC}"
    kill $SAM_PID 2>/dev/null
    exit 1
fi

# テスト関数
run_test() {
    local test_name="$1"
    local curl_command="$2"
    local expected_status="$3"
    
    echo -e "\n${BLUE}=== ${test_name} ===${NC}"
    echo "Command: ${curl_command}"
    
    response=$(eval $curl_command)
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo "Status Code: ${status_code}"
    echo "Response Body:"
    echo "$body" | jq . 2>/dev/null || echo "$body"
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ テスト成功${NC}"
    else
        echo -e "${RED}✗ テスト失敗 (Expected: $expected_status, Got: $status_code)${NC}"
    fi
}

echo -e "\n${YELLOW}=== API テスト開始 ===${NC}"

# テスト1: GETリクエスト（クエリパラメータ付き）
run_test "GET Request with Query Parameters" \
    "curl -s -w '\n%{http_code}' -X GET '${API_URL}/test?param1=value1&param2=value2' -H 'Content-Type: application/json' -H 'User-Agent: test-client'" \
    "200"

# テスト2: POSTリクエスト（JSONボディ付き）
run_test "POST Request with JSON Body" \
    "curl -s -w '\n%{http_code}' -X POST '${API_URL}/api/echo' -H 'Content-Type: application/json' -d '{\"message\": \"Hello, World!\", \"timestamp\": \"2023-01-01T00:00:00Z\", \"data\": {\"key\": \"value\"}}'" \
    "200"

# テスト3: ルートパスGETリクエスト
run_test "GET Request to Root Path" \
    "curl -s -w '\n%{http_code}' -X GET '${API_URL}/' -H 'Accept: application/json'" \
    "200"

# テスト4: 無効なメソッド（405エラー）
run_test "Invalid Method (DELETE) - Should Return 405" \
    "curl -s -w '\n%{http_code}' -X DELETE '${API_URL}/test'" \
    "405"

# テスト5: 無効なメソッド（PUT）
run_test "Invalid Method (PUT) - Should Return 405" \
    "curl -s -w '\n%{http_code}' -X PUT '${API_URL}/api/test' -H 'Content-Type: application/json' -d '{\"data\": \"test\"}'" \
    "405"

# テスト6: OPTIONSリクエスト（CORS preflight）
run_test "OPTIONS Request (CORS Preflight)" \
    "curl -s -w '\n%{http_code}' -X OPTIONS '${API_URL}/test' -H 'Origin: https://example.com'" \
    "200"

# テスト7: 複雑なパス
run_test "Complex Path with Multiple Segments" \
    "curl -s -w '\n%{http_code}' -X GET '${API_URL}/api/v1/echo/test?debug=true&format=json' -H 'Authorization: Bearer test-token' -H 'X-Custom-Header: custom-value'" \
    "200"

# テスト8: 空のPOSTボディ
run_test "POST Request with Empty Body" \
    "curl -s -w '\n%{http_code}' -X POST '${API_URL}/empty' -H 'Content-Type: application/json'" \
    "200"

echo -e "\n${GREEN}=== テスト完了 ===${NC}"

# SAMプロセスを終了
echo -e "${YELLOW}SAM local プロセスを終了中...${NC}"
kill $SAM_PID 2>/dev/null
wait $SAM_PID 2>/dev/null

echo -e "${GREEN}ローカルテストが完了しました${NC}"