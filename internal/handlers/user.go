package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"learningServerDB/internal/cfg"
	"learningServerDB/internal/logger"
	"learningServerDB/internal/models"
	"learningServerDB/internal/utils"
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
		logger.LogError(err, "Ошибка парса тела запроса:", logger.Error)
		utils.RespondWithError(w, http.StatusBadRequest, utils.ErrCodeServerError, "Некорректное тело запроса")
		return
	}
	err = userData.Validate()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.ErrCodeInvalidRequest, fmt.Sprintf("Ошибка валидации: %v", err))
		return
	}
	passwordHash, err := userData.GetHashPassword()
	if err != nil {
		logger.LogError(err, "Ошибка хеширования пароля:", logger.Error)
		utils.RespondWithError(w, http.StatusInternalServerError, utils.ErrCodeServerError, "Ошибка на сервере")
		return
	}
	tx, err := h.dataBasePool.Begin(context.Background())
	if err != nil {
		logger.LogError(err, "Ошибка начала транзакции:", logger.Error)
		utils.RespondWithError(w, http.StatusInternalServerError, utils.ErrCodeServerError, "Ошибка на сервере")
		return
	}
	defer tx.Rollback(context.Background())

	var userId int64
	err = tx.QueryRow(context.Background(),
		"INSERT INTO users (username, email, passwordhash, createdat, updatedat) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id",
		userData.Username, userData.Email, passwordHash,
	).Scan(&userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				utils.RespondWithError(w, http.StatusConflict, utils.ErrCodeEmailExists, fmt.Sprintf("Email %s уже используется", userData.Email))
				return
			case "users_username_key":
				utils.RespondWithError(w, http.StatusConflict, utils.ErrCodeUsernameExists, fmt.Sprintf("Username %s уже используется", userData.Username))
				return
			default:
				utils.RespondWithError(w, http.StatusConflict, utils.ErrCodeUserAlreadyExists, "Пользователь с такими данными уже существует")
				return
			}
		}
		logger.LogError(err, "Ошибка при создании пользователя:", logger.Error)
		utils.RespondWithError(w, http.StatusInternalServerError, utils.ErrCodeServerError, "Ошибка на сервере")
		return

	}
	var userProfilesId int64
	err = tx.QueryRow(context.Background(),
		"INSERT INTO userprofiles (userid, firstname, lastname, middlename, createdat, updatedat) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id",
		userId, userData.FirstName, userData.LastName, userData.MiddleName,
	).Scan(&userProfilesId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, utils.ErrCodeServerError, "Ошибка на сервере")
		logger.LogError(err, "Ошибка при создании профиля пользователя:", logger.Error)
		return
	}
	var refreshToken string
	err = tx.QueryRow(context.Background(),
		"INSERT INTO refresh_tokens (userid, createdat, expiresat) VALUES ($1, NOW(), $2) RETURNING token",
		userId, time.Now().Add(time.Hour*24*time.Duration(h.config.REFRESH_TOKEN_LIVE_DAY))).Scan(&refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, utils.ErrCodeServerError, "Ошибка на сервере")
		logger.LogError(err, "Ошибка при создании refresh токена:", logger.Error)
		return
	}

	if err = tx.Commit(context.Background()); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, utils.ErrCodeServerError, "Ошибка на сервере")
		logger.LogError(err, "Ошибка фиксации транзакции:", logger.Error)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"username": userData.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.config.SECRET))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, utils.ErrCodeServerError, "Ошибка на сервере")
		logger.LogError(err, "Ошибка создания JWT токена:", logger.Error)
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
		logger.LogError(err, "Ошибка кодирования ответа:", logger.Error)
		utils.RespondWithError(w, http.StatusInternalServerError, utils.ErrCodeServerError, "Ошибка на сервере")
		return
	}
}
