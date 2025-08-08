package handler

import (
	"go-restaurant-management/internal/domain/user"
	"go-restaurant-management/internal/shared/errors/exceptions"
	"go-restaurant-management/internal/shared/middleware"
	"go-restaurant-management/internal/shared/types"
	"go-restaurant-management/internal/shared/utils"
	"log"
	"net/http"
)

func AuthHandler(userService user.UserService) http.HandlerFunc {
	routes := map[string]map[string]middleware.HandlerFunc{
		"/api/auth/register": {
			"POST": func(w http.ResponseWriter, r *http.Request) error {
				return register(w, r, userService)
			},
		},
		// "/api/auth/login": {
		// 	"POST": func(w http.ResponseWriter, r *http.Request) error {
		// 		return login(w, r, userService)
		// 	},
		// },
	}

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		if pathRoutes, exists := routes[path]; exists {
			if handler, methodExists := pathRoutes[method]; methodExists {
				middleware.ErrorHandlerFunc(handler)(w, r)
				return
			}
			utils.WriteError(w, exceptions.NewMethodNotAllowedError(method, path))
			return
		}
		utils.WriteError(w, exceptions.NewRouteNotFoundError(path))
	}
}

func register(w http.ResponseWriter, r *http.Request, userService user.UserService) error {
	log.Println("-> new request to register user")
	var req types.RegisterUserRequest

	if err := utils.ParseAndValidateJson(r, &req); err != nil {
		log.Printf("error parsing json: %v", err)
		return err
	}

	if err := validateBusinessRules(&req, userService); err != nil {
		return err
	}

	user, err := userService.Register(user.RegisterToUser(req))
	if err != nil {
		log.Printf("error registering user: %v", err)
		return err
	}

	log.Printf("user %s registered successfully", user.Email)

	response := map[string]interface{}{
		"user":    user,
		"message": "User registered successfully",
	}

	utils.WriteJson(w, http.StatusCreated, response)
	return nil
}

func validateBusinessRules(req *types.RegisterUserRequest, userService user.UserService) error {
	// Check if user already exists by email
	_, err := userService.FindByEmail(req.Email)
	if err == nil {
		return exceptions.NewConflictError("email", "user already exists")
	}

	// Check if phone already exists (if necessary)
	// _, err = userService.FindByPhone(req.Phone)
	// if err == nil {
	// 	return exceptions.NewConflictError("phone", "phone number already in use")
	// }

	return nil
}
