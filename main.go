package main

import (
	"fish/config"
	"fish/db"
	"fish/handlers"
	"log"
	"net/http"
	"time"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	db.InitDB(cfg)
	defer db.CloseDB()

	// 设置WebSocket路由
	http.HandleFunc("/ws", handlers.WebSocketHandler)

	// 启动HTTP服务器
	server := &http.Server{
		Addr:         cfg.ServerAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("服务器正在运行，监听地址: %s", cfg.ServerAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
