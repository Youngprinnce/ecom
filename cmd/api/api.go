package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/youngprinnce/go-ecom/controller/user"
)

type APIServer struct {
	addr string
	db *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Run() error {
	router := s.setupRouter()

	log.Println("Server is running on", s.addr)
	return http.ListenAndServe(s.addr, router)
}

func (s *APIServer) setupRouter() *mux.Router {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userController := user.NewHandler()
	userController.RegisterRoutes(subrouter)

	return router
}
