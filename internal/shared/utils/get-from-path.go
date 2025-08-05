package utils

import (
	"go-restaurant-management/internal/shared/errors/exceptions"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetIntParamFromPath(r *http.Request, param string) int {
	vars := mux.Vars(r)

	str, ok := vars[param]
	if !ok {
		panic(exceptions.NewValidationError(param, "parameter not found"))
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		panic(exceptions.NewValidationError(param, "parameter not found"))
	}

	return val
}
