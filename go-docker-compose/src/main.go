package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
)

var (
	db  *sql.DB
	rdb *redis.Client
	ctx = context.Background()
)

func main() {
	// 读取环境变量
	mysqlURL := os.Getenv("MYSQL_URL")
	redisURL := os.Getenv("REDIS_URL")

	// 初始化 MySQL 连接
	var err error
	db, err = sql.Open("mysql", mysqlURL)
	if err != nil {
		log.Fatalf("Error opening MySQL connection: %v", err)
	}
	defer db.Close()

	// 初始化 Redis 连接
	// 使用redis.ParseURL函数解析Redis URI
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Error parsing Redis URL: %v", err)
	}

	// 创建Redis客户端
	rdb = redis.NewClient(opts)

	// 使用Redis客户端执行操作，例如Ping
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error pinging Redis: %v", err)
	}
	log.Println("Redis ping response:", pong)

	// 设置 HTTP 处理器
	http.HandleFunc("/", handleRequest)

	// 启动 HTTP 服务器
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 从 Redis 读取数据
	redisValue, err := rdb.Get(ctx, "key").Result()
	if err != nil && err != redis.Nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 从 MySQL 查询数据
	var now string
	err = db.QueryRow("SELECT NOW()").Scan(&now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 响应客户端
	fmt.Fprintf(w, "Current time from MySQL: %s\n", now)
	fmt.Fprintf(w, "Value from Redis: %s\n", redisValue)
}
