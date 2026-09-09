package http

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maadiii/taskmanager/config"
	"github.com/maadiii/taskmanager/internal/app/port"
	"github.com/maadiii/taskmanager/pkg/appcontext"
	"github.com/maadiii/taskmanager/pkg/errors"
	"go.uber.org/fx"
)

func Route(p Param) {
	routeTaskV1(p.RouterGroup, p.TaskSvc)
}

type Param struct {
	fx.In

	RouterGroup *gin.RouterGroup

	TaskSvc port.TaskService
}

func NewApiGroupRouter(lc fx.Lifecycle, cfg *config.Config) *gin.RouterGroup {
	g := gin.Default()
	g.Use(errors.HandleHttpError)
	rg := g.Group("/api")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: g.Handler(),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("server is running...")

				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("failed on server running")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("server is shutting down...")

			if err := srv.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed on server shutdown: %w", err)
			}

			return nil
		},
	})

	return rg
}

func handle[
	IN any, OUT any,
	FN func(*appcontext.Context, *IN) (*OUT, error),
](fn FN, statusCode int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body IN

		if err := bind(c, &body); err != nil {
			_ = c.Error(errors.BadRequest(err))
			c.Abort()

			return
		}

		requestID := rand.Text()
		c.Header("request-id", requestID)

		ctx := &appcontext.Context{
			Context:   c.Request.Context(),
			RequestID: requestID,
			Identity: appcontext.Identity{
				UserID:      c.GetString(userIdKey),
				Role:        appcontext.Role(c.GetString(roleKey)),
				Permissions: c.GetStringSlice(permissionsKey),
			},
		}
		out, err := fn(ctx, &body)
		if err != nil {
			_ = c.Error(err)
			c.Abort()

			return
		}

		c.JSON(statusCode, out)
	}
}

func bind(c *gin.Context, obj any) error {
	// URI parameters
	if err := c.ShouldBindUri(obj); err != nil {
		return err
	}

	// Headers
	if err := c.ShouldBindHeader(obj); err != nil {
		return err
	}

	// Query parameters
	if err := c.ShouldBindQuery(obj); err != nil {
		return err
	}

	// Body
	if hasBody(c) {
		if err := c.ShouldBind(obj); err != nil {
			return err
		}
	}

	return nil
}

func hasBody(c *gin.Context) bool {
	return c.Request.Body != nil &&
		c.Request.ContentLength != 0
}

const (
	userIdKey      = "USER_ID"
	roleKey        = "ROLES"
	permissionsKey = "PERM"
)

// Simple authorization middleware
func authorize(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get userId by real world authorization mechanism(seesion, jwt, paseto .etc)
		ctx.Set(userIdKey, "01a081e5-87cb-7f89-bbbc-5a6cd8cae373")

		// Get roles by real world authorization mechanism(seesion, jwt, paseto .etc)
		ctx.Set(roleKey, roles)

		// Get permissions by real world authorization mechanism(seesion, jwt, paseto .etc)
		// Permissions checking must be in the business layer, it just get permissions and sets to ctx
		ctx.Set(permissionsKey, []string{"create", "read", "update", "delete"})
	}
}
