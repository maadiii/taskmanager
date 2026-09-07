package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maadiii/taskmanager/internal/application/port"
)

func RouteTaskAPIs(rg *gin.RouterGroup, svc port.TaskService) {
	tasks := rg.Group("/v1/tasks/:id")

	tasks.POST("", handle(svc.Create, http.StatusCreated))
}
