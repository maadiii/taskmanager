package metrics

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareRecordsRequestMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m, err := New(&metricsExecer{count: 3})
	if !assert.NoError(t, err) {
		return
	}

	router := gin.New()
	router.Use(m.Middleware())
	router.GET("/tasks/:id", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})

	request := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
	families, err := m.registry.Gather()
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, families, 3) {
		return
	}

	assertMetricFamilyExists(t, families, "requests_total")
	assertMetricFamilyExists(t, families, "request_latency_histogram")
	assertMetricFamilyExists(t, families, "tasks_count")

	m.IncTasks()
	m.DecTasks()
}

func assertMetricFamilyExists(
	t *testing.T,
	families []*dto.MetricFamily,
	name string,
) {
	t.Helper()
	for _, family := range families {
		if family.GetName() == name {
			return
		}
	}
	t.Fatalf("metric family %q was not collected", name)
}

type metricsExecer struct {
	count int64
}

func (e *metricsExecer) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (e *metricsExecer) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, nil
}

func (e *metricsExecer) QueryRow(context.Context, string, ...any) pgx.Row {
	return metricsRow{value: e.count}
}

func (e *metricsExecer) BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error) {
	return nil, nil
}

type metricsRow struct {
	value int64
}

func (r metricsRow) Scan(dest ...any) error {
	*dest[0].(*int64) = r.value
	return nil
}
