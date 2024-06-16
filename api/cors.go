package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func allowAll() gin.HandlerFunc {
	corsDefault := cors.Default()
	return func(context *gin.Context) {
		corsDefault.HandlerFunc(context.Writer, context.Request)
	}
}
