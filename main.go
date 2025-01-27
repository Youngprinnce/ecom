package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/youngprinnce/go-ecom/cmd/api"
	"github.com/youngprinnce/go-ecom/config"
	"github.com/youngprinnce/go-ecom/db"

	_ "github.com/youngprinnce/go-ecom/docs"
)


//	@termsOfService	http://swagger.io/terms/

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

//	@securityDefinitions.apiKey	apiKey
//	@in							header
//	@name						Authorization
//	@description				JWT token for authentication

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

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.PORT), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
