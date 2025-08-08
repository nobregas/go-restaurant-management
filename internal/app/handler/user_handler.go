package handler

import (
	"go-restaurant-management/internal/domain/user"
	"go-restaurant-management/internal/shared/types"
	"go-restaurant-management/internal/shared/utils"
	"net/http"
)

func UserHandler(userService user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			register(w, r, userService)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func register(w http.ResponseWriter, r *http.Request, userService user.UserService) {
	var req types.RegisterUserRequest

	if err := utils.ParseJson(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := userService.Register(user.ToUser(req))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, user)
}
