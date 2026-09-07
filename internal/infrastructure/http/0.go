package http

import (
	"crypto/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maadiii/taskmanager/pkg/appcontext"
)

func handle[
	IN any, OUT any,
	FN func(*appcontext.Context, *IN) (*OUT, error),
](fn FN, statusCode int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body IN

		if err := bind(c, &body); err != nil {
			// NOTE: use a good error handler
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		requestID := rand.Text()

		ctx := &appcontext.Context{
			Context: c.Request.Context(),
			// UserID: "user_id",
			RequestID: requestID,
		}

		out, err := fn(ctx, &body)
		if err != nil {
			// NOTE: use a good error handler
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		c.Header("request-id", requestID)

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
