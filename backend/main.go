package backend

import (
	"context"
	docs "github.com/alhazenlabs/procure-hub/backend/docs"
	logger "github.com/alhazenlabs/procure-hub/backend/internal/logger"
	"github.com/alhazenlabs/procure-hub/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	logger.Info("got a request for hello world")
	g.JSON(http.StatusOK, "helloworld")
}

func StartServer() {
	cleanup := middleware.InitTracer()  // Initializing the tracer
	defer cleanup(context.Background()) // cancellation

	logger.Info("starting the backend server")
	r := gin.Default()
	r.Use(otelgin.Middleware(middleware.ServiceName))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")

}
