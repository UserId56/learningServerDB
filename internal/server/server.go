package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"learningServerDB/API/routers"
	"learningServerDB/internal/cfg"
	"learningServerDB/internal/handlers"
	"net/http"
	"time"
)

type Server struct {
	config *cfg.Cfg
	ctx    context.Context
	srv    *http.Server
	db     *pgxpool.Pool
}

func NewServer(ctx context.Context, config *cfg.Cfg) *Server {
	server := new(Server)
	server.ctx = ctx
	server.config = config
	return server

}

func (s *Server) Serv() {
	var err error
	s.db, err = pgxpool.New(context.Background(), s.config.GetDBString())
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	} else {
		fmt.Println("DB connected successfully!")
	}
	userHandler := handlers.NewUserHandler(s.db)
	router := routers.NewRouter(userHandler)
	s.srv = &http.Server{
		Addr:    ":" + s.config.PORT,
		Handler: router,
	}
	fmt.Println("Server is starting on port: " + s.config.PORT)
	err = s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Println("Failed to start server: " + err.Error())
		panic(err)
	}
	return
}

func (s *Server) Shutdown() {
	fmt.Println("Shutting down server...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	s.db.Close()
	defer func() {
		cancel()
	}()
	if err := s.srv.Shutdown(ctxShutdown); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}

}
