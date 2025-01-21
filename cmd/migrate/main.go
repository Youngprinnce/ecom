package main

import (
	"log"
	"os"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/youngprinnce/go-ecom/config"
	"github.com/youngprinnce/go-ecom/db"
)


func main() {
	cfg := mysqlCfg.Config{
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

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
		
}
