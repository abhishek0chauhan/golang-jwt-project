package routes

import (
	"github.com/abhishek0chauhan/golang-jwt-project/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(route *gin.Engine){
	route.POST("users/signup", controllers.Signup())
	route.POST("users/login", controllers.Login())
}