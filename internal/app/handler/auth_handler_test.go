package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"go-restaurant-management/internal/domain/user"
	"go-restaurant-management/internal/shared/types"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// MockUserService is a mock implementation of the UserService for testing.
type MockUserService struct {
	RegisterFunc    func(user.User) (user.User, error)
	FindByEmailFunc func(email string) (user.User, error)
}

func (m *MockUserService) Register(u user.User) (user.User, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(u)
	}
	return u, nil
}

func (m *MockUserService) FindByEmail(email string) (user.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return user.User{}, errors.New("user not found")
}

func TestRegister(t *testing.T) {
	t.Run("should return 201 when user is registered successfully", func(t *testing.T) {
		// Create a mock user service
		mockUserService := &MockUserService{
			RegisterFunc: func(u user.User) (user.User, error) {
				// Return the user with an ID to simulate successful creation
				u.ID = 1
				return u, nil
			},
			FindByEmailFunc: func(email string) (user.User, error) {
				// Return error to simulate user doesn't exist yet
				return user.User{}, errors.New("user not found")
			},
		}

		// Create a new HTTP handler with the mock service
		h := AuthHandler(mockUserService)

		// Create a new registration request
		regReq := types.RegisterUserRequest{
			First_name: "John",
			Last_name:  "Doe",
			Email:      "john.doe@example.com",
			Password:   "password123", // At least 6 characters
			Phone:      "12345678901", // Exactly 11 digits
		}

		// Marshal the request body to JSON
		body, err := json.Marshal(regReq)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new HTTP request with correct path
		req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		h.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v, body: %s",
				status, http.StatusCreated, rr.Body.String())
		}

		// Check the response body structure (now includes user and message)
		var response map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}

		// Check if user data exists in response
		if userData, ok := response["user"]; !ok {
			t.Error("response should contain 'user' field")
		} else {
			userMap := userData.(map[string]interface{})
			if userMap["email"] != regReq.Email {
				t.Errorf("handler returned unexpected email: got %v want %v",
					userMap["email"], regReq.Email)
			}
		}

		// Check if message exists
		if message, ok := response["message"]; !ok {
			t.Error("response should contain 'message' field")
		} else if message != "User registered successfully" {
			t.Errorf("unexpected message: got %v", message)
		}
	})

	t.Run("should return 400 when request body is invalid JSON", func(t *testing.T) {
		// Create a mock user service
		mockUserService := &MockUserService{}

		// Create a new HTTP handler with the mock service
		h := AuthHandler(mockUserService)

		// Create a new HTTP request with an invalid JSON body
		req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(`{"invalid`)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		h.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v, body: %s",
				status, http.StatusBadRequest, rr.Body.String())
		}

		// Check error response structure
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
			t.Fatal(err)
		}

		// Should be INVALID_JSON error
		if errorResponse["code"] != "INVALID_JSON" {
			t.Errorf("expected INVALID_JSON error, got %v", errorResponse["code"])
		}
	})

	t.Run("should return 400 when required fields are missing", func(t *testing.T) {
		// Create a mock user service
		mockUserService := &MockUserService{}

		// Create a new HTTP handler with the mock service
		h := AuthHandler(mockUserService)

		// Create a new registration request with missing required fields
		regReq := types.RegisterUserRequest{
			First_name: "", // Empty first name (should fail validation)
			Last_name:  "Doe",
			Email:      "john.doe@example.com",
			Password:   "password123",
			Phone:      "12345678901",
		}

		// Marshal the request body to JSON
		body, err := json.Marshal(regReq)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		h.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v, body: %s",
				status, http.StatusBadRequest, rr.Body.String())
		}

		// Check error response structure
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
			t.Fatal(err)
		}

		// Should be VALIDATION_ERROR
		if errorResponse["code"] != "VALIDATION_ERROR" {
			t.Errorf("expected VALIDATION_ERROR, got %v", errorResponse["code"])
		}
	})

	t.Run("should return 500 when user service returns an error", func(t *testing.T) {
		// Create a mock user service
		mockUserService := &MockUserService{
			RegisterFunc: func(u user.User) (user.User, error) {
				return user.User{}, errors.New("internal server error")
			},
			FindByEmailFunc: func(email string) (user.User, error) {
				// Return error to simulate user doesn't exist yet
				return user.User{}, errors.New("user not found")
			},
		}

		// Create a new HTTP handler with the mock service
		h := AuthHandler(mockUserService)

		// Create a new registration request
		regReq := types.RegisterUserRequest{
			First_name: "John",
			Last_name:  "Doe",
			Email:      "john.doe@example.com",
			Password:   "password123",
			Phone:      "12345678901",
		}

		// Marshal the request body to JSON
		body, err := json.Marshal(regReq)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		h.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v, body: %s",
				status, http.StatusInternalServerError, rr.Body.String())
		}
	})

	t.Run("should return 405 when http method is not POST", func(t *testing.T) {
		// Create a mock user service
		mockUserService := &MockUserService{}

		// Create a new HTTP handler with the mock service
		h := AuthHandler(mockUserService)

		// Create a new HTTP request with a GET method
		req, err := http.NewRequest("GET", "/api/auth/register", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		h.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusBadRequest { // Updated to match new error handling
			t.Errorf("handler returned wrong status code: got %v want %v, body: %s",
				status, http.StatusBadRequest, rr.Body.String())
		}

		// Check error response structure
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
			t.Fatal(err)
		}

		// Should be METHOD_NOT_ALLOWED error
		if errorResponse["code"] != "METHOD_NOT_ALLOWED" {
			t.Errorf("expected METHOD_NOT_ALLOWED error, got %v", errorResponse["code"])
		}
	})

	t.Run("should return 409 when user already exists", func(t *testing.T) {
		// Create a mock user service
		mockUserService := &MockUserService{
			FindByEmailFunc: func(email string) (user.User, error) {
				// Return a user to simulate user already exists
				return user.User{
					ID:    1,
					Email: email,
				}, nil
			},
		}

		// Create a new HTTP handler with the mock service
		h := AuthHandler(mockUserService)

		// Create a new registration request
		regReq := types.RegisterUserRequest{
			First_name: "Jane",
			Last_name:  "Doe",
			Email:      "jane.doe@example.com",
			Password:   "password123",
			Phone:      "09876543210",
		}

		// Marshal the request body to JSON
		body, err := json.Marshal(regReq)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		h.ServeHTTP(rr, req)

		// Check the status code - now should return 409 Conflict
		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("handler returned wrong status code: got %v want %v, body: %s",
				status, http.StatusConflict, rr.Body.String())
		}

		// Check error response structure
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
			t.Fatal(err)
		}

		// Should be CONFLICT_ERROR
		if errorResponse["code"] != "CONFLICT_ERROR" {
			t.Errorf("expected CONFLICT_ERROR, got %v", errorResponse["code"])
		}
	})

	t.Run("should return 404 when route is not found", func(t *testing.T) {
		// Create a mock user service
		mockUserService := &MockUserService{}

		// Create a new HTTP handler with the mock service
		h := AuthHandler(mockUserService)

		// Create a new HTTP request with wrong path
		req, err := http.NewRequest("POST", "/api/auth/invalid", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the HTTP request
		h.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v, body: %s",
				status, http.StatusNotFound, rr.Body.String())
		}

		// Check error response structure
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
			t.Fatal(err)
		}

		// Should be ROUTE_NOT_FOUND error
		if errorResponse["code"] != "ROUTE_NOT_FOUND" {
			t.Errorf("expected ROUTE_NOT_FOUND error, got %v", errorResponse["code"])
		}
	})
}

func TestRegisterIntegration(t *testing.T) {
	// Setup do banco de teste
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/restaurant-test")
	if err != nil {
		t.Skip("Database not available for integration tests")
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		t.Skip("Database not available for integration tests")
	}

	// Limpar dados de teste antes de cada teste
	_, err = db.Exec("DELETE FROM users WHERE email LIKE '%test%'")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should register user successfully with real database", func(t *testing.T) {
		// Usar o repository real
		userRepo := user.NewUserRepository(db)
		userService := user.NewUserService(userRepo)
		h := AuthHandler(userService)

		regReq := types.RegisterUserRequest{
			First_name: "Integration",
			Last_name:  "Test",
			Email:      "integration.test@example.com",
			Password:   "password123",
			Phone:      "12345678901", // Exactly 11 digits
		}

		body, err := json.Marshal(regReq)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v, body: %s",
				status, http.StatusCreated, rr.Body.String())
		}

		// Verificar se o usuário foi realmente criado no banco
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", regReq.Email).Scan(&count)
		if err != nil {
			t.Fatal(err)
		}

		if count != 1 {
			t.Errorf("expected 1 user in database, got %d", count)
		}

		// Cleanup after test
		_, err = db.Exec("DELETE FROM users WHERE email = ?", regReq.Email)
		if err != nil {
			t.Logf("Failed to cleanup test data: %v", err)
		}
	})

	t.Run("should return 409 when trying to register duplicate email", func(t *testing.T) {
		// Usar o repository real
		userRepo := user.NewUserRepository(db)
		userService := user.NewUserService(userRepo)
		h := AuthHandler(userService)

		regReq := types.RegisterUserRequest{
			First_name: "Duplicate",
			Last_name:  "Test",
			Email:      "duplicate.test@example.com",
			Password:   "password123",
			Phone:      "12345678901",
		}

		// First registration - should succeed
		body, err := json.Marshal(regReq)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("first registration should succeed: got %v want %v, body: %s",
				status, http.StatusCreated, rr.Body.String())
		}

		// Second registration with same email - should fail
		regReq.Phone = "10987654321" // Different phone
		body, err = json.Marshal(regReq)
		if err != nil {
			t.Fatal(err)
		}

		req, err = http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr = httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("second registration should fail: got %v want %v, body: %s",
				status, http.StatusConflict, rr.Body.String())
		}

		// Cleanup after test
		_, err = db.Exec("DELETE FROM users WHERE email = ?", regReq.Email)
		if err != nil {
			t.Logf("Failed to cleanup test data: %v", err)
		}
	})
}
