# 构建阶段
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o tdl-app ./cmd/tdl

# 运行阶段
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/tdl-app .
EXPOSE 8080
ENTRYPOINT ["./tdl-app"]