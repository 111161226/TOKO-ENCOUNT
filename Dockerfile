# --- ステージ1: Frontend Build ---
FROM node:16-alpine AS client-build
RUN apk add --update --no-cache openjdk8-jre-base
WORKDIR /client-build

ARG VITE_API_HOST
ENV VITE_API_HOST=$VITE_API_HOST

COPY ./client/package*.json ./
RUN npm ci --unsafe-perm
COPY ./client .
COPY ./docs ../docs
RUN npm run gen-api
RUN npm run build

# --- ステージ2: Backend Build ---
FROM golang:1.18-alpine AS server-build
RUN apk add --update --no-cache git
WORKDIR /go/src/github.com/111161226/TOKO-ENCOUNT

COPY ./go.* ./
RUN go mod download
COPY . .

# Swagger UI 準備
COPY docs/swagger.yaml /docs/swagger.yaml
COPY docs/swagger-ui-diff.patch /docs/swagger-ui-diff.patch
RUN git clone --depth 1 https://github.com/swagger-api/swagger-ui /docs/swagger-ui && \
    cd /docs/swagger-ui && \
    git apply /docs/swagger-ui-diff.patch

# バイナリ名を「toko-app」に固定[cite: 3, 4]
RUN go build -o toko-app main.go

# --- ステージ3: 実行環境 ---
FROM alpine:3.15.0
WORKDIR /app

RUN apk --update --no-cache add tzdata ca-certificates openssl && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    rm -rf /var/cache/apk/*

# 【修正】dockerizeが必要なら、ここで確実にパスを通す
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
 && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
 && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
 && chmod +x /usr/local/bin/dockerize

# ファイルのコピー
COPY --from=client-build /client-build/dist ./web/dist
COPY --from=server-build /docs/swagger-ui/dist ./docs/swagger-ui/dist
COPY --from=server-build /go/src/github.com/111161226/TOKO-ENCOUNT/docs/swagger.yaml ./docs/swagger.yaml
COPY --from=server-build /go/src/github.com/111161226/TOKO-ENCOUNT/toko-app .

# 【重要】Renderではdockerizeを介さず、バイナリを直接実行する設定にする
# ポート80固定ではなく、環境変数を参照するようにしたバイナリを叩く
ENTRYPOINT ["./toko-app"]