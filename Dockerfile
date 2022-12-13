FROM node:16-alpine AS client-build
RUN apk add --update --no-cache openjdk8-jre-base
WORKDIR /github.com/cs-sysimpl/SakataKintoki/client
COPY ./client/package*.json ./
COPY ./client/scripts ./scripts
COPY ./docs ../docs
RUN npm ci --unsafe-perm
COPY ./client .
RUN npm run gen-api
RUN npm run build

FROM golang:1.17-alpine AS server-build
RUN apk add --update --no-cache git curl make
WORKDIR /go/src/github.com/cs-sysimpl/SakataKintoki
COPY ./go.* ./
RUN go mod download
COPY . .
RUN go build

FROM alpine:3.15.0
WORKDIR /app
RUN apk --update --no-cache add tzdata ca-certificates openssl && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
  rm -rf /var/cache/apk/*
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
 && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
 && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz
COPY --from=client-build /github.com/cs-sysimpl/SakataKintoki/client/dist ./web/dist
COPY --from=server-build /go/src/github.com/cs-sysimpl/SakataKintoki ./
ENTRYPOINT ./SakataKintoki
