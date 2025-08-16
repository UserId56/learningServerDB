package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strconv"
)

type UserHandler struct {
	dataBasePool *pgxpool.Pool
}

func NewUserHandler(dbPool *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		dataBasePool: dbPool,
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	var user string
	variable := mux.Vars(r)
	if variable["id"] == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(variable["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid User ID format", http.StatusBadRequest)
		return
	}
	err = h.dataBasePool.QueryRow(context.Background(), "SELECT username FROM users WHERE id = $1", id).Scan(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User ID: " + variable["id"] + ", Username: " + user))
}
