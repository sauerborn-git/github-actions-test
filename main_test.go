package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAlbums(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.Default()
	router.GET("/albums", getAlbums)

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, "/albums", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check if the status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the response body contains the correct albums data
	expectedBody := `[{"id":"1","title":"Blue Train","artist":"John Coltrane","price":56.99},{"id":"2","title":"Jeru","artist":"Gerry Mulligan","price":17.99},{"id":"3","title":"Sarah Vaughan and Clifford Brown","artist":"Sarah Vaughan","price":39.99}]`
	assert.JSONEq(t, expectedBody, w.Body.String())
}
