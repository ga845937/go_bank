package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/v2"
)

func getTraceID(context *gin.Context) string {
	var traceID string
	tx := apm.TransactionFromContext(context.Request.Context())
	if tx != nil {
		traceID = tx.TraceContext().Trace.String()
	}
	return traceID
}

func ResponseOK[T any](context *gin.Context, data T) {
	traceID := getTraceID(context)

	context.JSON(http.StatusOK, gin.H{
		"traceID": traceID,
		"data":    data,
	})
}

func ResponseOKWithoutData(context *gin.Context) {
	traceID := getTraceID(context)

	context.JSON(http.StatusOK, gin.H{
		"traceID": traceID,
	})
}

func ResponseParamError(context *gin.Context, err error) {
	traceID := getTraceID(context)

	context.Error(err)

	context.JSON(http.StatusUnprocessableEntity, gin.H{
		"traceID": traceID,
		"error":   err.Error(),
	})
}

func ResponseServerError(context *gin.Context, err error) {
	traceID := getTraceID(context)

	context.Error(err)

	context.JSON(http.StatusInternalServerError, gin.H{
		"traceID": traceID,
		"error":   err.Error(),
	})
}
