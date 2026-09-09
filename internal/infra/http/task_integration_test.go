package http

import (
	"net/http"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
)

func TestTaskRoutes_CompleteLifecycle(t *testing.T) {
	router := integrationRouter()
	truncateTasks(t)

	createResponse := performRequest(t, router, http.MethodPost, "/api/v1/tasks",
		`{"title":"integration task","description":"created through HTTP"}`)
	assert.Equal(t, http.StatusCreated, createResponse.Code)

	var created struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Status   string `json:"status"`
		Priority string `json:"priority"`
	}
	decodeResponse(t, createResponse, &created)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "integration task", created.Title)
	assert.Equal(t, "TODO", created.Status)
	assert.Equal(t, "MEDIUM", created.Priority)

	getResponse := performRequest(t, router, http.MethodGet, "/api/v1/tasks/"+created.ID, "")
	assert.Equal(t, http.StatusOK, getResponse.Code)
	var fetched struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	decodeResponse(t, getResponse, &fetched)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.Title, fetched.Title)

	listResponse := performRequest(t, router, http.MethodGet, "/api/v1/tasks?limit=10", "")
	assert.Equal(t, http.StatusOK, listResponse.Code)
	var listed struct {
		Tasks []struct {
			ID     string `json:"id"`
			Title  string `json:"title"`
			Status string `json:"status"`
		} `json:"tasks"`
	}
	decodeResponse(t, listResponse, &listed)
	if assert.Len(t, listed.Tasks, 1) {
		assert.Equal(t, created.ID, listed.Tasks[0].ID)
		assert.Equal(t, "integration task", listed.Tasks[0].Title)
		assert.Equal(t, "TODO", listed.Tasks[0].Status)
	}

	updateResponse := performRequest(t, router, http.MethodPut, "/api/v1/tasks/"+created.ID,
		`{"title":"updated integration task","status":"DONE","priority":"HIGH"}`)
	assert.Equal(t, http.StatusOK, updateResponse.Code)
	var updated struct {
		Title    string `json:"title"`
		Status   string `json:"status"`
		Priority string `json:"priority"`
	}
	decodeResponse(t, updateResponse, &updated)
	assert.Equal(t, "updated integration task", updated.Title)
	assert.Equal(t, "DONE", updated.Status)
	assert.Equal(t, "HIGH", updated.Priority)

	deleteResponse := performRequest(t, router, http.MethodDelete, "/api/v1/tasks/"+created.ID, "")
	assert.Equal(t, http.StatusOK, deleteResponse.Code)
	var deleted struct {
		ID      string `json:"id"`
		Deleted bool   `json:"deleted"`
	}
	decodeResponse(t, deleteResponse, &deleted)
	assert.Equal(t, created.ID, deleted.ID)
	assert.True(t, deleted.Deleted)

	assert.Equal(t, http.StatusNotFound,
		performRequest(t, router, http.MethodGet, "/api/v1/tasks/"+created.ID, "").Code)
}

func TestTaskRoutes_ValidationAndOwnership(t *testing.T) {
	router := integrationRouter()
	truncateTasks(t)

	assert.Equal(t, http.StatusBadRequest,
		performRequest(t, router, http.MethodGet, "/api/v1/tasks?limit=-1", "").Code)
	assert.Equal(t, http.StatusBadRequest,
		performRequest(t, router, http.MethodGet, "/api/v1/tasks?limit=101", "").Code)
	assert.Equal(t, http.StatusNotFound,
		performRequest(t, router, http.MethodGet, "/api/v1/tasks/"+uuid.NewV7().String(), "").Code)
	assert.Equal(t, http.StatusNotFound,
		performRequest(t, router, http.MethodPut, "/api/v1/tasks/"+uuid.NewV7().String(),
			`{"title":"missing"}`).Code)
	assert.Equal(t, http.StatusNotFound,
		performRequest(t, router, http.MethodDelete, "/api/v1/tasks/"+uuid.NewV7().String(), "").Code)
}

func TestTaskRoutes_ListFiltersAndCursor(t *testing.T) {
	router := integrationRouter()
	truncateTasks(t)

	createTask(t, router, "alpha", "TODO")
	createTask(t, router, "beta", "DONE")
	createTask(t, router, "alpha second", "TODO")

	response := performRequest(t, router, http.MethodGet, "/api/v1/tasks?status=TODO&limit=1", "")
	assert.Equal(t, http.StatusOK, response.Code)
	var listed struct {
		Tasks []struct {
			ID     string `json:"id"`
			Title  string `json:"title"`
			Status string `json:"status"`
		} `json:"tasks"`
	}
	decodeResponse(t, response, &listed)
	if assert.Len(t, listed.Tasks, 1) {
		assert.Contains(t, listed.Tasks[0].Title, "alpha")
		assert.Equal(t, "TODO", listed.Tasks[0].Status)
	}
}

func createTask(t *testing.T, router http.Handler, title, status string) string {
	t.Helper()
	body := `{"title":"` + title + `","description":"integration"}`
	response := performRequest(t, router, http.MethodPost, "/api/v1/tasks", body)
	assert.Equal(t, http.StatusCreated, response.Code)

	var created struct {
		ID string `json:"id"`
	}
	decodeResponse(t, response, &created)
	assert.NotEmpty(t, created.ID)

	if status != "TODO" {
		update := performRequest(t, router, http.MethodPut, "/api/v1/tasks/"+created.ID,
			`{"status":"`+status+`"}`)
		assert.Equal(t, http.StatusOK, update.Code)
	}

	return created.ID
}
