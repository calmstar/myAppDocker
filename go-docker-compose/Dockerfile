FROM golang:1.19

WORKDIR /app

LABEL maintainer="cwx@email.com"

COPY ./src/ /app

#  || true 的作用是确保即使 go mod init myAppDocker 失败，Docker 构建过程也不会因为这一行而停止。
RUN go mod init myAppDocker.com || true

RUN go mod tidy

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]