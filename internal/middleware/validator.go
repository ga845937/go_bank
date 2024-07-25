package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Validator[T any](next func(*gin.Context, *T)) gin.HandlerFunc {
	return func(context *gin.Context) {
		param := new(T)
		b := binding.Default(context.Request.Method, context.ContentType())
		if err := context.ShouldBindWith(&param, b); err != nil {
			ResponseParamError(context, err)
			return
		}

		next(context, param)
	}
}
