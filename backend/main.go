package backend

import (
	"context"
	"net/http"

	"github.com/alhazenlabs/procure-hub/backend/docs"
	"github.com/alhazenlabs/procure-hub/backend/internal/api"
	"github.com/alhazenlabs/procure-hub/backend/internal/database"
	"github.com/alhazenlabs/procure-hub/backend/internal/logger"
	"github.com/alhazenlabs/procure-hub/backend/internal/middleware"
	// "github.com/alhazenlabs/procure-hub/backend/internal/models"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var BasePath = "/products/procurehub"

// @BasePath /products/procure-hub

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
	logger.Info("configuring middleware")
	cleanup := middleware.InitTracer()  // Initializing the tracer
	defer cleanup(context.Background()) // cancellation

	logger.Info("running migrations")
	database.RunMigrations()

	logger.Info("starting the backend server")
	r := gin.Default()
	r.Use(otelgin.Middleware(middleware.ServiceName))

	docs.SwaggerInfo.BasePath = BasePath
	v1 := r.Group(BasePath)
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
		}
		users := v1.Group("/users")
		{
			users.POST("/v1/signup", api.RegisterUser)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
