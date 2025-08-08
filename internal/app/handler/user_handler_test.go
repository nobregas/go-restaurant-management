package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-restaurant-management/internal/domain/user"
	"go-restaurant-management/internal/shared/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockUserService is a mock implementation of the UserService for testing.
type MockUserService struct {
	RegisterFunc func(user.User) (user.User, error)
}

func (m *MockUserService) Register(u user.User) (user.User, error) {
	return m.RegisterFunc(u)
}

func TestRegister(t *testing.T) {
	t.Run("should return 201 when user is registered successfully", func(t *testing.T) {
		// Create a mock user service
		mockUserService := &MockUserService{
			RegisterFunc: func(u user.User) (user.User, error) {
				// In this mock, we'll just return the user as is
				return u, nil
			},
		}

		// Create a new HTTP handler with the mock service
		h := UserHandler(mockUserService)

		// Create a new registration request
		regReq := types.RegisterUserRequest{
			First_name: "John",
			Last_name:  "Doe",
			Email:      "john.doe@example.com",
			Password:   "password",
			Phone:      "1234567890",
		}

		// Marshal the request body to JSON
		body, err := json.Marshal(regReq)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		h.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		// Check the response body
		var createdUser user.User
		if err := json.Unmarshal(rr.Body.Bytes(), &createdUser); err != nil {
			t.Fatal(err)
		}

		if createdUser.Email != regReq.Email {
			t.Errorf("handler returned unexpected body: got %v want %v",
				createdUser.Email, regReq.Email)
		}
	})

	t.Run("should return 400 when request body is invalid", func(t *testing.T) {
		// Create a mock user service
		mockUserService := &MockUserService{}

		// Create a new HTTP handler with the mock service
		h := UserHandler(mockUserService)

		// Create a new HTTP request with an invalid body
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"invalid`)))
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		h.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("should return 500 when user service returns an error", func(t *testing.T) {
		// Create a mock user service
		mockUserService := &MockUserService{
			RegisterFunc: func(u user.User) (user.User, error) {
				return user.User{}, errors.New("internal server error")
			},
		}

		// Create a new HTTP handler with the mock service
		h := UserHandler(mockUserService)

		// Create a new registration request
		regReq := types.RegisterUserRequest{
			First_name: "John",
			Last_name:  "Doe",
			Email:      "john.doe@example.com",
			Password:   "password",
			Phone:      "1234567890",
		}

		// Marshal the request body to JSON
		body, err := json.Marshal(regReq)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		h.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}
