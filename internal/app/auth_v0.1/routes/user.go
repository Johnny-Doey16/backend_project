package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/steve-mir/diivix_backend/internal/app/auth_v0.1/controllers"
	"github.com/steve-mir/diivix_backend/internal/app/auth_v0.1/middleware"
)

func User(r *gin.Engine) {
	r.Use(middleware.Authenticate())
	r.GET("/users", controllers.GetUsers())
	r.GET("/users/:user_id", controllers.GetUser())
}
