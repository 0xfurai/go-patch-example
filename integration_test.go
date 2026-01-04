package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang-http-patch/database"
	"golang-http-patch/handlers"
	"golang-http-patch/models"
	"golang-http-patch/validation"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates a temporary in-memory database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = models.AutoMigrate(db)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// setupTestRouter creates a router with test database
func setupTestRouter(t *testing.T, db *gorm.DB) *mux.Router {
	// Set the global DB for handlers to use
	database.DB = db

	r := mux.NewRouter()
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handlers.PatchUser).Methods("PATCH")
	return r
}

func TestMain(m *testing.M) {
	// Initialize validator once
	validation.InitValidator()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(t, db)

	// Test data
	createReq := models.CreateUserDTO{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    30,
		Phone:  stringPtr("1234567890"),
		Active: true,
		Bio:    "Software developer",
		Role:   "user",
		Score:  85.5,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var user models.User
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
	if user.Name != createReq.Name {
		t.Errorf("Expected name %s, got %s", createReq.Name, user.Name)
	}
	if user.Email != createReq.Email {
		t.Errorf("Expected email %s, got %s", createReq.Email, user.Email)
	}
	if user.Age != createReq.Age {
		t.Errorf("Expected age %d, got %d", createReq.Age, user.Age)
	}
	if user.Phone == nil || *user.Phone != *createReq.Phone {
		t.Errorf("Expected phone %s, got %v", *createReq.Phone, user.Phone)
	}
	if user.Active != createReq.Active {
		t.Errorf("Expected active %v, got %v", createReq.Active, user.Active)
	}
	if user.Bio != createReq.Bio {
		t.Errorf("Expected bio %s, got %s", createReq.Bio, user.Bio)
	}
	if user.Role != createReq.Role {
		t.Errorf("Expected role %s, got %s", createReq.Role, user.Role)
	}
	if user.Score != createReq.Score {
		t.Errorf("Expected score %f, got %f", createReq.Score, user.Score)
	}
}

func TestGetUser(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(t, db)

	// First create a user
	user := models.User{
		Name:   "Jane Doe",
		Email:  "jane@example.com",
		Age:    25,
		Phone:  stringPtr("9876543210"),
		Active: true,
		Bio:    "Test bio",
		Role:   "admin",
		Score:  92.0,
	}
	db.Create(&user)

	// Test getting the user
	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d", user.ID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var retrievedUser models.User
	if err := json.Unmarshal(w.Body.Bytes(), &retrievedUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, retrievedUser.ID)
	}
	if retrievedUser.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, retrievedUser.Name)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrievedUser.Email)
	}

	// Test getting non-existent user
	req = httptest.NewRequest("GET", "/users/99999", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(t, db)

	// Create a user first
	user := models.User{
		Name:   "Original Name",
		Email:  "original@example.com",
		Age:    20,
		Phone:  stringPtr("1111111111"),
		Active: false,
		Bio:    "Original bio",
		Role:   "user",
		Score:  50.0,
	}
	db.Create(&user)

	// Update the user
	updateReq := models.UpdateUserDTO{
		Name:   "Updated Name",
		Age:    35,
		Phone:  stringPtr("2222222222"),
		Active: true,
		Bio:    "Updated bio",
		Role:   "admin",
		Score:  95.0,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var updatedUser models.User
	if err := json.Unmarshal(w.Body.Bytes(), &updatedUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify all fields were updated
	if updatedUser.Name != updateReq.Name {
		t.Errorf("Expected name %s, got %s", updateReq.Name, updatedUser.Name)
	}
	if updatedUser.Age != updateReq.Age {
		t.Errorf("Expected age %d, got %d", updateReq.Age, updatedUser.Age)
	}
	if updatedUser.Phone == nil || *updatedUser.Phone != *updateReq.Phone {
		t.Errorf("Expected phone %s, got %v", *updateReq.Phone, updatedUser.Phone)
	}
	if updatedUser.Active != updateReq.Active {
		t.Errorf("Expected active %v, got %v", updateReq.Active, updatedUser.Active)
	}
	if updatedUser.Bio != updateReq.Bio {
		t.Errorf("Expected bio %s, got %s", updateReq.Bio, updatedUser.Bio)
	}
	if updatedUser.Role != updateReq.Role {
		t.Errorf("Expected role %s, got %s", updateReq.Role, updatedUser.Role)
	}
	if updatedUser.Score != updateReq.Score {
		t.Errorf("Expected score %f, got %f", updateReq.Score, updatedUser.Score)
	}

	// Verify email was not changed (immutable)
	if updatedUser.Email != user.Email {
		t.Errorf("Email should be immutable, expected %s, got %s", user.Email, updatedUser.Email)
	}
}

func TestPatchUser_UnsetFields(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(t, db)

	// Create a user with initial values
	user := models.User{
		Name:   "Initial Name",
		Email:  "initial@example.com",
		Age:    25,
		Phone:  stringPtr("1111111111"),
		Active: true,
		Bio:    "Initial bio",
		Role:   "user",
		Score:  75.0,
	}
	db.Create(&user)

	// PATCH with only one field (others are unset - should remain unchanged)
	patchReq := map[string]interface{}{
		"name": "Updated Name Only",
		// All other fields are unset (not in JSON)
	}

	body, _ := json.Marshal(patchReq)
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var patchedUser models.User
	if err := json.Unmarshal(w.Body.Bytes(), &patchedUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Only name should be updated
	if patchedUser.Name != "Updated Name Only" {
		t.Errorf("Expected name 'Updated Name Only', got %s", patchedUser.Name)
	}

	// All other fields should remain unchanged
	if patchedUser.Age != user.Age {
		t.Errorf("Age should remain unchanged: expected %d, got %d", user.Age, patchedUser.Age)
	}
	if patchedUser.Phone == nil || *patchedUser.Phone != *user.Phone {
		t.Errorf("Phone should remain unchanged: expected %v, got %v", user.Phone, patchedUser.Phone)
	}
	if patchedUser.Active != user.Active {
		t.Errorf("Active should remain unchanged: expected %v, got %v", user.Active, patchedUser.Active)
	}
	if patchedUser.Bio != user.Bio {
		t.Errorf("Bio should remain unchanged: expected %s, got %s", user.Bio, patchedUser.Bio)
	}
	if patchedUser.Role != user.Role {
		t.Errorf("Role should remain unchanged: expected %s, got %s", user.Role, patchedUser.Role)
	}
	if patchedUser.Score != user.Score {
		t.Errorf("Score should remain unchanged: expected %f, got %f", user.Score, patchedUser.Score)
	}
	if patchedUser.Email != user.Email {
		t.Errorf("Email should remain unchanged: expected %s, got %s", user.Email, patchedUser.Email)
	}
}

func TestPatchUser_NullFields(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(t, db)

	// Create a user with all fields set
	user := models.User{
		Name:   "Test User",
		Email:  "test@example.com",
		Age:    30,
		Phone:  stringPtr("1234567890"),
		Active: true,
		Bio:    "Some bio",
		Role:   "admin",
		Score:  80.0,
	}
	db.Create(&user)

	// PATCH with null values - should set fields to null/zero values
	patchReq := map[string]interface{}{
		"phone":  nil, // Nullable field - should be set to null
		"bio":    nil, // String field - should be set to empty string
		"active": nil, // Boolean field - should be set to false
		"score":  nil, // Float field - should be set to 0.0
	}

	body, _ := json.Marshal(patchReq)
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var patchedUser models.User
	if err := json.Unmarshal(w.Body.Bytes(), &patchedUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Phone should be null (nil pointer)
	if patchedUser.Phone != nil {
		t.Errorf("Phone should be null, got %v", patchedUser.Phone)
	}

	// Bio should be empty string
	if patchedUser.Bio != "" {
		t.Errorf("Bio should be empty string, got %s", patchedUser.Bio)
	}

	// Active should be false (zero value for bool)
	if patchedUser.Active != false {
		t.Errorf("Active should be false, got %v", patchedUser.Active)
	}

	// Score should be 0.0
	if patchedUser.Score != 0.0 {
		t.Errorf("Score should be 0.0, got %f", patchedUser.Score)
	}

	// Fields not in patch should remain unchanged
	if patchedUser.Name != user.Name {
		t.Errorf("Name should remain unchanged: expected %s, got %s", user.Name, patchedUser.Name)
	}
	if patchedUser.Age != user.Age {
		t.Errorf("Age should remain unchanged: expected %d, got %d", user.Age, patchedUser.Age)
	}
	if patchedUser.Role != user.Role {
		t.Errorf("Role should remain unchanged: expected %s, got %s", user.Role, patchedUser.Role)
	}
}

func TestPatchUser_ValueFields(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(t, db)

	// Create a user with initial values
	user := models.User{
		Name:   "Original Name",
		Email:  "original@example.com",
		Age:    20,
		Phone:  stringPtr("1111111111"),
		Active: false,
		Bio:    "Original bio",
		Role:   "user",
		Score:  50.0,
	}
	db.Create(&user)

	// PATCH with new values
	patchReq := map[string]interface{}{
		"name":   "New Name",
		"age":    35,
		"phone":  "9999999999",
		"active": true,
		"bio":    "New bio text",
		"role":   "admin",
		"score":  95.5,
	}

	body, _ := json.Marshal(patchReq)
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var patchedUser models.User
	if err := json.Unmarshal(w.Body.Bytes(), &patchedUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify all fields were updated to new values
	if patchedUser.Name != "New Name" {
		t.Errorf("Expected name 'New Name', got %s", patchedUser.Name)
	}
	if patchedUser.Age != 35 {
		t.Errorf("Expected age 35, got %d", patchedUser.Age)
	}
	if patchedUser.Phone == nil || *patchedUser.Phone != "9999999999" {
		t.Errorf("Expected phone '9999999999', got %v", patchedUser.Phone)
	}
	if patchedUser.Active != true {
		t.Errorf("Expected active true, got %v", patchedUser.Active)
	}
	if patchedUser.Bio != "New bio text" {
		t.Errorf("Expected bio 'New bio text', got %s", patchedUser.Bio)
	}
	if patchedUser.Role != "admin" {
		t.Errorf("Expected role 'admin', got %s", patchedUser.Role)
	}
	if patchedUser.Score != 95.5 {
		t.Errorf("Expected score 95.5, got %f", patchedUser.Score)
	}

	// Email should remain unchanged (immutable)
	if patchedUser.Email != user.Email {
		t.Errorf("Email should remain unchanged: expected %s, got %s", user.Email, patchedUser.Email)
	}
}

func TestPatchUser_MixedUnsetNullValue(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(t, db)

	// Create a user with all fields set
	user := models.User{
		Name:   "Mixed Test",
		Email:  "mixed@example.com",
		Age:    25,
		Phone:  stringPtr("1111111111"),
		Active: true,
		Bio:    "Original bio",
		Role:   "user",
		Score:  70.0,
	}
	db.Create(&user)

	// PATCH with mix of unset, null, and value
	// - name: unset (should remain unchanged)
	// - age: value (should be updated)
	// - phone: null (should be set to null)
	// - active: value (should be updated)
	// - bio: unset (should remain unchanged)
	// - role: null (should be set to empty string)
	// - score: value (should be updated)
	patchReq := map[string]interface{}{
		"age":    40,    // value
		"phone":  nil,   // null
		"active": false, // value
		"role":   nil,   // null (will be empty string)
		"score":  88.5,  // value
		// name and bio are unset
	}

	body, _ := json.Marshal(patchReq)
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var patchedUser models.User
	if err := json.Unmarshal(w.Body.Bytes(), &patchedUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Unset fields should remain unchanged
	if patchedUser.Name != user.Name {
		t.Errorf("Name (unset) should remain unchanged: expected %s, got %s", user.Name, patchedUser.Name)
	}
	if patchedUser.Bio != user.Bio {
		t.Errorf("Bio (unset) should remain unchanged: expected %s, got %s", user.Bio, patchedUser.Bio)
	}

	// Null fields should be set to null/zero
	if patchedUser.Phone != nil {
		t.Errorf("Phone (null) should be null, got %v", patchedUser.Phone)
	}
	if patchedUser.Role != "" {
		t.Errorf("Role (null) should be empty string, got %s", patchedUser.Role)
	}

	// Value fields should be updated
	if patchedUser.Age != 40 {
		t.Errorf("Age (value) should be 40, got %d", patchedUser.Age)
	}
	if patchedUser.Active != false {
		t.Errorf("Active (value) should be false, got %v", patchedUser.Active)
	}
	if patchedUser.Score != 88.5 {
		t.Errorf("Score (value) should be 88.5, got %f", patchedUser.Score)
	}
}

func TestPatchUser_EmptyPatch(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(t, db)

	// Create a user
	user := models.User{
		Name:   "Test User",
		Email:  "test@example.com",
		Age:    30,
		Active: true,
	}
	db.Create(&user)

	// PATCH with empty body (no fields)
	patchReq := map[string]interface{}{}

	body, _ := json.Marshal(patchReq)
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for empty patch, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestPatchUser_NonExistentUser(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(t, db)

	// Try to patch non-existent user
	patchReq := map[string]interface{}{
		"name": "New Name",
	}

	body, _ := json.Marshal(patchReq)
	req := httptest.NewRequest("PATCH", "/users/99999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for non-existent user, got %d. Body: %s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestPatchUser_ValidationRules(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(t, db)

	// Create a user for testing
	user := models.User{
		Name:   "Valid Name",
		Email:  "valid@example.com",
		Age:    25,
		Active: true,
		Role:   "user",
		Score:  75.0,
	}
	db.Create(&user)

	tests := []struct {
		name           string
		patchData      map[string]interface{}
		expectedStatus int
		description    string
	}{
		{
			name:           "Name too short",
			patchData:      map[string]interface{}{"name": "A"},
			expectedStatus: http.StatusBadRequest,
			description:    "Name with 1 character should fail min=2",
		},
		{
			name:           "Name too long",
			patchData:      map[string]interface{}{"name": string(make([]byte, 101))}, // 101 chars
			expectedStatus: http.StatusBadRequest,
			description:    "Name with >100 characters should fail max=100",
		},
		{
			name:           "Name valid min boundary",
			patchData:      map[string]interface{}{"name": "AB"},
			expectedStatus: http.StatusOK,
			description:    "Name with exactly 2 characters should pass",
		},
		{
			name:           "Age negative",
			patchData:      map[string]interface{}{"age": -1},
			expectedStatus: http.StatusBadRequest,
			description:    "Negative age should fail gte=0",
		},
		{
			name:           "Age too large",
			patchData:      map[string]interface{}{"age": 151},
			expectedStatus: http.StatusBadRequest,
			description:    "Age >150 should fail lte=150",
		},
		{
			name:           "Age valid boundaries",
			patchData:      map[string]interface{}{"age": 0},
			expectedStatus: http.StatusOK,
			description:    "Age 0 should pass gte=0",
		},
		{
			name:           "Age valid max boundary",
			patchData:      map[string]interface{}{"age": 150},
			expectedStatus: http.StatusOK,
			description:    "Age 150 should pass lte=150",
		},
		{
			name:           "Phone too short",
			patchData:      map[string]interface{}{"phone": "123456789"}, // 9 chars
			expectedStatus: http.StatusBadRequest,
			description:    "Phone with <10 characters should fail min=10",
		},
		{
			name:           "Phone too long",
			patchData:      map[string]interface{}{"phone": "123456789012345678901"}, // 21 chars
			expectedStatus: http.StatusBadRequest,
			description:    "Phone with >20 characters should fail max=20",
		},
		{
			name:           "Phone valid",
			patchData:      map[string]interface{}{"phone": "1234567890"}, // 10 chars
			expectedStatus: http.StatusOK,
			description:    "Phone with 10 characters should pass",
		},
		{
			name:           "Role invalid",
			patchData:      map[string]interface{}{"role": "superadmin"},
			expectedStatus: http.StatusBadRequest,
			description:    "Role not in enum should fail oneof=admin user guest",
		},
		{
			name:           "Role valid admin",
			patchData:      map[string]interface{}{"role": "admin"},
			expectedStatus: http.StatusOK,
			description:    "Role 'admin' should pass",
		},
		{
			name:           "Role valid user",
			patchData:      map[string]interface{}{"role": "user"},
			expectedStatus: http.StatusOK,
			description:    "Role 'user' should pass",
		},
		{
			name:           "Role valid guest",
			patchData:      map[string]interface{}{"role": "guest"},
			expectedStatus: http.StatusOK,
			description:    "Role 'guest' should pass",
		},
		{
			name:           "Score negative",
			patchData:      map[string]interface{}{"score": -0.1},
			expectedStatus: http.StatusBadRequest,
			description:    "Negative score should fail gte=0",
		},
		{
			name:           "Score too large",
			patchData:      map[string]interface{}{"score": 100.1},
			expectedStatus: http.StatusBadRequest,
			description:    "Score >100 should fail lte=100",
		},
		{
			name:           "Score valid min",
			patchData:      map[string]interface{}{"score": 0.0},
			expectedStatus: http.StatusOK,
			description:    "Score 0.0 should pass",
		},
		{
			name:           "Score valid max",
			patchData:      map[string]interface{}{"score": 100.0},
			expectedStatus: http.StatusOK,
			description:    "Score 100.0 should pass",
		},
		{
			name:           "Bio too long",
			patchData:      map[string]interface{}{"bio": string(make([]byte, 501))}, // 501 chars
			expectedStatus: http.StatusBadRequest,
			description:    "Bio with >500 characters should fail max=500",
		},
		{
			name:           "Bio valid max",
			patchData:      map[string]interface{}{"bio": string(make([]byte, 500))}, // 500 chars
			expectedStatus: http.StatusOK,
			description:    "Bio with 500 characters should pass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.patchData)
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: Expected status %d, got %d. Body: %s",
					tt.description, tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
