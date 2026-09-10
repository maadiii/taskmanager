package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
)

const (
	benchmarkUserID = "01a081e5-87cb-7f89-bbbc-5a6cd8cae373"
	benchmarkRows   = 10000
)

func BenchmarkTaskListCursorPagination(b *testing.B) {
	router := integrationRouter()
	seedBenchmarkTasks(b, benchmarkRows)
	lastID := benchmarkCursor(router, b)

	b.Run("first-page", func(b *testing.B) {
		benchmarkTaskListRequest(b, router, "/api/v1/tasks?status=TODO&limit=20")
	})
	b.Run("cursor-page", func(b *testing.B) {
		benchmarkTaskListRequest(b, router, "/api/v1/tasks?status=TODO&limit=20&last_id="+lastID)
	})
}

func seedBenchmarkTasks(b *testing.B, count int) {
	b.Helper()

	_, err := integrationDB.Exec(context.Background(), "TRUNCATE TABLE tasks")
	if err != nil {
		b.Fatal(err)
	}

	rows := make([][]any, 0, count)
	now := time.Now().UTC()
	for i := 0; i < count; i++ {
		id := uuid.NewV7().String()
		status := "TODO"
		if i%2 == 0 {
			status = "DONE"
		}
		rows = append(rows, []any{
			id,
			benchmarkUserID,
			"benchmark-task-" + id,
			"benchmark task",
			status,
			"MEDIUM",
			now,
			now,
		})
	}

	_, err = integrationDB.CopyFrom(
		context.Background(),
		pgx.Identifier{"tasks"},
		[]string{"id", "user_id", "title", "description", "status", "priority", "created_at", "updated_at"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		b.Fatal(err)
	}
}

func benchmarkCursor(router http.Handler, b *testing.B) string {
	b.Helper()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?status=TODO&limit=20", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		b.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body struct {
		Tasks []struct {
			ID string `json:"id"`
		} `json:"tasks"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		b.Fatal(err)
	}
	if len(body.Tasks) == 0 {
		b.Fatal("expected benchmark seed to produce tasks")
	}

	return body.Tasks[len(body.Tasks)-1].ID
}

func benchmarkTaskListRequest(b *testing.B, router http.Handler, path string) {
	b.Helper()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			b.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
		}
	}
}
