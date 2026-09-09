package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maadiii/taskmanager/internal/app/port"
)

func routeTaskV1(rg *gin.RouterGroup, svc port.TaskService) {
	v1 := rg.Group("/v1/tasks")
	v1Tasks := v1.Use(authorize())

	v1Tasks.POST("", handle(svc.Create, http.StatusCreated))
	v1Tasks.GET("/:id", handle(svc.GetByID, http.StatusOK))
	v1Tasks.PUT("/:id", handle(svc.Update, http.StatusOK))
	v1Tasks.DELETE("/:id", handle(svc.Delete, http.StatusOK))
}
