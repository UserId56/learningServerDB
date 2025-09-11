package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"learningServerDB/internal/cfg"
	"learningServerDB/internal/models"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	dataBasePool *pgxpool.Pool
	config       *cfg.Cfg
}

func NewUserHandler(dbPool *pgxpool.Pool, config *cfg.Cfg) *UserHandler {
	return &UserHandler{
		dataBasePool: dbPool,
		config:       config,
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

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userData models.UserRegistrationData
	var err error
	err = json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Не валидное тело запроса", http.StatusBadRequest)
		return
	}
	passwordHash, err := userData.GetHashPassword()
	if err != nil {
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		fmt.Printf("Ошибка хеширования пароля: %v", err)
		return
	}
	tx, err := h.dataBasePool.Begin(context.Background())
	if err != nil {
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		fmt.Printf("Ошибка создания транзакции: %v", err)
		return
	}
	defer tx.Rollback(context.Background())

	var userId int64
	err = tx.QueryRow(context.Background(),
		"INSERT INTO users (username, email, passwordhash, createdat, updatedat) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id",
		userData.Username, userData.Email, passwordHash,
	).Scan(&userId)
	if err != nil {
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		fmt.Printf("Ошибка при создании пользователя: %v", err)
		return
	}
	var userProfilesId int64
	err = tx.QueryRow(context.Background(),
		"INSERT INTO userprofiles (userid, firstname, lastname, middlename, createdat, updatedat) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id",
		userId, userData.FirstName, userData.LastName, userData.MiddleName,
	).Scan(&userProfilesId)
	if err != nil {
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		fmt.Printf("Ошибка при создании профиля пользователя: %v", err)
		return
	}
	var refreshToken string
	err = tx.QueryRow(context.Background(),
		"INSERT INTO refresh_tokens (userid, createdat, expiresat) VALUES ($1, NOW(), $2) RETURNING token",
		userId, time.Now().Add(time.Hour*24*time.Duration(h.config.REFRESH_TOKEN_LIVE_DAY))).Scan(&refreshToken)
	if err != nil {
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		fmt.Printf("Ошибка при создании refresh токена: %v", err)
		return
	}

	if err = tx.Commit(context.Background()); err != nil {
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		fmt.Printf("Ошибка при фиксации транзакции: %v", err)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"username": userData.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.config.SECRET))
	if err != nil {
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		fmt.Printf("Ошибка при создании JWT токена: %v", err)
		return
	}
	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := models.AuthResponse{
		UserID:         userId,
		AccessToken:    tokenString,
		RefreshToken:   refreshToken,
		UserProfilesID: userProfilesId,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Printf("Ошибка парса response или отправки ответа: %v", err)
		return
	}
	fmt.Printf("Result: %+v\n", userId)
	fmt.Printf("Result: %+v\n", userProfilesId)
	//fmt.Printf("Received user data: %+v\n", userData)
}
