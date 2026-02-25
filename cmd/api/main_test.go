package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "pong - final auto deploy via docker hub & hook success!", response["message"])
}

func TestHealthCheck(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup router
	r := setupRouter()

	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/health", nil)

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "TaskMaster API is healthy", response["message"])
}

func TestLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()

	// Setup dummy test secret
	t.Setenv("JWT_SECRET", "test_secret")

	// Create request body
	body := map[string]string{
		"username": "admin",
		"password": "password",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.True(t, response["success"].(bool))
	
	// Check token existence
	data := response["data"].(map[string]interface{})
	tokenString, ok := data["token"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, tokenString)

	// Verify the token structure
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("test_secret"), nil
	})
	assert.NotNil(t, token)
	assert.True(t, token.Valid)
}

func TestLoginFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()

	// Create request body with invalid credentials
	body := map[string]string{
		"username": "wronguser",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Invalid credentials", response["message"])
}
