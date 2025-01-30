package routes

import (
	"ecommerce/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
		auth.POST("/refresh", controllers.RefreshToken)
		auth.POST("/logout", controllers.Logout)
	}
}
