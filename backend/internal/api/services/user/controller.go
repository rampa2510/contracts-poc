package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/rampa2510/contracts-poc/internal/api/middleware"
	"github.com/rampa2510/contracts-poc/internal/utils"
)

type UserController struct {
	storage *UserStorage
}

func NewUserController(s *UserStorage) *UserController {
	slog.Info("Initalised User controller")
	return &UserController{storage: s}
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var user userDb

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(&middleware.APIError{
			Message: "Invalid request body",
			Status:  http.StatusBadRequest,
		})
	}

	userId, err := uc.storage.createUser(user.Name, user.Email)
	if err != nil {
		slog.Error("Error while creating new user", err)
		panic(&middleware.APIError{Message: "Internal Server Error", Status: http.StatusInternalServerError})
	}

	slog.Info("User created", slog.String("userId", userId))

	response := map[string]string{
		"newUserId": userId,
	}

	utils.SendResponse(w, http.StatusCreated, response)
}

func (uc *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uc.storage.getUsers()
	if err != nil {
		panic(&middleware.APIError{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
	}
	utils.SendResponse(w, http.StatusOK, users)
}

func (uc *UserController) GetAUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		panic(&middleware.APIError{
			Message: "Invalid User Id",
			Status:  http.StatusBadRequest,
		})
	}

	user, err := uc.storage.getUser(userId)
	if err != nil {
		panic(&middleware.APIError{
			Message: "Internal Server Error",
			Status:  http.StatusInternalServerError,
		})
	}
	utils.SendResponse(w, http.StatusOK, user)
}
