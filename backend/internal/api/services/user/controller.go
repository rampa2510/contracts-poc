package user

import (
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

type createUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest

	err, validationErrors := utils.ValidateRequest(r, &req)

	if err != nil {
		panic(&middleware.APIError{
			Message: "Invalid request body",
			Status:  http.StatusBadRequest,
		})
	}

	if len(validationErrors) > 0 {
		utils.SendResponse(w, http.StatusBadRequest, validationErrors)
		return
	}

	userId, err := uc.storage.createUser(req.Name, req.Email)
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
