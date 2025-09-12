# 軽量化Dockerfileの例（参考用）
# マルチステージビルド
FROM --platform=linux/amd64 golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 静的リンクで完全に独立したバイナリを作成
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags '-extldflags "-static"' \
    -o bootstrap ./cmd/lambda

# Lambda Runtime提供のベースイメージを使用
FROM --platform=linux/amd64 public.ecr.aws/lambda/provided:al2

# バイナリをコピー
COPY --from=builder /app/bootstrap /var/runtime/

# 実行権限を設定
RUN chmod +x /var/runtime/bootstrap

# Lambda関数のハンドラーを設定
CMD ["bootstrap"]
