// db/db.go
package db

import (
	"database/sql"
	"fish/config"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(cfg *config.Config) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	// 测试数据库连接
	err = DB.Ping()
	if err != nil {
		log.Fatalf("无法ping数据库: %v", err)
	}

	log.Println("成功连接到数据库")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
