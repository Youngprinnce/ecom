package api

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/youngprinnce/go-ecom/controller/order"
	"github.com/youngprinnce/go-ecom/controller/product"
	"github.com/youngprinnce/go-ecom/controller/user"
	"github.com/youngprinnce/go-ecom/middleware"
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

	log.Println("Listening on", s.addr)

	return router.Run(s.addr)
}
