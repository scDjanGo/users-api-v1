package router

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-contrib/cors"

	_ "users-api-v1/src/docs"
	"users-api-v1/src/handlers"
	"users-api-v1/src/media"
	"users-api-v1/src/middlewares"
)

func RouterInit(userHandler *handlers.UserHandler, mediaDir string, maxUploadSize int64) *gin.Engine {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			if name == "-" || name == "" {
				name = fld.Tag.Get("form")
			}
			if comma := strings.Index(name, ","); comma != -1 {
				name = name[:comma]
			}
			if name == "" || name == "-" {
				return fld.Name
			}
			return name
		})
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE",},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	authRouter := router.Group("/auth")
	{
		authRouter.GET("/me", userHandler.GetMe)
		authRouter.POST("/login", userHandler.Login)
		authRouter.POST("/registration", middlewares.LimitBody(maxUploadSize), userHandler.RegisterUser)
	}

	usersRouter := router.Group("/users")
	{
		usersRouter.POST("", middlewares.LimitBody(maxUploadSize), userHandler.Create)
		usersRouter.GET("/:id", userHandler.ReadById)
		usersRouter.GET("", userHandler.ReadAll)
		usersRouter.PATCH("/:id", middlewares.LimitBody(maxUploadSize), userHandler.Update)
		usersRouter.DELETE("/:id", userHandler.Delete)

	}

	mediaGroup := router.Group(media.PublicPrefix(), middlewares.MediaHeaders())
	mediaGroup.StaticFS("/", http.Dir(mediaDir))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
