package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOpenAPISpecIsValidJSON(t *testing.T) {
	var document map[string]any
	assert.NoError(t, json.Unmarshal([]byte(openAPISpec), &document))
	assert.Equal(t, "3.0.3", document["openapi"])
}

func TestDocsRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routeDocs(router, 8000)

	jsonResponse := httptest.NewRecorder()
	router.ServeHTTP(jsonResponse, httptest.NewRequest(http.MethodGet, "/openapi.json", nil))
	assert.Equal(t, http.StatusOK, jsonResponse.Code)
	assert.Equal(t, "application/json; charset=utf-8", jsonResponse.Header().Get("Content-Type"))
	assert.Contains(t, jsonResponse.Body.String(), `"url": "http://localhost:8000"`)

	docsResponse := httptest.NewRecorder()
	router.ServeHTTP(docsResponse, httptest.NewRequest(http.MethodGet, "/docs", nil))
	assert.Equal(t, http.StatusOK, docsResponse.Code)
	assert.Contains(t, docsResponse.Body.String(), "SwaggerUIBundle")
}
