package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maadiii/taskmanager/internal/app/port"
	taskservice "github.com/maadiii/taskmanager/internal/app/service/task"
	"github.com/maadiii/taskmanager/internal/infra/persistence/postgres"
	taskrepo "github.com/maadiii/taskmanager/internal/infra/persistence/postgres/task"
	pkgerrors "github.com/maadiii/taskmanager/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	postgresUser     = "postgres"
	postgresPassword = "postgres"
	postgresDatabase = "taskmanager"
)

var integrationDB *pgxpool.Pool

func TestMain(m *testing.M) {
	if err := startIntegrationDatabase(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code := m.Run()
	stopIntegrationDatabase()
	os.Exit(code)
}

func startIntegrationDatabase() error {
	containerName := integrationContainerName()
	_ = exec.Command("docker", "rm", "-f", containerName).Run()

	port := integrationDatabasePort()
	args := []string{
		"run", "--detach", "--rm",
		"--name", containerName,
		"-e", "POSTGRES_USER=" + postgresUser,
		"-e", "POSTGRES_PASSWORD=" + postgresPassword,
		"-e", "POSTGRES_DB=" + postgresDatabase,
		"-p", port + ":5432",
		"postgres:16-alpine",
	}
	if output, err := exec.Command("docker", args...).CombinedOutput(); err != nil {
		return fmt.Errorf("start postgres container: %w: %s", err, output)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@127.0.0.1:%s/%s?sslmode=disable",
		postgresUser, postgresPassword, port, postgresDatabase,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var err error
	for ctx.Err() == nil {
		integrationDB, err = pgxpool.New(ctx, dsn)
		if err == nil {
			err = integrationDB.Ping(ctx)
			if err == nil {
				break
			}
			integrationDB.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		stopIntegrationDatabase()
		return fmt.Errorf("connect to postgres container: %w", err)
	}

	_, sourceFile, _, _ := runtime.Caller(0)
	migrationPath := filepath.Join(filepath.Dir(sourceFile), "../../../setup/migrations/000001_tasks_table.up.sql")
	migration, err := os.ReadFile(migrationPath)
	if err != nil {
		stopIntegrationDatabase()
		return fmt.Errorf("read migration: %w", err)
	}
	if _, err := integrationDB.Exec(context.Background(), string(migration)); err != nil {
		stopIntegrationDatabase()
		return fmt.Errorf("run migration: %w", err)
	}

	return nil
}

func stopIntegrationDatabase() {
	if integrationDB != nil {
		integrationDB.Close()
	}
	_ = exec.Command("docker", "rm", "-f", integrationContainerName()).Run()
}

func integrationContainerName() string {
	return "taskmanager-integration-postgres"
}

func integrationDatabasePort() string {
	if port := os.Getenv("TEST_POSTGRES_PORT"); port != "" {
		return port
	}
	return "55432"
}

func integrationRouter() http.Handler {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(pkgerrors.HandleHttpError)
	routeTaskV1(engine.Group("/api"), newTaskService())
	return engine
}

func newTaskService() port.TaskService {
	repository := taskrepo.NewRepository(integrationDB)
	return taskservice.NewService(repository, postgres.NewUoW(integrationDB))
}

func truncateTasks(t *testing.T) {
	t.Helper()
	_, err := integrationDB.Exec(context.Background(), "TRUNCATE TABLE tasks")
	assert.NoError(t, err)
}

func performRequest(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var requestBody *bytes.Reader
	if body == "" {
		requestBody = bytes.NewReader(nil)
	} else {
		requestBody = bytes.NewReader([]byte(body))
	}
	request := httptest.NewRequest(method, path, requestBody)
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func decodeResponse(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	assert.NoError(t, json.NewDecoder(response.Body).Decode(target))
}
