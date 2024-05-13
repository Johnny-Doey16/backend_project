package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/steve-mir/diivix_backend/internal/app/auth_v0.1/controllers"
)

func Auth(r *gin.Engine) {
	r.POST("user/register", controllers.Register())
	r.POST("user/login", controllers.Login())
}
