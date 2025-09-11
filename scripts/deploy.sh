#!/bin/bash

# デプロイスクリプト for Echo API
set -e

# カラー設定
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# デフォルト設定
ENVIRONMENT=${1:-"prod"}
AWS_REGION=${AWS_REGION:-"ap-northeast-1"}
REPOSITORY_NAME="echo-api"

echo -e "${GREEN}=== Echo API Deploy Script ===${NC}"
echo "Environment: ${ENVIRONMENT}"
echo "AWS Region: ${AWS_REGION}"

# 引数チェック
if [[ ! "$ENVIRONMENT" =~ ^(prod)$ ]]; then
    echo -e "${RED}エラー: 無効な環境です。prod を指定してください${NC}"
    echo "使用方法: $0 [prod]"
    exit 1
fi

# AWS Account IDを取得
echo -e "${YELLOW}AWS Account IDを取得中...${NC}"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
if [ $? -ne 0 ]; then
    echo -e "${RED}AWS Account IDの取得に失敗しました${NC}"
    exit 1
fi

ECR_URI="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${REPOSITORY_NAME}"

echo "AWS Account ID: ${AWS_ACCOUNT_ID}"
echo "ECR Repository: ${ECR_URI}"

# ECRイメージが存在するかチェック
echo -e "${YELLOW}ECRイメージをチェック中...${NC}"
aws ecr describe-images --repository-name ${REPOSITORY_NAME} --image-ids imageTag=latest --region ${AWS_REGION} > /dev/null 2>&1

if [ $? -ne 0 ]; then
    echo -e "${RED}ECRに最新のイメージが見つかりません。${NC}"
    echo -e "${YELLOW}まず build.sh を実行してください${NC}"
    exit 1
fi

# SAM build
echo -e "${YELLOW}SAM build を実行中...${NC}"
sam build --config-env ${ENVIRONMENT}
if [ $? -ne 0 ]; then
    echo -e "${RED}SAM buildに失敗しました${NC}"
    exit 1
fi

# SAM deploy
echo -e "${YELLOW}SAM deploy を実行中...${NC}"
sam deploy \
    --config-env ${ENVIRONMENT} \
    --image-repositories EchoFunction=${ECR_URI} \
    --capabilities CAPABILITY_IAM \
    --parameter-overrides \
        Environment=${ENVIRONMENT} \
        ECRRepository=${ECR_URI}:latest

if [ $? -eq 0 ]; then
    echo -e "${GREEN}=== デプロイ完了 ===${NC}"
    
    # スタック情報を取得
    STACK_NAME="echo-api-prod"
    echo -e "${BLUE}=== デプロイ情報 ===${NC}"
    
    # API Gateway URLを取得
    API_URL=$(aws cloudformation describe-stacks \
        --stack-name ${STACK_NAME} \
        --region ${AWS_REGION} \
        --query 'Stacks[0].Outputs[?OutputKey==`EchoApiUrl`].OutputValue' \
        --output text)
    
    if [ -n "$API_URL" ]; then
        echo -e "${GREEN}API Gateway URL: ${API_URL}${NC}"
        echo ""
        echo -e "${YELLOW}テスト用コマンド:${NC}"
        echo "# GETリクエスト"
        echo "curl -X GET \"${API_URL}test?param1=value1&param2=value2\" -H \"Content-Type: application/json\""
        echo ""
        echo "# POSTリクエスト"
        echo "curl -X POST \"${API_URL}api/echo\" -H \"Content-Type: application/json\" -d '{\"message\": \"Hello, World!\", \"timestamp\": \"2023-01-01T00:00:00Z\"}'"
        echo ""
        echo "# 無効なメソッドテスト (405エラー)"
        echo "curl -X DELETE \"${API_URL}test\""
    else
        echo -e "${YELLOW}API Gateway URLの取得に失敗しました${NC}"
    fi
    
else
    echo -e "${RED}デプロイに失敗しました${NC}"
    exit 1
fi