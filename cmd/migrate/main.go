package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/youngprinnce/go-ecom/config"
	"github.com/youngprinnce/go-ecom/db"
)

func main() {
	cfg := mysqlDriver.Config{
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

	driver, err := mysqlMigrate.WithInstance(db, &mysqlMigrate.Config{})
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

	v, d, _ := m.Version()
	log.Printf("Version: %d, dirty: %v", v, d)

	cmd := os.Args[len(os.Args)-1]
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
