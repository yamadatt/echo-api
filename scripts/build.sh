#!/bin/bash

# ビルドスクリプト for Echo API
set -e

# カラー設定
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 設定
REPOSITORY_NAME="echo-api"
AWS_REGION=${AWS_REGION:-"ap-northeast-1"}
AWS_ACCOUNT_ID=${AWS_ACCOUNT_ID:-""}

echo -e "${GREEN}=== Echo API Build Script ===${NC}"

# AWS Account IDを取得
if [ -z "$AWS_ACCOUNT_ID" ]; then
    echo -e "${YELLOW}AWS Account IDを取得中...${NC}"
    AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
    if [ $? -ne 0 ]; then
        echo -e "${RED}AWS Account IDの取得に失敗しました${NC}"
        exit 1
    fi
fi

ECR_URI="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${REPOSITORY_NAME}"

echo "AWS Account ID: ${AWS_ACCOUNT_ID}"
echo "AWS Region: ${AWS_REGION}"
echo "ECR Repository: ${ECR_URI}"

# ECRリポジトリが存在するかチェック
echo -e "${YELLOW}ECRリポジトリをチェック中...${NC}"
aws ecr describe-repositories --repository-names ${REPOSITORY_NAME} --region ${AWS_REGION} > /dev/null 2>&1

if [ $? -ne 0 ]; then
    echo -e "${YELLOW}ECRリポジトリが存在しません。作成中...${NC}"
    aws ecr create-repository --repository-name ${REPOSITORY_NAME} --region ${AWS_REGION}
    if [ $? -ne 0 ]; then
        echo -e "${RED}ECRリポジトリの作成に失敗しました${NC}"
        exit 1
    fi
    echo -e "${GREEN}ECRリポジトリを作成しました${NC}"
fi

# ECRにログイン
echo -e "${YELLOW}ECRにログイン中...${NC}"
aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${ECR_URI}
if [ $? -ne 0 ]; then
    echo -e "${RED}ECRログインに失敗しました${NC}"
    exit 1
fi

# Dockerイメージをビルド
echo -e "${YELLOW}Dockerイメージをビルド中...${NC}"
docker build -t ${REPOSITORY_NAME} .
if [ $? -ne 0 ]; then
    echo -e "${RED}Dockerビルドに失敗しました${NC}"
    exit 1
fi

# イメージにタグを付与
echo -e "${YELLOW}イメージにタグを付与中...${NC}"
docker tag ${REPOSITORY_NAME}:latest ${ECR_URI}:latest

# 追加のタグ（日時）
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
docker tag ${REPOSITORY_NAME}:latest ${ECR_URI}:${TIMESTAMP}

# ECRにプッシュ
echo -e "${YELLOW}ECRにプッシュ中...${NC}"
docker push ${ECR_URI}:latest
docker push ${ECR_URI}:${TIMESTAMP}

if [ $? -eq 0 ]; then
    echo -e "${GREEN}=== ビルド完了 ===${NC}"
    echo "イメージURI: ${ECR_URI}:latest"
    echo "タイムスタンプタグ: ${ECR_URI}:${TIMESTAMP}"
    
    # samconfig.tomlを更新するためのECR URIを出力
    echo -e "${YELLOW}samconfig.tomlの更新用URI:${NC}"
    echo "EchoFunction=${ECR_URI}"
else
    echo -e "${RED}ECRプッシュに失敗しました${NC}"
    exit 1
fi