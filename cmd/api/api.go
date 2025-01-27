package api

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/youngprinnce/go-ecom/controller/order"
	"github.com/youngprinnce/go-ecom/controller/product"
	"github.com/youngprinnce/go-ecom/controller/user"
	"github.com/youngprinnce/go-ecom/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/youngprinnce/go-ecom/docs" // Import the generated docs
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

// @title Go E-Commerce API
// @version 1.0
// @description This is a sample e-commerce API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@ecom.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func (s *APIServer) Run() error {
	router := gin.Default()

	// Apply the Logging middleware to all routes
	router.Use(middleware.Logging())

	api := router.Group("/api/v1")

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(api)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(api)

	orderStore := order.NewStore(s.db)
	orderHandler := order.NewHandler(productStore, orderStore, userStore)
	orderHandler.RegisterRoutes(api)

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Listening on", s.addr)

	return router.Run(s.addr)
}
