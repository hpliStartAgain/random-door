package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	_ = godotenv.Load("../.env") // load .env from root

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "117.72.99.55"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "53306"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "123456"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pass, host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open mysql: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping mysql: %v", err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS city_roam CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	fmt.Println("Database city_roam created or already exists.")
}
