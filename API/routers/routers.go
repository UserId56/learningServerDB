package routers

import (
	"github.com/gorilla/mux"
	"learningServerDB/internal/handlers"
	"net/http"
)

func testRespons(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, this is a test response!"))
}

func NewRouter(userHandler *handlers.UserHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/test", testRespons).Methods("GET")
	router.HandleFunc("/user/{id}", userHandler.GetUserByID).Methods("GET")
	router.HandleFunc("/user/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/user/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/user/logout", userHandler.Logout).Methods("POST")
	router.HandleFunc("/user/refresh", userHandler.Refresh).Methods("POST")
	router.HandleFunc("/user", userHandler.Delete).Methods("DELETE")
	// Define your routes here
	// Example: router.HandleFunc("/api/resource", resourceHandler).Methods("GET", "POST")
	return router
}
