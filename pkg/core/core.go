package core

import (
	"github.com/brucejrc/utils/pkg/errorsx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrResponse struct {
	Reason   string            `json:"reason,omitempty"`
	Message  string            `json:"message,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func WriteResponse(c *gin.Context, data any, err error) {
	if err != nil {
		errx := errorsx.FromError(err)
		c.JSON(errx.Code, ErrResponse{
			Reason:   errx.Reason,
			Metadata: errx.Metadata,
			Message:  errx.Message,
		})
		return
	}
	c.JSON(http.StatusOK, data)
}
