package http

import (
	"bytes"
	"context"
	"database/sql"
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
	"github.com/maadiii/goutils/uow"
	"github.com/maadiii/taskmanager/internal/app/dto"
	"github.com/maadiii/taskmanager/internal/app/port"
	taskservice "github.com/maadiii/taskmanager/internal/app/service/task"
	"github.com/maadiii/taskmanager/internal/domain/task"
	taskrepo "github.com/maadiii/taskmanager/internal/infra/persistence/postgres/task"
	pkgerrors "github.com/maadiii/taskmanager/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
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
		integrationDB, err = newIntegrationPool(ctx, dsn)
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

func newIntegrationPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 20

	return pgxpool.NewWithConfig(context.Background(), config)
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
	tracer := noop.NewTracerProvider().Tracer("http-integration-test")
	repository := contextFreeTaskRepo{repo: taskrepo.NewRepository(integrationDB, tracer)}
	return taskservice.NewService(repository, integrationUoW{
		factory: contextFreeRepoFactory{repo: repository},
	}, noopTaskCache{}, tracer, noopTaskMetrics{})
}

type integrationUoW struct {
	factory port.RepoFactory
}

type contextFreeRepoFactory struct {
	repo port.TaskRepo
}

func (f contextFreeRepoFactory) Tasks() port.TaskRepo {
	return f.repo
}

type contextFreeTaskRepo struct {
	repo port.TaskRepo
}

func (r contextFreeTaskRepo) CreateNew(_ context.Context, entity *task.Entity) error {
	return r.repo.CreateNew(context.Background(), entity)
}

func (r contextFreeTaskRepo) GetTaskByIdAndOwner(_ context.Context, id, ownerID string) (*task.Entity, error) {
	return r.repo.GetTaskByIdAndOwner(context.Background(), id, ownerID)
}

func (r contextFreeTaskRepo) List(_ context.Context, status, userID string, limit int, lastID string) ([]task.Entity, error) {
	return r.repo.List(context.Background(), status, userID, limit, lastID)
}

func (r contextFreeTaskRepo) UpdateByIdAndOwner(_ context.Context, entity *task.Entity) error {
	return r.repo.UpdateByIdAndOwner(context.Background(), entity)
}

func (r contextFreeTaskRepo) DeleteByIdAndOwner(_ context.Context, id, ownerID string) error {
	return r.repo.DeleteByIdAndOwner(context.Background(), id, ownerID)
}

func (u integrationUoW) Do(
	ctx context.Context,
	fn func(context.Context, port.RepoFactory) error,
	_ ...*sql.TxOptions,
) error {
	return fn(context.Background(), u.factory)
}

func (u integrationUoW) Begin(ctx context.Context, _ ...*sql.TxOptions) (context.Context, port.RepoFactory, error) {
	return ctx, u.factory, nil
}

func (integrationUoW) SavePoint(context.Context, string) error {
	return nil
}

func (integrationUoW) Commit(context.Context) error {
	return nil
}

func (integrationUoW) Rollback(context.Context) error {
	return nil
}

var _ uow.UoW[port.RepoFactory] = integrationUoW{}

type noopTaskCache struct{}

type noopTaskMetrics struct{}

func (noopTaskMetrics) IncTasks() {}

func (noopTaskMetrics) DecTasks() {}

func (noopTaskCache) Get(context.Context, string, string) (*dto.GetByIdRs, error) {
	return nil, nil
}

func (noopTaskCache) Set(context.Context, string, string, *dto.GetByIdRs) error {
	return nil
}

func (noopTaskCache) Delete(context.Context, string, string) error {
	return nil
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
