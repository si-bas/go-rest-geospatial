package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/si-bas/go-rest-geospatial/shared/helper/response"
)

func (h *Handler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	result := response.NewJSONResponse(ctx)

	c.JSON(result.APIStatusSuccess().StatusCode, result.SetData("Success"))
}
