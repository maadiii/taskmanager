package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleHttpError(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	err := c.Errors.Last().Err

	appErr := new(Error)
	if As(err, &appErr) {
		handleAppError(c, appErr)

		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"key": "INTERNAL_SERVER_ERROR"})
}

func BadRequest(err error) error {
	return WrapCK(err, badRequest, "bad request")
}

func handleAppError(c *gin.Context, err *Error) {
	var status int

	switch err.code {
	case alreadyExists:
		status = http.StatusConflict
	case notFound:
		status = http.StatusNotFound
	case forbidden:
		status = http.StatusForbidden
	case badRequest:
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(status, gin.H{
		"key": err.key,
	})
}
