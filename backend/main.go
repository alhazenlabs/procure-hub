package backend

import (
	"context"
	"github.com/gin-contrib/cors" // Todo Remove in the production version
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

	// Todo remove the below in the prod version of the code
	// Create a CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}

	// Use CORS middleware with the specified configuration
	r.Use(cors.New(config))

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
			users.POST("/v1/login", api.LoginHandler)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
