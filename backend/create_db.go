package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env") // load .env from root

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		log.Fatal("DB_USER is required")
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		log.Fatal("DB_PASSWORD is required")
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "city_roam"
	}
	if !regexp.MustCompile(`^[A-Za-z0-9_]+$`).MatchString(name) {
		log.Fatalf("invalid DB_NAME %q", name)
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

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", name))
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	fmt.Printf("Database %s created or already exists.\n", name)
}
