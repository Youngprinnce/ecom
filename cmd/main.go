package main

import (
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/youngprinnce/go-ecom/cmd/api"
	"github.com/youngprinnce/go-ecom/config"
	"github.com/youngprinnce/go-ecom/db"
)

func main() {
	cfg := mysql.Config{
		User:                 config.Envs.DB.User,
		Passwd:               config.Envs.DB.Passwd,
		Net:                  config.Envs.DB.Net,
		Addr:                 config.Envs.DB.Addr,
		DBName:               config.Envs.DB.DBName,
		AllowNativePasswords: config.Envs.DB.AllowNativePasswords,
		ParseTime:            config.Envs.DB.ParseTime,
	}
	
	db, err := db.NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected")

	server := api.NewAPIServer(":8080", nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
