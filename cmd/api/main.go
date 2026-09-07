package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maadiii/taskmanager/internal/application/service/task"
	"github.com/maadiii/taskmanager/internal/infrastructure/http"
)

func main() {
	svc := task.NewService(nil)

	g := gin.Default()
	rg := g.Group("/api")

	http.RouteTaskAPIs(rg, svc)

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
